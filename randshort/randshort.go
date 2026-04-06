// randshort is a simple predictable pseudo-random int16 number generator.
package randshort

var next = int64(0)

func Seed(seed int64) {
	next = seed
}

func Get() int16 {
	next *= 1103515245
	next += 12345

	x := next >> 16
	x &= 0x7fff

	return int16(x)
}
