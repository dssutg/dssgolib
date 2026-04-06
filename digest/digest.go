// Package digest converts raw byte representation of data to their textual
// representation.
package digest

// Uint16ToHex converts a 16-bit unsigned integer into a hexadecimal string. The
// function produces a 4-character string representation of the number using
// lowercase hexadecimal digits.
func Uint16ToHex(x uint16) string {
	const hexDigits = "0123456789abcdef"

	hex := []byte{
		hexDigits[(x>>12)&0xF],
		hexDigits[(x>>8)&0xF],
		hexDigits[(x>>4)&0xF],
		hexDigits[x&0xF],
	}

	return string(hex)
}
