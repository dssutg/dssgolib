// Package set implements simple Set collection.
package set

import (
	"bytes"
	"cmp"
	"encoding/gob"
	"encoding/json/v2"
	"iter"
	"maps"
	"slices"
)

// Set is a generic set data structure that holds unique elements of any
// comparable type. It uses a map to store elements, where the key is the
// element and the value is an empty struct. Zero value is a ready to use
// set.
type Set[T comparable] struct {
	// elements is the underlying map that stores the set elements.
	// The key is the element and the value is an empty struct to save memory.
	elements map[T]struct{}
}

// Make is a convenience function that creates a new Set instance (by value)
// with its underlying map properly initialized.
func Make[T comparable]() Set[T] {
	return Set[T]{elements: make(map[T]struct{})}
}

// Clone returns the shallow copy of the set.
func (s *Set[T]) Clone() Set[T] {
	return Set[T]{elements: maps.Clone(s.elements)}
}

// MakeWithCap is a convenience function that creates a new Set
// instance (by value) with its underlying map of the given capacity properly
// initialized.
func MakeWithCap[T comparable](capacity int) Set[T] {
	return Set[T]{elements: make(map[T]struct{}, capacity)}
}

// New creates and returns a pointer to a new Set instance
// with its underlying map properly initialized.
func New[T comparable]() *Set[T] {
	return &Set[T]{elements: make(map[T]struct{})}
}

// Add inserts an element into the set.
// If the element is already present, it won't be added again.
func (s *Set[T]) Add(element T) {
	if s.elements == nil {
		s.elements = make(map[T]struct{})
	}
	s.elements[element] = struct{}{}
}

// AddFromArray inserts all provided array elements into the set.
// If an element is already present, it won't be added again.
func (s *Set[T]) AddFromArray(elements []T) {
	for _, element := range elements {
		s.Add(element)
	}
}

// Remove deletes an element from the set.
// If the element is not present, the function has no effect.
func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

// Has checks if the given element exists in the set.
// It returns true if the element is found, otherwise false.
func (s *Set[T]) Has(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Iter returns an iterator function (iter.Seq[T]) that yields each element in
// the set. The iterator follows the iter.Seq[T] interface, calling the provided
// yield function on each element until the yield function returns false or all
// elements have been processed.
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		// Iterate over all keys in the underlying map.
		for k := range s.elements {
			// Call the yield function with the current element.
			// If yield returns false, stop the iteration early.
			if !yield(k) {
				return
			}
		}
	}
}

// Array returns a slice containing all elements present in the set. This is a
// convenient way to retrieve all elements if needed for further operations.
func (s *Set[T]) Array() []T {
	keys := make([]T, 0, len(s.elements))
	for k := range s.elements {
		keys = append(keys, k)
	}
	return keys
}

// SortedArray returns a sorted slice containing all elements present in the set.
// This is a convenient way to retrieve all elements if needed for further operations.
func SortedArray[T cmp.Ordered](s Set[T]) []T {
	keys := make([]T, 0, len(s.elements))
	for k := range s.elements {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// MapToArray returns a slice containing all elements present in the set.
// The elements are transformed with the mapper.
func MapToArray[T comparable, R any](s Set[T], mapper func(element T) R) []R {
	result := make([]R, 0, len(s.elements))
	for k := range s.elements {
		result = append(result, mapper(k))
	}
	return result
}

// Size returns the size of the set.
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Marshal encodes the set to JSON as an array.
// The order of the array elements is not deterministic.
func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Array())
}

// Unmarshal decodes the set from JSON stored as an array.
// The order of the array elements does not matter,
// as well as the duplicate elements.
func (s *Set[T]) UnmarshalJSON(in []byte) error {
	var arr []T
	if err := json.Unmarshal(in, &arr); err != nil {
		return err
	}
	*s = FromArray(arr)
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler for Set[T].
func (s Set[T]) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(s.Array()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for Set[T].
func (s *Set[T]) UnmarshalBinary(data []byte) error {
	var arr []T
	buf := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buf).Decode(&arr); err != nil {
		return err
	}
	*s = FromArray(arr)
	return nil
}

// FromArray makes a new set with all the array elements in it.
// The order of the array elements does not matter,
// as well as the duplicate elements.
func FromArray[T comparable](elements []T) Set[T] {
	set := Set[T]{
		elements: make(map[T]struct{}, len(elements)),
	}
	for _, element := range elements {
		set.elements[element] = struct{}{}
	}
	return set
}
