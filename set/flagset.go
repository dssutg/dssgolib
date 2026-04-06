package set

import (
	"fmt"
	"math/bits"
	"strings"
)

// FlagSet is a set of non-negative integers packed into a single integer value.
// Its zero value represents the empty set.
type (
	FlagSet8  uint8
	FlagSet16 uint16
	FlagSet32 uint32
	FlagSet64 uint64
)

// Has reports whether the set contains the non-negative value x.
func (s *FlagSet8) Has(x uint) bool { return *s&(1<<x) != 0 }

// Has reports whether the set contains the non-negative value x.
func (s *FlagSet16) Has(x uint) bool { return *s&(1<<x) != 0 }

// Has reports whether the set contains the non-negative value x.
func (s *FlagSet32) Has(x uint) bool { return *s&(1<<x) != 0 }

// Has reports whether the set contains the non-negative value x.
func (s *FlagSet64) Has(x uint) bool { return *s&(1<<x) != 0 }

// Add adds the non-negative value x to the set.
func (s *FlagSet8) Add(x uint) { *s |= 1 << x }

// Add adds the non-negative value x to the set.
func (s *FlagSet16) Add(x uint) { *s |= 1 << x }

// Add adds the non-negative value x to the set.
func (s *FlagSet32) Add(x uint) { *s |= 1 << x }

// Add adds the non-negative value x to the set.
func (s *FlagSet64) Add(x uint) { *s |= 1 << x }

// Remove removes the non-negative value from the set. If the value is not in
// the set, Remove is no-op.
func (s *FlagSet8) Remove(x uint) { *s &^= 1 << x }

// Remove removes the non-negative value from the set. If the value is not in
// the set, Remove is no-op.
func (s *FlagSet16) Remove(x uint) { *s &^= 1 << x }

// Remove removes the non-negative value from the set. If the value is not in
// the set, Remove is no-op.
func (s *FlagSet32) Remove(x uint) { *s &^= 1 << x }

// Remove removes the non-negative value from the set. If the value is not in
// the set, Remove is no-op.
func (s *FlagSet64) Remove(x uint) { *s &^= 1 << x }

// Union sets s to the union of s and t.
func (s *FlagSet8) Union(t *FlagSet8) { *s |= *t }

// Union sets s to the union of s and t.
func (s *FlagSet16) Union(t *FlagSet16) { *s |= *t }

// Union sets s to the union of s and t.
func (s *FlagSet32) Union(t *FlagSet32) { *s |= *t }

// Union sets s to the union of s and t.
func (s *FlagSet64) Union(t *FlagSet64) { *s |= *t }

// Sub sets s to the difference of s and t.
func (s *FlagSet8) Sub(t *FlagSet8) {
	x, _ := bits.Sub32(uint32(*s), uint32(*t), 0)
	*s = FlagSet8(x)
}

// Sub sets s to the difference of s and t.
func (s *FlagSet16) Sub(t *FlagSet16) {
	x, _ := bits.Sub32(uint32(*s), uint32(*t), 0)
	*s = FlagSet16(x)
}

// Sub sets s to the difference of s and t.
func (s *FlagSet32) Sub(t *FlagSet32) {
	x, _ := bits.Sub32(uint32(*s), uint32(*t), 0)
	*s = FlagSet32(x)
}

// Sub sets s to the difference of s and t.
func (s *FlagSet64) Sub(t *FlagSet64) {
	x, _ := bits.Sub64(uint64(*s), uint64(*t), 0)
	*s = FlagSet64(x)
}

// Len return the number of elements in the set.
func (s *FlagSet8) Len() int { return bits.OnesCount8(uint8(*s)) }

// Len return the number of elements in the set.
func (s *FlagSet16) Len() int { return bits.OnesCount16(uint16(*s)) }

// Len return the number of elements in the set.
func (s *FlagSet32) Len() int { return bits.OnesCount32(uint32(*s)) }

// Len return the number of elements in the set.
func (s *FlagSet64) Len() int { return bits.OnesCount64(uint64(*s)) }

// Clear clears the set.
func (s *FlagSet8) Clear() { *s = 0 }

// Clear clears the set.
func (s *FlagSet16) Clear() { *s = 0 }

// Clear clears the set.
func (s *FlagSet32) Clear() { *s = 0 }

// Clear clears the set.
func (s *FlagSet64) Clear() { *s = 0 }

// String returns the set as a string of the form "{1 2 3}".
func (s *FlagSet8) String() string { return toString(uint64(*s), 8) }

// String returns the set as a string of the form "{1 2 3}".
func (s *FlagSet16) String() string { return toString(uint64(*s), 16) }

// String returns the set as a string of the form "{1 2 3}".
func (s *FlagSet32) String() string { return toString(uint64(*s), 32) }

// String returns the set as a string of the form "{1 2 3}".
func (s *FlagSet64) String() string { return toString(uint64(*s), 64) }

func toString(reg uint64, width int) string {
	var sb strings.Builder

	sb.WriteByte('{')

	if reg != 0 {
		for i := range width {
			if reg&(1<<i) == 0 {
				continue
			}
			if sb.Len() > len("{") {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", i)
		}
	}

	sb.WriteByte('}')

	return sb.String()
}
