// Basic mathematics routines.
package utils

import (
	"math"
	"math/bits"
	"unsafe"
)

const (
	Float32Epsilon = 1.1920929e-07         // pre-computed math.Nextafter32(1, 2) - 1
	Float64Epsilon = 2.220446049250313e-16 // pre-computed math.Nextafter(1, 2) - 1
)

const (
	// MinFloat64SafeInt is the minimum exactly representable integer in float64.
	MinFloat64SafeInt = -(1 << 53)

	// MaxFloat64SafeInt the is maximum exactly representable integer in float64.
	MaxFloat64SafeInt = 1 << 53
)

const MaxUintptr = ^uintptr(0)

// SignedInt is a type constraint that permits all signed integer types.
type SignedInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInt is a type constraint that permits all unsigned integer types.
type UnsignedInt interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a type constraint that permits all signed or unsigned integer types.
type Integer interface {
	SignedInt | UnsignedInt
}

// Float is a type constraint for floating-point types.
type Float interface {
	~float32 | ~float64
}

// SignedNumber is a type constraint for any signed numeric type: floating-point
// or signed integer.
type SignedNumber interface {
	Float | SignedInt
}

// Number is a type constraint for all built-in numeric types: floats, signed
// ints, unsigned ints, and uintptr.
type Number interface {
	Float | SignedInt | UnsignedInt
}

// Sign returns +1 if x > 0, -1 if x < 0, or 0 if x == 0.
func Sign[T SignedNumber](x T) T {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

// Max returns the larger of a and b.
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller of a and b.
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Clamp restricts value to the inclusive range [minValue, maxValue].
// If minValue > maxValue, their order is swapped.
func Clamp[T Number](value, minValue, maxValue T) T {
	if minValue > maxValue {
		minValue, maxValue = maxValue, minValue
	}
	if value < minValue {
		return minValue
	} else if value > maxValue {
		return maxValue
	}
	return value
}

// Lerp performs linear interpolation between start and end by t in [0,1].
func Lerp[T Float](start, end, t T) T {
	return start + (end-start)*t
}

// Unlerp computes the normalized position of value in [minValue, maxValue].
// If maxValue == minValue, it returns NaN.
func Unlerp[T Float](minValue, maxValue, value T) T {
	if maxValue == minValue {
		return T(math.NaN())
	}
	return (value - minValue) / (maxValue - minValue)
}

// LerpRange maps value from the range [min0, max0] into the range [min1, max1].
func LerpRange[T Float](min0, max0, min1, max1, value T) T {
	t := Unlerp(min0, max0, value)
	return Lerp(min1, max1, t)
}

// Modulo returns a true modulus result in [0, b). b must be > 0.
func Modulo[T SignedInt](a, b T) T {
	if b <= 0 {
		panic("Modulo: b must be > 0")
	}
	// This is basically ((a % b) + b) % b
	// but more readable and avoiding overflow.
	r := a % b
	if r < 0 {
		r += b
	}
	return r
}

// WrapWithOffset returns (i + offset) wrapped into [0, n). n must be > 0.
func WrapWithOffset[T SignedInt](i, offset, n T) T {
	if n <= 0 {
		panic("WrapWithOffset: n must be > 0")
	}
	return Modulo(i+offset%n, n)
}

// WrapIndex returns a true modulus result in [0, b) even if a is negative.
// This effectively wraps array index.
// If length is zero, zero is returned.
func WrapIndex[T SignedInt](index, length T) T {
	if length == 0 {
		return 0
	}
	return Modulo(index, length)
}

// Step returns 1 if x >= threshold, otherwise 0.
func Step[T Number](x, threshold T) T {
	if x >= threshold {
		return 1
	}
	return 0
}

// Deg2Rad converts degrees to radians.
func Deg2Rad[T Float](degrees T) T {
	return degrees * (math.Pi / 180)
}

// Rad2Deg converts radians to degrees.
func Rad2Deg[T Float](radians T) T {
	return radians * (180 / math.Pi)
}

// Fract returns the fractional part of x. Sign is preserved.
func Fract(x float64) float64 {
	_, f := math.Modf(x)
	return f
}

// IsWholeFloat64 reports whether x is an integer value.
// It returns false for NaN and infinities.
func IsWholeFloat64(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0) && Fract(x) == 0
}

