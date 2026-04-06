package set

import (
	"fmt"
	"math/bits"
	"strings"
)

// BitSet is a set of non-negative integers.
// Its zero value represents the empty set.
type BitSet struct {
	words []uint64
}

// pos returns the position inside the bit set of the provided value.
// It panics if x is negative.
func (s *BitSet) pos(x int) (word int, bit uint) {
	if x < 0 {
		panic("x cannot be negative")
	}
	return x / 64, uint(x % 64)
}

// Has reports whether the set contains the non-negative value x.
// It panics if x is negative.
func (s *BitSet) Has(x int) bool {
	word, bit := s.pos(x)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// AddOne adds the non-negative value x to the set.
// It panics if x is negative.
func (s *BitSet) AddOne(x int) {
	word, bit := s.pos(x)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// Add adds the non-negative values stored in xv to the set.
// It panics if any of the values is negative.
func (s *BitSet) Add(x int, rest ...int) {
	s.AddOne(x)
	for x := range rest {
		s.AddOne(x)
	}
}

// Remove removes the non-negative value from the set.
// It panics if x is negative. If the value is not in the set,
// Remove is no-op. Backing memory is not freed or shrunk.
func (s *BitSet) Remove(x int) {
	word, bit := s.pos(x)
	if word < len(s.words) {
		s.words[word] &^= 1 << bit
	}
}

// Union sets s to the union of s and t.
func (s *BitSet) Union(t *BitSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Len return the number of elements in the set.
func (s *BitSet) Len() (count int) {
	for _, word := range s.words {
		count += bits.OnesCount64(word)
	}
	return count
}

// Clear clears the set without freeing the backing memory.
func (s *BitSet) Clear() {
	clear(s.words)
}

// String returns the set as a string of the form "{1 2 3}".
func (s *BitSet) String() string {
	var sb strings.Builder

	sb.WriteByte('{')

	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := range 64 {
			if word&(1<<j) == 0 {
				continue
			}
			if sb.Len() > len("{") {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", j+i*64)
		}
	}

	sb.WriteByte('}')

	return sb.String()
}
