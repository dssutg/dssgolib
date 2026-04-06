// Map manipulation routines.
package utils

import (
	"cmp"
	"iter"
	"slices"
)

// Entry represents a key-value pair from a map.
type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// MapEntries represents a list of key-value pairs from a map.
type MapEntries[K comparable, V any] []MapEntry[K, V]

// Keys returns a slice containing all keys of the given map entries.
func (e MapEntries[K, V]) Keys() []K {
	entries := make([]K, 0, len(e))
	for _, v := range e {
		entries = append(entries, v.Key)
	}
	return entries
}

// Values returns a slice containing all values of the given map entries.
func (e MapEntries[K, V]) Values() []V {
	entries := make([]V, 0, len(e))
	for _, v := range e {
		entries = append(entries, v.Value)
	}
	return entries
}

// GetMapKeys returns a slice containing all keys of the given map. It accepts a
// map with keys of any comparable type K and values of any type V, then
// iterates over the map to collect all keys into a slice.
//
// Benchmarks show that it is approximately two times faster than
// slices.Collect(maps.Keys(m)). The standard version uses iterators,
// and does not preallocate the result slice. This is why it's slower.
// So this function is not only clearer but also faster than the
// standard notation.
func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetSortedMapKeys returns a sorted slice containing all keys of the given map.
// It accepts a map with keys of any comparable type K and values of any type V,
// then iterates over the map to collect all keys into a slice.
func GetSortedMapKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// GetMapValues returns a slice containing all values of the given map.
// It accepts a map with keys of any comparable type K and values of any type V,
// then iterates over the map to collect all values into a slice.
//
// Benchmarks show that it is approximately two times faster than
// slices.Collect(maps.Keys(m)). The standard version uses iterators,
// and does not preallocate the result slice. This is why it's slower.
// So this function is not only clearer but also faster than the
// standard notation.
func GetMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// GetMapEntries returns a slice containing all entries of the given map. It
// accepts a map with keys of any comparable type K and values of any type V,
// then iterates over the map to collect all key-value pairs into a slice of
// Entry structs.
func GetMapEntries[K comparable, V any](m map[K]V) MapEntries[K, V] {
	entries := make(MapEntries[K, V], 0, len(m))
	for k, v := range m {
		entries = append(entries, MapEntry[K, V]{Key: k, Value: v})
	}
	return entries
}

// GetSortedMapEntries returns a slice containing all entries of the given map,
// sorted by key in ascending order. The key type K must be ordered.
func GetSortedMapEntries[K cmp.Ordered, V any](m map[K]V) MapEntries[K, V] {
	entries := GetMapEntries(m)
	slices.SortStableFunc(entries, func(a, b MapEntry[K, V]) int {
		return cmp.Compare(a.Key, b.Key)
	})
	return entries
}

// MapKeyExists checks if a key exists in the provided map.
func MapKeyExists[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

// IterSeq2Keys returns an iterator over keys in seq.
func IterSeq2Keys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range seq {
			if !yield(k) {
				return
			}
		}
	}
}

// IterSeq2Values returns an iterator over values in seq.
func IterSeq2Values[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}

// CollectIterSeq2 returns map entries from iterator seq2.
func CollectIterSeq2[K comparable, V any](seq iter.Seq2[K, V]) MapEntries[K, V] {
	var e MapEntries[K, V]
	for k, v := range seq {
		e = append(e, MapEntry[K, V]{Key: k, Value: v})
	}
	return e
}

// ReversedMap returns the reverse of the map m that maps values
// to keys.
func ReversedMap[K comparable, V comparable](m map[K]V) map[V]K {
	rv := make(map[V]K, len(m))
	for k, v := range m {
		rv[v] = k
	}
	return rv
}

// InitMap ensures map is initialized. If the map has already
// been initialized, this function is no-op.
func InitMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return make(map[K]V)
	}
	return m
}