// CanSafelyConvertFloat64ToInt reports whether a floating
// point number can be safely converted to an integer
// without being truncated.
func CanSafelyConvertFloat64ToInt(x float64, minX, maxX float64) bool {
	return IsWholeFloat64(x) && x >= minX && x <= maxX
}

// Abs returns the absolute value of x. Works for signed ints and floats.
func Abs[T SignedInt | Float](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Gcd computes the greatest common divisor of a and b using the Euclidean
// algorithm. The result is always non-negative.
func Gcd[T SignedInt](a T, b T) T {
	if b == 0 {
		return a
	}

	x, y := a, b
	for y != 0 {
		x, y = y, x%y
	}

	return Abs(x)
}

// ReduceFrac returns the reduced fraction a/b by dividing
// by greatest common divisor.
func ReduceFrac[T SignedInt](a T, b T) (T, T) {
	k := Gcd(a, b)
	return a / k, b / k
}

// Det3x3 returns the determinant of the provided 3x3 matrix.
func Det3x3[T Float](a [3][3]T) T {
	return a[0][0]*a[1][1]*a[2][2] +
		a[0][1]*a[1][2]*a[2][0] +
		a[0][2]*a[1][0]*a[2][1] -
		a[2][0]*a[1][1]*a[0][2] -
		a[2][1]*a[1][2]*a[0][0] -
		a[2][2]*a[1][0]*a[0][1]
}

func MinNumValue[T Number](value T) T {
	switch v := any(value).(type) { //nolint:ineffassign
	case int:
		v = math.MinInt
		return T(v)
	case int8:
		v = math.MinInt8
		return T(v)
	case int16:
		v = math.MinInt16
		return T(v)
	case int32:
		v = math.MinInt32
		return T(v)
	case int64:
		v = math.MinInt64
		return T(v)
	case uint, uint8, uint16, uint32, uint64, uintptr:
		return 0
	case float32:
		v = -math.MaxFloat32
		return T(v)
	case float64:
		v = -math.MaxFloat64
		return T(v)
	default:
		panic("unreachable")
	}
}

func MaxNumValue[T Number](v T) T {
	switch v := any(v).(type) { //nolint:ineffassign
	case int:
		v = math.MaxInt
		return T(v)
	case int8:
		v = math.MaxInt8
		return T(v)
	case int16:
		v = math.MaxInt16
		return T(v)
	case int32:
		v = math.MaxInt32
		return T(v)
	case int64:
		v = math.MaxInt64
		return T(v)
	case uint:
		v = math.MaxUint
		return T(v)
	case uint8:
		v = math.MaxUint8
		return T(v)
	case uint16:
		v = math.MaxUint16
		return T(v)
	case uint32:
		v = math.MaxUint32
		return T(v)
	case uint64:
		v = math.MaxUint64
		return T(v)
	case uintptr:
		v = MaxUintptr
		return T(v)
	case float32:
		v = math.MaxFloat32
		return T(v)
	case float64:
		v = math.MaxFloat64
		return T(v)
	default:
		panic("unreachable")
	}
}

func CanSafelyDecrement[T Integer](v T) bool {
	switch v := any(v).(type) {
	case int:
		return v != math.MinInt
	case int8:
		return v != math.MinInt8
	case int16:
		return v != math.MinInt16
	case int32:
		return v != math.MinInt32
	case int64:
		return v != math.MinInt64
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case uintptr:
		return v != 0
	default:
		panic("unreachable")
	}
}

func CanSafelyIncrement[T Integer](v T) bool {
	switch v := any(v).(type) {
	case int:
		return v != math.MaxInt
	case int8:
		return v != math.MaxInt8
	case int16:
		return v != math.MaxInt16
	case int32:
		return v != math.MaxInt32
	case int64:
		return v != math.MaxInt64
	case uint:
		return v != math.MaxUint
	case uint8:
		return v != math.MaxUint8
	case uint16:
		return v != math.MaxUint16
	case uint32:
		return v != math.MaxUint32
	case uint64:
		return v != math.MaxUint64
	case uintptr:
		return v != MaxUintptr
	default:
		panic("unreachable")
	}
}

func SafelyDecrement[T Integer](v T) T {
	switch v := any(v).(type) {
	case int:
		if v != math.MinInt {
			v--
		}
		return T(v)
	case int8:
		if v != math.MinInt8 {
			v--
		}
		return T(v)
	case int16:
		if v != math.MinInt16 {
			v--
		}
		return T(v)
	case int32:
		if v != math.MinInt32 {
			v--
		}
		return T(v)
	case int64:
		if v != math.MinInt64 {
			v--
		}
		return T(v)
	case uint:
		if v != 0 {
			v--
		}
		return T(v)
	case uint8:
		if v != 0 {
			v--
		}
		return T(v)
	case uint16:
		if v != 0 {
			v--
		}
		return T(v)
	case uint32:
		if v != 0 {
			v--
		}
		return T(v)
	case uint64:
		if v != 0 {
			v--
		}
		return T(v)
	case uintptr:
		if v != 0 {
			v--
		}
		return T(v)
	default:
		panic("unreachable")
	}
}

func SafelyIncrement[T Integer](v T) T {
	switch v := any(v).(type) {
	case int:
		if v != math.MaxInt {
			v++
		}
		return T(v)
	case int8:
		if v != math.MaxInt8 {
			v++
		}
		return T(v)
	case int16:
		if v != math.MaxInt16 {
			v++
		}
		return T(v)
	case int32:
		if v != math.MaxInt32 {
			v++
		}
		return T(v)
	case int64:
		if v != math.MaxInt64 {
			v++
		}
		return T(v)
	case uint:
		if v != math.MaxUint {
			v++
		}
		return T(v)
	case uint8:
		if v != math.MaxUint8 {
			v++
		}
		return T(v)
	case uint16:
		if v != math.MaxUint16 {
			v++
		}
		return T(v)
	case uint32:
		if v != math.MaxUint32 {
			v++
		}
		return T(v)
	case uint64:
		if v != math.MaxUint64 {
			v++
		}
		return T(v)
	case uintptr:
		if v != MaxUintptr {
			v++
		}
		return T(v)
	default:
		panic("unreachable")
	}
}

// GreateGreaterl returns 1 if a > b; -1 if a < b; 0 if a = b.
// The function compiles to branchless code and is inline-friendly.
func CompareBool(a, b bool) int {
	ai := 0
	if a {
		ai = 1
	}
	bi := 0
	if b {
		bi = 1
	}
	return ai - bi
}

// Pow10TabUint64[i] = 10^i for uint64.
var Pow10TabUint64 = [20]uint64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
}

