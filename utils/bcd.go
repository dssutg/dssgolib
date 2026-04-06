package utils

import "strings"

// StringToBCDSliceBE converts the decimal digits in s into packed
// Binary Coded Digits (BCD) stored in arr that is assumed to have enough length.
// Each digit is encoded sequentially as 4-bit BCD nibble in Big Endian (first digit
// in high nibble, second - low). Non-digits are ignored. The BCD stream is
// terminated by writing the low 4 bits of end as the final nibble. The terminator
// should be greater than 9.
func StringToBCDSliceBE(arr []byte, s string, end byte) {
	b := MakeNibbleArrayBuilder(arr)

	for i := range len(s) {
		c := s[i]
		if c < '0' || c > '9' {
			continue // ignore non-digit
		}
		b.WriteNibbleBE(c - '0')
	}

	b.WriteNibbleBE(end & 0xF) // terminate BCD
}

// StringToBCDSliceLE converts the decimal digits in s into packed
// Binary Coded Digits (BCD) stored in arr that is assumed to have enough length.
// Each digit is encoded sequentially as 4-bit BCD nibble in Little Endian (first digit
// in low nibble, second - high). Non-digits are ignored. The BCD stream is
// terminated by writing the low 4 bits of end as the final nibble. The terminator
// should be greater than 9.
func StringToBCDSliceLE(arr []byte, s string, end byte) {
	b := MakeNibbleArrayBuilder(arr)

	for i := range len(s) {
		c := s[i]
		if c < '0' || c > '9' {
			continue // ignore non-digit
		}
		b.WriteNibbleLE(c - '0')
	}

	b.WriteNibbleLE(end & 0xF) // terminate BCD
}

// BCDToStringBE returns the string representation of the Binary Encoded
// Digit (BCD) byte stream.  The function stops decoding once it encounters
// the first nibble that is greater than 9.  The order of nibbles is Big
// Endian (first digit in MSB nibble, second - LSB).
func BCDToStringBE(bcd []byte) string {
	var sb strings.Builder
	sb.Grow(len(bcd) * 2)

	for _, d := range bcd {
		high, low := d>>4, d&0xF

		// Decode high.
		if high > 9 {
			break
		}
		sb.WriteByte(high + '0')

		// Decode low.
		if low > 9 {
			break
		}
		sb.WriteByte(low + '0')
	}

	return sb.String()
}

// BCDToBytesBE is the same as [BCDToStringBE] but returns byte slice
// instead of string.
func BCDToBytesBE(bcd []byte) []byte {
	out := make([]byte, 0, len(bcd)*2)

	for _, d := range bcd {
		high, low := d>>4, d&0xF

		// Decode high.
		if high > 9 {
			break
		}
		out = append(out, high+'0')

		// Decode low.
		if low > 9 {
			break
		}
		out = append(out, low+'0')
	}

	return out
}

// BCDToStringBE returns the string representation of the Binary Encoded
// Digit (BCD) byte stream.  The function stops decoding once it encounters
// the first nibble that is greater than 9.  The order of nibbles is Little
// Endian (first digit in LSB nibble, second - MSB).
func BCDToStringLE(bcd []byte) string {
	var sb strings.Builder
	sb.Grow(len(bcd) * 2)

	for _, d := range bcd {
		low, high := d>>4, d&0xF

		// Decode high.
		if high > 9 {
			break
		}
		sb.WriteByte(high + '0')

		// Decode low.
		if low > 9 {
			break
		}
		sb.WriteByte(low + '0')
	}

	return sb.String()
}

// BCDToBytesLE is the same as [BCDToStringLE] but returns byte slice
// instead of string.
func BCDToBytesLE(bcd []byte) []byte {
	out := make([]byte, 0, len(bcd)*2)

	for _, d := range bcd {
		low, high := d>>4, d&0xF

		// Decode high.
		if high > 9 {
			break
		}
		out = append(out, high+'0')

		// Decode low.
		if low > 9 {
			break
		}
		out = append(out, low+'0')
	}

	return out
}
