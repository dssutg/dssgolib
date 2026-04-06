package txt

import (
	"strconv"
	"unicode/utf8"
)

type Buf []byte

func NewBuf(initCap int) Buf {
	return make([]byte, 0, initCap)
}

func (b *Buf) Clear() {
	*b = (*b)[:0]
}

func (b *Buf) AddString(s string) {
	*b = append(*b, s...)
}

func (b *Buf) AddBool(t bool) {
	if t {
		b.AddString("true")
	} else {
		b.AddString("false")
	}
}

func (b *Buf) AddSRune(r rune) {
	var tmp [utf8.UTFMax]byte
	n := utf8.EncodeRune(tmp[:], r)
	*b = append(*b, tmp[:n]...)
}

func (b *Buf) AddSByte(c byte) {
	*b = append(*b, c)
}

func (b *Buf) AddInt(x int, base int) {
	*b = strconv.AppendInt(*b, int64(x), base)
}

func (b *Buf) AddInt8(x int8, base int) {
	*b = strconv.AppendInt(*b, int64(x), base)
}

func (b *Buf) AddInt16(x int16, base int) {
	*b = strconv.AppendInt(*b, int64(x), base)
}

func (b *Buf) AddInt32(x int32, base int) {
	*b = strconv.AppendInt(*b, int64(x), base)
}

func (b *Buf) AddInt64(x int64, base int) {
	*b = strconv.AppendInt(*b, x, base)
}

func (b *Buf) AddUint(x uint, base int) {
	*b = strconv.AppendUint(*b, uint64(x), base)
}

func (b *Buf) AddUint8(x uint8, base int) {
	*b = strconv.AppendUint(*b, uint64(x), base)
}

func (b *Buf) AddUint16(x uint16, base int) {
	*b = strconv.AppendUint(*b, uint64(x), base)
}

func (b *Buf) AddUint32(x uint32, base int) {
	*b = strconv.AppendUint(*b, uint64(x), base)
}

func (b *Buf) AddUint64(x uint64, base int) {
	*b = strconv.AppendUint(*b, x, base)
}

func (b *Buf) AddFloat32(f float64, fmt byte, prec int) {
	*b = strconv.AppendFloat(*b, f, fmt, prec, 32)
}

func (b *Buf) AddFloat64(f float64, fmt byte, prec int) {
	*b = strconv.AppendFloat(*b, f, fmt, prec, 64)
}