// Pow10Uint64 returns the power of ten as uint64.
// If the power is negative, 0 is returned,
// If the power results in the integer that cannot be
// represented by uint64, [math.MaxUint64] is returned.
func Pow10Uint64(pow int) uint64 {
	if pow < 0 {
		return 0
	}
	if pow >= len(Pow10TabUint64) {
		return math.MaxUint64
	}
	return Pow10TabUint64[pow]
}

// bitLenToNearPow10Tab[bits.Len64(i)] = max(1, decimalDigitCount(i) - 1) for uint64.
var bit64LenToNearPow10Tab = [64 + 1]byte{
	1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10,
	11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 14, 14, 14, 15,
	15, 15, 15, 16, 16, 16, 17, 17, 17, 18, 18, 18, 18, 19,
}

// DigitCountUint64 returns the number of decimal digits
// in the provided number. The digit count for zero is one.
// The computation time is constant.
func DigitCountUint64(x uint64) int {
	// Get the nearest exponent of ten based on the bit length
	// of the number.
	bits := bits.Len64(x)
	pow := int(bit64LenToNearPow10Tab[bits])

	// If the number is less than the power of ten,
	// return the nearest exponent. Otherwise,
	// this number is one digit greater than the
	// nearest exponent.
	//
	// If it were greater by more digits than one, it would snap
	// to the greater power of ten in the table, so such
	// a case is actually impossible. Each snap
	// grows by 2 whereas the exponent grows by 10,
	// i.e., much more frequent, so we won't miss the
	// numbers within the snaps.
	// Also, the exponent cannot be outside tables,
	// because they cover the whole possible range.
	//
	// In other words, this is basically this kind of
	// if-else chain:
	// - if x < 10^1: return 1
	// - if x < 10^2: return 2
	// - if x < 10^3: return 3
	// - etc.
	// but it is for any representable exponent and is table-driven.
	//
	// Let's illustrate this for small numbers:
	// - Let x=9 => bits=4 => pow=1, 9<10^1 => 1
	// - Let x=12 => bits=4 => pow=1, 12>=10^1 => 1 + 1 => 2
	// - Let x=63 => bits=6 => pow=1, 63>=10^1 => 1 + 1 => 2
	// - Let x=99 => bits=7 => pow=2, 99<10^2 => 2
	// - Let x=100 => bits=7 => pow=2, 100>=10^2 => 2 + 1 => 3
	if x < Pow10TabUint64[pow] {
		return pow
	}
	return pow + 1
}

