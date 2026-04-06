package bitfield

import (
	"fmt"
	"reflect"
	"strconv"
)

// MustSizeof calls [Sizeof] but panics on error.
func MustSizeof(s any) int {
	size, err := Sizeof(s)
	if err != nil {
		panic(fmt.Sprintf("MustSizeof: %v", err))
	}
	return size
}

// Sizeof returns the number of bytes occupied by the value s
// marshaled into bytes by using [Marshal].
func Sizeof(s any) (int, error) {
	return doSizeof(reflect.ValueOf(s))
}

// doSizeof recursively processes the reflect.Value (which should be a
// struct or an array) and returns the number of serialized bytes.
func doSizeof(rv reflect.Value) (int, error) {
	byteCount := 0

	// If our value is a pointer, dereference it.
	if rv.Kind() == reflect.Pointer {
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

			// If the field is a bitfield union.
			if fieldStruct.Tag.Get("bitfield") == "union" {
				// We expect a nested struct
				if fieldValue.Kind() == reflect.Pointer {
					fieldValue = fieldValue.Elem()
				}
				if fieldValue.Kind() != reflect.Struct {
					return 0, fmt.Errorf("bitfield:\"union\" on non-struct field %s", fieldStruct.Name)
				}

				// Marshal all members of the union struct, merging them.
				unionDataSize := 0
				for j := range fieldValue.NumField() {
					member := fieldValue.Field(j)
					memberStruct := fieldValue.Type().Field(j)

					if memberStruct.Tag.Get("bitfield") == "-" {
						continue
					}

					// If member is itself a union, process it recursively.
					if memberStruct.Tag.Get("bitfield") == "union" {
						if member.Kind() == reflect.Pointer {
							member = member.Elem()
						}
						if member.Kind() != reflect.Struct {
							return 0, fmt.Errorf("bitfield:\"union\" on non-struct member %s.%s",
								fieldStruct.Name, memberStruct.Name)
						}

						nestedUnionSize := 0
						for k := range member.NumField() {
							nestedMember := member.Field(k)
							nestedMemberStruct := member.Type().Field(k)

							if nestedMemberStruct.Tag.Get("bitfield") == "-" {
								continue
							}

							dataSize, err := doSizeof(nestedMember)
							if err != nil {
								return 0, fmt.Errorf("failed to sizeof nested union member %s.%s.%s: %w",
									fieldStruct.Name, memberStruct.Name, nestedMemberStruct.Name, err)
							}

							if nestedUnionSize < dataSize {
								nestedUnionSize = dataSize
							}
						}

						if unionDataSize < nestedUnionSize {
							unionDataSize = nestedUnionSize
						}
						continue
					}

					dataSize, err := doSizeof(member)
					if err != nil {
						return 0, fmt.Errorf("failed to marshal union member %s.%s: %w",
							fieldStruct.Name, fieldValue.Type().Field(j).Name, err)
					}

					if unionDataSize < dataSize {
						unionDataSize = dataSize
					}
				}

				byteCount += unionDataSize

				i++
				continue
			}

			// If the field is a bitfield, we may have several consecutive bitfields.
			if _, ok := fieldStruct.Tag.Lookup("bitfield"); ok {
				// Start accumulating bits.
				var bitCount int // number of bits accumulated in bitBuffer

				// Process consecutive fields that are tagged as bitfields.
				for i < rt.NumField() {
					fieldStruct = rt.Field(i)
					tagStr, isBitField := fieldStruct.Tag.Lookup("bitfield")
					if !isBitField {
						break
					}
					if tagStr == "-" {
						i++
						continue // ignore field
					}
					// Parse the bit width.
					width, err := strconv.Atoi(tagStr)
					if err != nil || width <= 0 || width > 8 {
						return 0, fmt.Errorf("invalid bitfield tag for field %s: %w", fieldStruct.Name, err)
					}

					// Get the field's unsigned numeric value.
					switch fieldValue.Kind() {
					case
						reflect.Bool,
						reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64,
						reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
					default:
						return 0, fmt.Errorf("unsupported bitfield type %s in field %s", fieldValue.Kind(), fieldStruct.Name)
					}

					// Pack width bits into the current bitBuffer.
					bitsRemaining := width
					for bitsRemaining > 0 {
						available := 8 - bitCount // free bits in current byte.
						// if no free space flush current byte.
						if available == 0 {
							byteCount++
							bitCount = 0
							available = 8
						}

						// How many bits can we write from this field into current byte?
						toWrite := min(bitsRemaining, available)

						// Update counters.
						bitCount += toWrite
						bitsRemaining -= toWrite
					}

					i++
					// If there are remaining fields, update fieldValue.
					if i < rt.NumField() {
						fieldValue = rv.Field(i)
					}
				}

				// Flush any leftover bits.
				if bitCount > 0 {
					byteCount++
				}
				continue
			}

			// If the field is a slice/array, process each element.
			if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
				// Marshal each element recursively.
				for j := range fieldValue.Len() {
					elementBytes, err := doSizeof(fieldValue.Index(j))
					if err != nil {
						return 0, fmt.Errorf("failed to marshal slice/array element %d of field %s: %w", j, fieldStruct.Name, err)
					}
					byteCount += elementBytes
				}
				i++
				continue
			}

			// If the field is a nested struct, process recursively.
			if fieldValue.Kind() == reflect.Struct {
				nestedBytes, err := doSizeof(fieldValue)
				if err != nil {
					return 0, fmt.Errorf("failed to marshal nested struct %s: %w", fieldStruct.Name, err)
				}
				byteCount += nestedBytes
				i++
				continue
			}

			// Otherwise, write non-bitfield, non-complex fields.
			fieldBytes, err := doSizeof(fieldValue)
			if err != nil {
				return 0, fmt.Errorf("failed to marshal field %s: %w", fieldStruct.Name, err)
			}
			byteCount += fieldBytes
			i++
		}
	case reflect.Slice, reflect.Array:
		// For slices/arrays, process each element.
		for i := range rv.Len() {
			elBytes, err := doSizeof(rv.Index(i))
			if err != nil {
				return 0, err
			}
			byteCount += elBytes
		}
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		// For basic types, get their size directly.
		byteCount += int(rv.Type().Size()) //nolint:gosec
	default:
		return 0, ErrInvalidKind
	}

	return byteCount, nil
}
