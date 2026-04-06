package crc

import "github.com/dssutg/dssgolib/digest"

// ComputeCRC16ByTable calculates the CRC-16 checksum of the given byte slice.
// It uses the provided precomputed lookup table to update the checksum for each
// byte. The algorithm works by iterating over each byte in the slice and
// performing a bitwise operation that combines the current CRC value with the
// current byte. The final computed CRC is then returned as a 16-bit unsigned
// integer.
func ComputeCRC16ByTable(crc16Table [256]uint16, bytes []byte) uint16 {
	var crc uint16
	for i := range bytes {
		crc = (crc >> 8) ^ crc16Table[(crc^uint16(bytes[i]))&0xFF]
	}
	return crc
}

// ComputeCRC16ByTableHexString computes the CRC-16 checksum of the given byte
// slice. It uses the provided precomputed lookup table to update the checksum
// for each byte. The algorithm works by iterating over each byte in the slice
// and performing a bitwise operation that combines the current CRC value
// with the current byte. The 16-bit checksum is then stringified using
// the [digest.Uint16ToHex] function to represent that sum as a lowercase
// 4-character hex string.
func ComputeCRC16ByTableHexString(crc16Table [256]uint16, bytes []byte) string {
	return digest.Uint16ToHex(ComputeCRC16ByTable(crc16Table, bytes))
}