// AbsToUint8 returns the absolute value of the signed value
// in unsigned type of the same bit width. Two's complement
// is considered so no overflow on the minimum possible values.
func AbsToUint8(x int8) uint8 {
	if x < 0 {
		// Special case for the minimum possible value.
		if x == math.MinInt8 {
			return math.MaxInt8 + 1
		}

		return uint8(-x)
	}

	return uint8(x)
}

// AbsToUint16 returns the absolute value of the signed value
// in unsigned type of the same bit width. Two's complement
// is considered so no overflow on the minimum possible values.
func AbsToUint16(x int16) uint16 {
	if x < 0 {
		// Special case for the minimum possible value.
		if x == math.MinInt16 {
			return math.MaxInt16 + 1
		}

		return uint16(-x)
	}

	return uint16(x)
}

// AbsToUint32 returns the absolute value of the signed value
// in unsigned type of the same bit width. Two's complement
// is considered so no overflow on the minimum possible values.
func AbsToUint32(x int32) uint32 {
	if x < 0 {
		// Special case for the minimum possible value.
		if x == math.MinInt32 {
			return math.MaxInt32 + 1
		}

		return uint32(-x)
	}

	return uint32(x)
}

// AbsToUint64 returns the absolute value of the signed value
// in unsigned type of the same bit width. Two's complement
// is considered so no overflow on the minimum possible values.
func AbsToUint64(x int64) uint64 {
	if x < 0 {
		// Special case for the minimum possible value.
		if x == math.MinInt64 {
			return math.MaxInt64 + 1
		}

		return uint64(-x)
	}

	return uint64(x)
}

// ClampUint8 converts x of int type to uint8
// ensuring that x is clamped.
// The result value is of uint8 type so no type cast is
// needed that could fire overflow detectors.
func ClampUint8(x int64) uint8 {
	switch {
	case x < 0:
		return 0
	case x > math.MaxUint8:
		return math.MaxUint8
	default:
		return uint8(x)
	}
}

// ClampInt8 converts x of int type to int8
// ensuring that x is clamped.
// The result value is of int8 type so no type cast is
// needed that could fire overflow detectors.
func ClampInt8(x int64) int8 {
	switch {
	case x < math.MinInt8:
		return math.MinInt8
	case x > math.MaxInt8:
		return math.MaxInt8
	default:
		return int8(x)
	}
}

// ClampUint16 converts x of int type to uint16
// ensuring that x is clamped.
// The result value is of uint16 type so no type cast is
// needed that could fire overflow detectors.
func ClampUint16(x int64) uint16 {
	switch {
	case x < 0:
		return 0
	case x > math.MaxUint16:
		return math.MaxUint16
	default:
		return uint16(x)
	}
}

