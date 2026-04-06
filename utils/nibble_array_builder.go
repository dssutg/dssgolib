package utils

// NibbleArrayBuilder is a builder
// to conveniently construct an array of nibbles.
type NibbleArrayBuilder struct {
	Arr    []byte // underlying array
	Idx    int    // current index in the array
	Second bool   // whether the next nibble is second or first in the byte
}

// MakeNibbleArrayBuilder returns a properly initialized
// nibble array builder.
func MakeNibbleArrayBuilder(arr []byte) NibbleArrayBuilder {
	return NibbleArrayBuilder{Arr: arr}
}

// WriteNibbleBE appends a nibble to array. The order is Big Endian, i.e.,
// the first nibble is in the most significant 4 bits of byte, the
// second is in the least significant ones.
func (b *NibbleArrayBuilder) WriteNibbleBE(n byte) {
	n &= 0xF // ensure valid range

	// First nibble in byte.
	if !b.Second {
		b.Arr[b.Idx] = n << 4
		b.Second = true
		return
	}
	// Second nibble in byte.
	b.Arr[b.Idx] |= n
	b.Idx++ // advance to next byte
	b.Second = false
}

// WriteNibbleLE appends a nibble to array. The order is Little Endian, i.e.,
// the first nibble is in the least significant 4 bits of byte, the
// second is in the most significant ones.
func (b *NibbleArrayBuilder) WriteNibbleLE(n byte) {
	n &= 0xF // ensure valid range

	// First nibble in byte.
	if !b.Second {
		b.Arr[b.Idx] = n
		b.Second = true
		return
	}
	// Second nibble in byte.
	b.Arr[b.Idx] |= n << 4
	b.Idx++ // advance to next byte
	b.Second = false
}
