package bitfield

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
	"slices"
	"strconv"
)

// MarshalToHexString serializes any struct containing bitfields
// (designated by tag "bitfield"), nested structs, and arrays
// to a hex string. Bitfields are packed in order (little-endian bit
// order, similar to C/C++ with one‐byte packing).
func MarshalToHexString(s any) (string, error) {
	data, err := doMarshal(reflect.ValueOf(s))
	return hex.EncodeToString(data), err
}

// Marshal serializes any struct containing bitfields
// (designated by tag "bitfield"), nested structs, and arrays.
// Bitfields are packed in order (little-endian bit order, similar to
// C/C++ with one‐byte packing).
func Marshal(s any) ([]byte, error) {
	return doMarshal(reflect.ValueOf(s))
}

// doMarshal recursively processes the reflect.Value (which should be a
// struct or an array) and returns the serialized bytes.
func doMarshal(rv reflect.Value) ([]byte, error) {
	var buf bytes.Buffer

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
					return nil, fmt.Errorf("bitfield:\"union\" on non-struct field %s", fieldStruct.Name)
				}

				// Marshal all members of the union struct, merging them.
				unionData := make([]byte, 0, 512)
				for j := range fieldValue.NumField() {
					member := fieldValue.Field(j)
					memberStruct := fieldValue.Type().Field(j)

					if memberStruct.Tag.Get("bitfield") == "-" {
						continue
					}

					// If member is itself a union, process it recursively as a union.
					if memberStruct.Tag.Get("bitfield") == "union" {
						if member.Kind() == reflect.Pointer {
							member = member.Elem()
						}
						if member.Kind() != reflect.Struct {
							return nil, fmt.Errorf("bitfield:\"union\" on non-struct member %s.%s",
								fieldStruct.Name, memberStruct.Name)
						}

						// Process nested union - merge its data.
						nestedUnionData := make([]byte, 0, 512)
						for k := range member.NumField() {
							nestedMember := member.Field(k)
							nestedMemberStruct := member.Type().Field(k)

							if nestedMemberStruct.Tag.Get("bitfield") == "-" {
								continue
							}

							data, err := doMarshal(nestedMember)
							if err != nil {
								return nil, fmt.Errorf("failed to marshal nested union member %s.%s.%s: %w",
									fieldStruct.Name, memberStruct.Name, nestedMemberStruct.Name, err)
							}

							if len(nestedUnionData) < len(data) {
								nestedUnionData = slices.Grow(nestedUnionData, len(data)-len(nestedUnionData))
								nestedUnionData = nestedUnionData[:len(data)]
							}

							for i := range data {
								nestedUnionData[i] |= data[i]
							}
						}

						if len(unionData) < len(nestedUnionData) {
							unionData = slices.Grow(unionData, len(nestedUnionData)-len(unionData))
							unionData = unionData[:len(nestedUnionData)]
						}

						for i := range nestedUnionData {
							unionData[i] |= nestedUnionData[i]
						}
						continue
					}

					data, err := doMarshal(member)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal union member %s.%s: %w",
							fieldStruct.Name, fieldValue.Type().Field(j).Name, err)
					}

					if len(unionData) < len(data) {
						unionData = slices.Grow(unionData, len(data)-len(unionData))
						unionData = unionData[:len(data)]
					}

					for i := range data {
						unionData[i] |= data[i]
					}
				}

				if _, err := buf.Write(unionData); err != nil {
					return nil, err
				}

				i++
				continue
			}

			// If the field is a bitfield, we may have several consecutive bitfields.
			if _, ok := fieldStruct.Tag.Lookup("bitfield"); ok {
				// Start accumulating bits.
				var bitBuffer byte // current accumulation
				var bitCount int   // number of bits accumulated in bitBuffer

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
						return nil, fmt.Errorf("invalid bitfield tag for field %s: %w", fieldStruct.Name, err)
					}

					// Get the field's unsigned numeric value.
					var fieldUint uint64
					switch fieldValue.Kind() {
					case reflect.Bool:
						if fieldValue.Bool() {
							fieldUint = 1
						} else {
							fieldUint = 0
						}
					case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
						fieldUint = fieldValue.Uint()
					case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
						fieldUint = uint64(fieldValue.Int()) // #nosec G115: encoding int64 as binary. No overflow is possible.
					default:
						return nil, fmt.Errorf("unsupported bitfield type %s in field %s", fieldValue.Kind(), fieldStruct.Name)
					}

					// Limit the value to the allowed bits.
					maxVal := uint64(1<<width) - 1
					fieldUint &= maxVal

					// Pack width bits into the current bitBuffer.
					bitsRemaining := width
					for bitsRemaining > 0 {
						available := 8 - bitCount // free bits in current byte.
						// if no free space flush current byte.
						if available == 0 {
							if err := buf.WriteByte(bitBuffer); err != nil {
								return nil, err
							}
							bitBuffer = 0
							bitCount = 0
							available = 8
						}

						// How many bits can we write from this field into current byte?
						toWrite := min(bitsRemaining, available)

						// Extract the bits to write from fieldUint.
						// (They come from the lower portion.)
						part := byte(fieldUint & ((1 << toWrite) - 1))
						// Put them into bitBuffer at the proper offset.
						bitBuffer |= part << bitCount

						// Update counters.
						bitCount += toWrite
						bitsRemaining -= toWrite
						fieldUint >>= toWrite
					}

					i++
					// If there are remaining fields, update fieldValue.
					if i < rt.NumField() {
						fieldValue = rv.Field(i)
					}
				}

				// Flush any leftover bits.
				if bitCount > 0 {
					if err := buf.WriteByte(bitBuffer); err != nil {
						return nil, err
					}
				}
				continue
			}

			// If the field is a slice/array, process each element.
			if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
				// Marshal each element recursively.
				for j := range fieldValue.Len() {
					elementBytes, err := doMarshal(fieldValue.Index(j))
					if err != nil {
						return nil, fmt.Errorf("failed to marshal slice/array element %d of field %s: %w", j, fieldStruct.Name, err)
					}
					if _, err := buf.Write(elementBytes); err != nil {
						return nil, err
					}
				}
				i++
				continue
			}

			// If the field is a nested struct, process recursively.
			if fieldValue.Kind() == reflect.Struct {
				nestedBytes, err := doMarshal(fieldValue)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal nested struct %s: %w", fieldStruct.Name, err)
				}
				_, err = buf.Write(nestedBytes)
				if err != nil {
					return nil, err
				}
				i++
				continue
			}

			// Otherwise, write non-bitfield, non-complex fields.
			fieldBytes, err := doMarshal(fieldValue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal field %s: %w", fieldStruct.Name, err)
			}
			if err := binary.Write(&buf, binary.LittleEndian, fieldBytes); err != nil {
				return nil, fmt.Errorf("failed to write field %s: %w", fieldStruct.Name, err)
			}
			i++
		}
	case reflect.Slice, reflect.Array:
		// For slices/arrays, process each element.
		for i := range rv.Len() {
			elBytes, err := doMarshal(rv.Index(i))
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(elBytes)
			if err != nil {
				return nil, err
			}
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
		// For basic types, write them out directly.
		if err := binary.Write(&buf, binary.LittleEndian, rv.Interface()); err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidKind
	}

	return buf.Bytes(), nil
}
