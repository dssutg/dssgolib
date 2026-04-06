// Slice manipulation routines.
package utils

import (
	"encoding/json/v2"
	"fmt"
	"math"
	"strings"
)

// GetSliceWithUniqueElements returns a new slice containing only the unique
// elements from the provided slice. It preserves the order in which the
// elements first appeared in the original slice.
func GetSliceWithUniqueElements[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	var unique []T
	for _, value := range slice {
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			unique = append(unique, value)
		}
	}
	return unique
}

// TrimSliceToMax ensures that the given slice has at most maxLength elements.
// If the slice length exceeds maxLength, it returns a new slice that only
// contains the first maxLength elements. Otherwise, it returns the original
// slice.
func TrimSliceToMax[T any](s []T, maxLength int) []T {
	if len(s) > maxLength {
		// Creating a new slice that contains only the first maxLength elements.
		return s[:maxLength]
	}
	return s
}

// SqueezeSlice returns a slice holding only the non-zero value elements. The
// underlying array is modified during the call. The rest of the underlying array
// is filled with zero value. The complexity of this function is O(n).
func SqueezeSlice[T comparable](v []T) []T {
	var zero T
	i := 0
	for j, e := range v {
		if e == zero {
			continue
		}
		v[i] = e
		if j != i {
			// Fill with zero value to make GC collect it.
			v[j] = zero
		}
		i++
	}
	return v[:i]
}

// SliceEvery returns true if f returns true for all elements
// of the array. The function always returns true for empty arrays.
func SliceEvery[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, x := range s {
		if !f(x) {
			return false
		}
	}
	return true
}

// SliceGetElemOrZero returns the element at the slice index or
// zero value if the index is outside the slice boundaries.
func SliceGetElemOrZero[V any](slice []V, index int) V {
	if index < 0 || index >= len(slice) {
		var zero V
		return zero
	}
	return slice[index]
}

// MapSlice returns the new slice whose each element
// transformed with the mapper function.
func MapSlice[V, R any](slice []V, mapper func(x V) R) []R {
	result := make([]R, len(slice))
	for i := range slice {
		result[i] = mapper(slice[i])
	}
	return result
}

// MapIntSliceToString returns the new slice whose each element
// is converted to decimal int string.
func MapIntSliceToString[V Integer](slice []V) []string {
	result := make([]string, len(slice))
	for i := range slice {
		result[i] = Itoa(slice[i])
	}
	return result
}

// MapToStringSlice returns the new slice whose each element
// is converted to string as defined by [ToString].
func MapToStringSlice[V any](slice []V) []string {
	return MapSlice(slice, ToString)
}

// JoinSliceToString converts all the slice elements to strings as defined by [fmt.Print],
// and joins them with the delim.
func JoinSliceToString[V any](slice []V, delim string) string {
	// Fast path for empty slice.
	if len(slice) == 0 {
		return ""
	}

	var sb strings.Builder
	const avg = 16 // heuristic bytes per element
	sb.Grow(len(slice)*avg + len(delim)*(len(slice)-1))

	for i := range slice {
		if i > 0 {
			sb.WriteString(delim)
		}
		fmt.Fprint(&sb, slice[i])
	}

	return sb.String()
}

// Uint8Slice is []byte but marshals/unmarshals as a JSON array of numbers
// instead of Base64.
type Uint8Slice []byte

// MarshalJSON implements json.Marshaler.
// It encodes the slice as a JSON array of numbers (e.g. [1,2,3]) instead of Base64.
func (a Uint8Slice) MarshalJSON() ([]byte, error) {
	if a == nil {
		// Marshal nil slice to either null or [] depending on
		// whether it is jsonv2 or the legacy JSON API.
		return json.Marshal([]uint16(nil))
	}

	// The closest to the byte array without Base64 is uint16 array,
	// so first copy to it, then marshal it.
	tmp := make([]uint16, len(a))
	for i := range a {
		tmp[i] = uint16(a[i])
	}

	return json.Marshal(tmp)
}

