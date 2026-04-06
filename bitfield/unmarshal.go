package bitfield

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const (
	FlagTreatEndAsZero = 1 << iota
)

// UnmarshalHexString deserializes data encoded as a hex string into
// the provided pointer to a struct/array/slice/value.
// It returns the number of bytes consumed and an error, if any.
// The unmarshaling behavior is the same as [Unmarshal].
func UnmarshalHexString(hexData string, out any, flags int) (int, error) {
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return 0, err
	}
	return Unmarshal(data, out, flags)
}

// Unmarshal deserializes data into the provided pointer to a struct/array/slice/value.
// It returns the number of bytes consumed and an error, if any.
func Unmarshal(data []byte, out any, flags int) (int, error) {
	if out == nil {
		return 0, errors.New("nil output")
	}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return 0, errors.New("out must be a non-nil pointer")
	}
	return doUnmarshal(data, rv.Elem(), flags)
}

// doUnmarshal fills rv from data, returning bytes consumed.
func doUnmarshal(data []byte, rv reflect.Value, flags int) (int, error) {
	origLen := len(data)
	// If pointer, dereference.
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		rt := rv.Type()
		i := 0
		for i < rt.NumField() {
			fieldStruct := rt.Field(i)
			fieldValue := rv.Field(i)

			if fieldStruct.Tag.Get("bitfield") == "-" {
				i++
				continue
			}

			// Expect nested struct (or pointer to struct).
			if fieldValue.Kind() == reflect.Pointer {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			// union: read data for members and OR-merge into the destination struct fields.
			if fieldStruct.Tag.Get("bitfield") == "union" {
				if fieldValue.Kind() != reflect.Struct {
					return 0, fmt.Errorf("bitfield:\"union\" on non-struct field %s", fieldStruct.Name)
				}

				maxConsumed := 0
				for j := range fieldValue.NumField() {
					memberStruct := fieldValue.Type().Field(j)
					memberValue := fieldValue.Field(j)

					if memberStruct.Tag.Get("bitfield") == "-" {
						continue
					}

					// If member is itself a union, process it recursively.
					if memberStruct.Tag.Get("bitfield") == "union" {
						if memberValue.Kind() == reflect.Pointer {
							if memberValue.IsNil() {
								memberValue.Set(reflect.New(memberValue.Type().Elem()))
							}
							memberValue = memberValue.Elem()
						}
						if memberValue.Kind() != reflect.Struct {
							return 0, fmt.Errorf("bitfield:\"union\" on non-struct member %s.%s",
								fieldStruct.Name, memberStruct.Name)
						}

						nestedMaxConsumed := 0
						for k := range memberValue.NumField() {
							nestedMemberStruct := memberValue.Type().Field(k)
							if nestedMemberStruct.Tag.Get("bitfield") == "-" {
								continue
							}
							consumed, err := doUnmarshal(data, memberValue.Field(k), flags)
							if err != nil {
								return 0, fmt.Errorf("bitfield:\"union\" on ObjectID member: %w", err)
							}
							if nestedMaxConsumed < consumed {
								nestedMaxConsumed = consumed
							}
						}

						if maxConsumed < nestedMaxConsumed {
							maxConsumed = nestedMaxConsumed
						}
						continue
					}

					// Attempt to unmarshal from the full remaining data.
					consumed, err := doUnmarshal(data, fieldValue.Field(j), flags)
					if err != nil {
						return 0, fmt.Errorf("bitfield:\"union\" on %s member: %w", memberStruct.Name, err)
					}
					if maxConsumed < consumed {
						maxConsumed = consumed
					}
				}

				i++

				data = safeAdvance(data, maxConsumed)

				continue
			}

			// Bitfield group? handle consecutive bitfield-tagged fields together.
			if _, ok := fieldStruct.Tag.Lookup("bitfield"); ok {
				// We'll consume bytes as needed to fill consecutive bitfields.
				bitIdx := 0 // index in data
				bitBuffer := byte(0)
				bitCount := 0 // bits available in bitBuffer (already loaded)
				loadByte := func() error {
					if bitIdx >= len(data) {
						if flags&FlagTreatEndAsZero == 0 {
							return errors.New("not enough data for bitfield")
						}
						bitBuffer = 0
					} else {
						bitBuffer = data[bitIdx]
					}
					bitIdx++
					bitCount = 8
					return nil
				}

				for i < rt.NumField() {
					fieldStruct = rt.Field(i)
					tagStr, isBitField := fieldStruct.Tag.Lookup("bitfield")
					if !isBitField {
						break
					}
					if tagStr == "-" {
						i++
						continue
					}
					width, err := strconv.Atoi(tagStr)
					if err != nil || width <= 0 || width > 8 {
						return 0, fmt.Errorf("invalid bitfield tag for field %s: %w", fieldStruct.Name, err)
					}

					// Ensure we have at least width bits available, loading bytes as needed.
					needed := width
					var val uint64
					var shift int
					for needed > 0 {
						if bitCount == 0 {
							if err := loadByte(); err != nil {
								return 0, err
							}
						}
						take := min(needed, bitCount)
						mask := byte((1 << take) - 1)
						part := (bitBuffer & mask)
						val |= uint64(part) << shift

						// consume bits from bitBuffer
						bitBuffer >>= take
						bitCount -= take
						needed -= take
						shift += take
					}

					// Set the value into the fieldValue.
					fv := fieldValue
					// If the field is addressable and settable (should be), set it.
					switch fv.Kind() {
					case reflect.Bool:
						fv.SetBool(val != 0)
					case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
						fv.SetUint(val)
					case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
						// preserve sign by interpreting as signed of width bits
						// If top bit of width is set, sign-extend.
						if width < int(64) {
							signBit := uint64(1 << (width - 1))
							if (val & signBit) != 0 {
								// sign extend
								mask := ^uint64(0) << width
								val |= mask
							}
						}
						fv.SetInt(int64(val)) // #nosec G115: decoding from binary; no overflow possible
					default:
						return 0, fmt.Errorf("unsupported bitfield type %s in field %s", fv.Kind(), fieldStruct.Name)
					}

					i++
					// update fieldValue for next iteration if there is one.
					if i < rt.NumField() {
						fieldValue = rv.Field(i)
					}
				}

				// After finishing the consecutive group, advance the data slice by the number of bytes consumed.
				consumed := bitIdx

				data = safeAdvance(data, consumed)

				continue
			}

			// Slice or array: iterate elements and unmarshal each.
			if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
				// For slices, we need to know the length. For arrays, it's fixed.
				// Also, we cannot infer length from slice type, so assume the
				// remaining bytes exactly match the serialized slice elements.
				// We'll treat slice length as full capacity by attempting to
				// fill until doMarshal would consume 0 bytes or error. However
				// this is ambiguous; since Marshal wrote all elements, caller
				// should initialize slice length beforehand.
				length := fieldValue.Len()
				for j := range length {
					consumed, err := doUnmarshal(data, fieldValue.Index(j), flags)
					if err != nil {
						return 0, fmt.Errorf("failed to unmarshal slice/array element %d of field %s: %w", j, fieldStruct.Name, err)
					}
					data = safeAdvance(data, consumed)
				}
				i++
				continue
			}

			// Nested struct: unmarshal recursively.
			if fieldValue.Kind() == reflect.Struct {
				consumed, err := doUnmarshal(data, fieldValue, flags)
				if err != nil {
					return 0, fmt.Errorf("failed to unmarshal nested struct %s: %w", fieldStruct.Name, err)
				}
				data = safeAdvance(data, consumed)
				i++
				continue
			}

			// Otherwise, basic numeric/bool types: read using binary.LittleEndian.
			switch fieldValue.Kind() {
			case reflect.Bool:
				if len(data) < 1 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for bool field %s", fieldStruct.Name)
					}
					fieldValue.SetBool(false)
					data = data[:0] // consume everything
				} else {
					fieldValue.SetBool(data[0] != 0)
					data = data[1:]
				}
			case reflect.Uint8:
				if len(data) < 1 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for uint8 field %s", fieldStruct.Name)
					}
					fieldValue.SetUint(0)
					data = data[:0] // consume everything
				} else {
					fieldValue.SetUint(uint64(data[0])) // #nosec G115: decoding from binary; no overflow possible
					data = data[1:]
				}
			case reflect.Int8:
				if len(data) < 1 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for int8 field %s", fieldStruct.Name)
					}
					fieldValue.SetInt(0)
					data = data[:0] // consume everything
				} else {
					fieldValue.SetInt(int64(int8(data[0]))) // #nosec G115: decoding from binary; no overflow possible
					data = data[1:]
				}
			case reflect.Uint16:
				if len(data) < 2 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for uint16 field %s", fieldStruct.Name)
					}
					fieldValue.SetUint(0)
					data = data[:0] // consume everything
				} else {
					v := binary.LittleEndian.Uint16(data)
					fieldValue.SetUint(uint64(v)) // #nosec G115: decoding from binary; no overflow possible
					data = data[2:]
				}
			case reflect.Int16:
				if len(data) < 2 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for int16 field %s", fieldStruct.Name)
					}
					fieldValue.SetInt(0)
					data = data[:0] // consume everything
				} else {
					v := int16(binary.LittleEndian.Uint16(data)) // #nosec G115: decoding from binary; no overflow possible
					fieldValue.SetInt(int64(v))
					data = data[2:]
				}
			case reflect.Uint32:
				if len(data) < 4 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for uint32 field %s", fieldStruct.Name)
					}
					fieldValue.SetUint(0)
					data = data[:0] // consume everything
				} else {
					v := binary.LittleEndian.Uint32(data) // #nosec G115: decoding from binary; no overflow possible
					fieldValue.SetUint(uint64(v))
					data = data[4:]
				}
			case reflect.Int32:
				if len(data) < 4 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for int32 field %s", fieldStruct.Name)
					}
					fieldValue.SetInt(0)
					data = data[:0] // consume everything
				} else {
					v := int32(binary.LittleEndian.Uint32(data)) // #nosec G115: decoding from binary; no overflow possible
					fieldValue.SetInt(int64(v))
					data = data[4:]
				}
			case reflect.Uint64, reflect.Uint:
				if len(data) < 8 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for uint64 field %s", fieldStruct.Name)
					}
					fieldValue.SetUint(0)
					data = data[:0] // consume everything
				} else {
					v := binary.LittleEndian.Uint64(data) // #nosec G115: decoding from binary; no overflow possible
					fieldValue.SetUint(v)
					data = data[8:]
				}
			case reflect.Int64, reflect.Int:
				if len(data) < 8 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for int64 field %s", fieldStruct.Name)
					}
					fieldValue.SetInt(0)
					data = data[:0] // consume everything
				} else {
					v := int64(binary.LittleEndian.Uint64(data)) // #nosec G115: decoding from binary; no overflow possible
					fieldValue.SetInt(v)
					data = data[8:]
				}
			case reflect.Float32:
				if len(data) < 4 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for float32 field %s", fieldStruct.Name)
					}
					fieldValue.SetFloat(0)
					data = data[:0] // consume everything
				} else {
					bits := binary.LittleEndian.Uint32(data)
					fieldValue.SetFloat(float64(math.Float32frombits(bits)))
					data = data[4:]
				}
			case reflect.Float64:
				if len(data) < 8 {
					if flags&FlagTreatEndAsZero == 0 {
						return 0, fmt.Errorf("not enough data for float64 field %s", fieldStruct.Name)
					}
					fieldValue.SetFloat(0)
					data = data[:0] // consume everything
				} else {
					bits := binary.LittleEndian.Uint64(data)
					fieldValue.SetFloat(math.Float64frombits(bits))
					data = data[8:]
				}
			default:
				return 0, ErrInvalidKind
			}
			i++
		}
	case reflect.Array:
		// iterate elements
		for idx := range rv.Len() {
			consumed, err := doUnmarshal(data, rv.Index(idx), flags)
			if err != nil {
				return 0, err
			}
			data = safeAdvance(data, consumed)
		}
	case reflect.Slice:
		// Clear slice
		rv.SetLen(0)

		sliceTarget := rv

		for idx := 0; len(data) > 0; idx++ {
			// Create a new zero value element to unmarshal into
			elem := reflect.New(rv.Type().Elem()).Elem()

			consumed, err := doUnmarshal(data, elem, flags)
			if err != nil {
				return 0, err
			}

			// Append element to sliceTarget
			sliceTarget = reflect.Append(sliceTarget, elem)
			// advance input
			data = safeAdvance(data, consumed)

			// If doUnmarshal consumed 0 bytes, avoid infinite loop: break.
			if consumed == 0 {
				break
			}
		}

		// set the slice value to the appended result
		rv.Set(sliceTarget)
	default:
		// basic scalar types: read directly
		switch rv.Kind() {
		case reflect.Bool:
			if len(data) < 1 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for bool")
				}
				rv.SetBool(false)
				data = data[:0] // consume everything
			} else {
				rv.SetBool(data[0] != 0)
				data = data[1:]
			}
		case reflect.Uint8:
			if len(data) < 1 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for uint8")
				}
				rv.SetUint(0)
				data = data[:0] // consume everything
			} else {
				rv.SetUint(uint64(data[0]))
				data = data[1:]
			}
		case reflect.Int8:
			if len(data) < 1 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for int8")
				}
				rv.SetInt(0)
				data = data[:0] // consume everything
			} else {
				rv.SetInt(int64(int8(data[0])))
				data = data[1:]
			}
		case reflect.Uint16:
			if len(data) < 2 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for uint16")
				}
				rv.SetUint(0)
				data = data[:0] // consume everything
			} else {
				v := binary.LittleEndian.Uint16(data) // #nosec G115: decoding from binary; no overflow possible
				rv.SetUint(uint64(v))
				data = data[2:]
			}
		case reflect.Int16:
			if len(data) < 2 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for int16")
				}
				rv.SetInt(0)
				data = data[:0] // consume everything
			} else {
				v := int16(binary.LittleEndian.Uint16(data)) // #nosec G115: decoding from binary; no overflow possible
				rv.SetInt(int64(v))
				data = data[2:]
			}
		case reflect.Uint32:
			if len(data) < 4 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for uint32")
				}
				rv.SetUint(0)
				data = data[:0] // consume everything
			} else {
				v := binary.LittleEndian.Uint32(data) // #nosec G115: decoding from binary; no overflow possible
				rv.SetUint(uint64(v))
				data = data[4:]
			}
		case reflect.Int32:
			if len(data) < 4 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for int32")
				}
				rv.SetInt(0)
				data = data[:0] // consume everything
			} else {
				v := int32(binary.LittleEndian.Uint32(data)) // #nosec G115: decoding from binary; no overflow possible
				rv.SetInt(int64(v))
				data = data[4:]
			}
		case reflect.Uint64, reflect.Uint:
			if len(data) < 8 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for uint64")
				}
				rv.SetUint(0)
				data = data[:0] // consume everything
			} else {
				v := binary.LittleEndian.Uint64(data)
				rv.SetUint(v)
				data = data[8:]
			}
		case reflect.Int64, reflect.Int:
			if len(data) < 8 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for int64")
				}
				rv.SetInt(0)
				data = data[:0] // consume everything
			} else {
				v := int64(binary.LittleEndian.Uint64(data)) // #nosec G115: decoding from binary; no overflow possible
				rv.SetInt(v)
				data = data[8:]
			}
		case reflect.Float32:
			if len(data) < 4 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for float32")
				}
				rv.SetFloat(0)
				data = data[:0] // consume everything
			} else {
				bits := binary.LittleEndian.Uint32(data)
				rv.SetFloat(float64(math.Float32frombits(bits)))
				data = data[4:]
			}
		case reflect.Float64:
			if len(data) < 8 {
				if flags&FlagTreatEndAsZero == 0 {
					return 0, errors.New("not enough data for float64")
				}
				rv.SetFloat(0)
				data = data[:0] // consume everything
			} else {
				bits := binary.LittleEndian.Uint64(data)
				rv.SetFloat(math.Float64frombits(bits))
				data = data[8:]
			}
		default:
			return 0, ErrInvalidKind
		}
	}

	consumed := origLen - len(data)
	return consumed, nil
}

func safeAdvance[T any](p []T, n int) []T {
	if n <= 0 {
		return p
	}
	if n >= len(p) {
		return p[:0]
	}
	return p[n:]
}