// ClampInt16 converts x of int type to int16
// ensuring that x is clamped.
// The result value is of int16 type so no type cast is
// needed that could fire overflow detectors.
func ClampInt16(x int64) int16 {
	switch {
	case x < math.MinInt16:
		return math.MinInt16
	case x > math.MaxInt16:
		return math.MaxInt16
	default:
		return int16(x)
	}
}

// ClampUint32 converts x of int type to uint32
// ensuring that x is clamped.
// The result value is of uint32 type so no type cast is
// needed that could fire overflow detectors.
func ClampUint32(x int64) uint32 {
	switch {
	case x < 0:
		return 0
	case x > math.MaxUint32:
		return math.MaxUint32
	default:
		return uint32(x)
	}
}

// ClampInt32 converts x of int type to int32
// ensuring that x is clamped.
// The result value is of int32 type so no type cast is
// needed that could fire overflow detectors.
func ClampInt32(x int64) int32 {
	switch {
	case x < math.MinInt32:
		return math.MinInt32
	case x > math.MaxInt32:
		return math.MaxInt32
	default:
		return int32(x)
	}
}

// ClampPort converts port of int type to uint16
// ensuring that the port is clamped.
func ClampPort(port int64) uint16 {
	return ClampUint16(port)
}

// BitSizeOf returns the bit size of the given value.
func BitSizeOf[T any](x T) int {
	const bitsPerByte = 8
	return int(unsafe.Sizeof(x)) * bitsPerByte //nolint:gosec
}

// NearlyEqual reports whether a nearly equals b.
// The distance between a and b is less or equals to epsilon eps.
// If eps is zero, then the result is the same as a == b.
// This function is designed to be used with [Float32Epsilon] and
// [Float64Epsilon] epsilons but it is not necessary.
func NearlyEqual[T Number](a, b, eps T) bool {
	if a > b {
		return a-b <= eps
	}
	return b-a <= eps
}

// IntSinAngles is the total angles represented by [IntSin] and [IntCos]
// functions.
const IntSinAngles = 2048

// intSinTable is a sine table with [IntSinAngles] angles for higher precision.
// The table is all integers so the range is multiplied by 16383 insteaad of
// -1 to 1 range. The range of the table is -16383 to 16383.
// The size must be power of two for the bitwise-AND and division by two.
var intSinTable [IntSinAngles]int16

// InitIntSinTable initializes the look-up table for sine.
func InitIntSinTable() {
	for i := range intSinTable {
		intSinTable[i] = int16(math.Sin(float64(i)*math.Pi/(float64(len(intSinTable))/2)) * 16383)
	}
}

// IntSin returns an integer sine within [-16383..16383] range.
// The angle is within [0..[IntSinAngles]) range.
// Degrees can be converted this way: deg * [IntSinAngles] / 360.
// Both positive and negative angles are wrapped automatically.
// This function is intended for high performance integer-only computations
// involving trigonometric functions. It must be called only after
// [InitIntSinTable] is called, otherwise the returned value is always zero.
func IntSin(angle int) int16 {
	// NOTE: we must use bitwise AND here instead of modulo,
	// not only because of performance, but also because
	// AND correctly wraps negative indicies unlike modulo.
	// For example, -1 % 16, gives -1, whereas we expect 15;
	// -1 & 0xF gives 15, which is what we want.
	return intSinTable[angle&(len(intSinTable)-1)]
}

// IntCos is similar to [IntSin] but returns cosine instead of sine.
func IntCos(angle int) int16 {
	// The offset of cosine. We know that cos(a) = sin(a + 360/4 deg)
	// Saving the 1/4 turn offset proportion, it becomes:
	// cos(a) = sin(a + len(sintab)/4)
	const cosOff = len(intSinTable) / 4

	// NOTE: we must use bitwise AND here instead of modulo,
	// not only because of performance, but also because
	// AND correctly wraps negative indicies unlike modulo.
	// For example, -1 % 16, gives -1, whereas we expect 15;
	// -1 & 0xF gives 15, which is what we want.
	return intSinTable[(angle+cosOff)&(len(intSinTable)-1)]
}