// UnmarshalJSON implements json.Unmarshaler.
// It encodes the slice as a JSON array of numbers (e.g. [1,2,3]) instead of Base64.
// It complies with jsonv2 so nil slice is marshaled as empty array instead of null.
func (a *Uint8Slice) UnmarshalJSON(p []byte) error {
	if a == nil {
		return nil
	}

	// The closest to the byte array without Base64 is uint16 array,
	// so first unmarshal to it.
	var tmp []uint16
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	*a = nil // clear old value

	// Map and validate the temporary uint16 array to the destination array.
	dst := make([]byte, len(tmp))
	for i, x := range tmp {
		if x > math.MaxUint8 {
			return fmt.Errorf("element #%d (value %d) is out of byte range", i, x)
		}
		dst[i] = byte(x)
	}

	*a = dst // set the new array

	return nil
}

// FindOneMissingInt returns the missing integer in a sorted, unique ascending slice.
// Assumes exactly one element is missing from what would otherwise be a consecutive
// sequence. If no missing element is detected, returns ok=false.
func FindOneMissingInt(v []int) (missing int, ok bool) {
	// If sequence is too short to have anything missing.
	if len(v) < 2 {
		return 0, false
	}

	// If missing at the boundaries.
	if v[1]-v[0] != 1 {
		return v[0] + 1, true
	}
	if v[len(v)-1]-v[len(v)-2] != 1 {
		return v[len(v)-2] + 1, true
	}

	// Binary search the missing slot.
	lo, hi := 0, len(v)-1
	for lo <= hi {
		mid := int(uint(lo+hi) >> 1) // avoid overflow when computing mid
		want := v[0] + mid           // wanted value at index mid if sequence were consecutive from a[0]
		switch {
		case v[mid] == want:
			lo = mid + 1 // left side has no missing; move right
		case mid == 0, v[mid-1] == v[0]+(mid-1):
			// Missing is at or left of mid.
			// Check if this is 1st index where mismatch occurs.
			return want, true
		default:
			hi = mid - 1 // right side has no missing; move left
		}
	}

	return 0, false
}

// MoveInSliceTo moves slice s element at index src to dst.
// If src is out of range, this function is no-op.
// If dest <= 0, it is interpreted as 0, if more than slice length,
// it is interpreted as the last element index.
func MoveInSliceTo[T any](s []T, dst, src int) {
	n := len(s)
	if src < 0 || src >= n || n <= 1 {
		return // index out of range or nothing to move
	}

	// Clamp the destination index.
	switch {
	case dst < 0:
		dst = 0
	case dst >= n:
		dst = n - 1
	}

	shiftSlice(s, dst, src)
}

// MoveInSlice moves element at index src by offset positions within the slice s in-place.
// The destination index is clampped. Negative offset moves left, positive moves right.
func MoveInSlice[T any](s []T, src, offset int) {
	n := len(s)
	if src < 0 || src >= n || offset == 0 || n <= 1 {
		return // index out of range or nothing to move
	}

	offset %= n // reduce offset

	dst := src + offset // find destination index

	// Clamp the destination index.
	switch {
	case dst < 0:
		dst = 0
	case dst >= n:
		dst = n - 1
	}

	shiftSlice(s, dst, src)
}

// MoveInSliceCircular moves element at index src by offset positions within the slice s in-place,
// wrapping around circularly. Negative offset moves left, positive moves right.
func MoveInSliceCircular[T any](s []T, src, offset int) {
	n := len(s)
	if src < 0 || src >= n || offset == 0 || n <= 1 {
		return // index out of range or nothing to move
	}

	offset %= n // reduce offset

	// Find the destination index.
	dst := (src + offset) % n
	if dst < 0 {
		dst += n
	}

	shiftSlice(s, dst, src)
}

// shiftSlice moves the slice element at src index to dst index.
// No safety checks are performed.
func shiftSlice[T any](s []T, dst, src int) {
	if dst == src {
		return // nothing to move
	}

	val := s[src] // save element to move

	if dst < src { // move right to free the spot for the saved element
		copy(s[dst+1:src+1], s[dst:src])
	} else {
		copy(s[src:dst], s[src+1:dst+1]) // move left to free the spot for the saved element
	}

	s[dst] = val // put the saved element to the freed spot
}
