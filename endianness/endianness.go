package endianness

import (
	"unsafe"
)

func IsLittleEndian() bool {
	x := uint16(1)
	return *(*byte)(unsafe.Pointer(&x)) != 0
}

func IsBigEndian() bool {
	x := uint16(1)
	return *(*byte)(unsafe.Pointer(&x)) == 0
}

func EnsureLittleEndian() {
	if !IsLittleEndian() {
		panic("only Little Endian byte order is supported")
	}
}
