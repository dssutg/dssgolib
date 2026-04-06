package mapx

import (
	"sync"
)

// Map is a type-safe wrapper around sync.Map.
type Map[K comparable, V any] struct {
	m sync.Map
}

// Store stores the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Load returns the value for a key and whether it was found.
// If not found, the zero value is returned.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	if m == nil {
		var zero V
		return zero, false
	}
	v, ok := m.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

// Has reports whether the key exists in the map.
func (m *Map[K, V]) Has(key K) bool {
	if m == nil {
		return false
	}
	_, has := m.m.Load(key)
	return has
}

// Get behaves the same as the Load method but drops the second (ok) argument.
// Use it when zero value is fine to use for non-existent entries.
func (m *Map[K, V]) Get(key K) V {
	value, _ := m.Load(key)
	return value
}

// LoadOrStore returns the existing value if present; otherwise stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.m.LoadOrStore(key, value)
	return a.(V), loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := m.m.LoadAndDelete(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// CompareAndSwap swaps the value for key from old to new if the current value equals old.
// Note: old must be comparable (runtime compare done by sync.Map).
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.m.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the entry for key if its value equals old.
func (m *Map[K, V]) CompareAndDelete(key K, old V) bool {
	return m.m.CompareAndDelete(key, old)
}

// Swap sets the value for key and returns the previous value if any.
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	p, loaded := m.m.Swap(key, value)
	if !loaded {
		var zero V
		return zero, false
	}
	return p.(V), true
}

// Range calls fn for each key and value present in the map.
// If fn returns false, range stops.
func (m *Map[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil {
		return
	}
	m.m.Range(func(k, v any) bool {
		return fn(k.(K), v.(V))
	})
}

// UnsortedKeys returns an unsorted slice of keys of the map.
func (m *Map[K, V]) UnsortedKeys() (keys []K) {
	m.Range(func(k K, _ V) bool {
		keys = append(keys, k)
		return true
	})
	return keys
}

// UnsortedValues returns an unsorted slice of values of the map.
func (m *Map[K, V]) UnsortedValues() (values []V) {
	m.Range(func(_ K, v V) bool {
		values = append(values, v)
		return true
	})
	return values
}

// Clear deletes all entries by ranging and deleting each key.
func (m *Map[K, V]) Clear() {
	if m == nil {
		return
	}
	m.m.Clear()
}

func (m *Map[K, V]) RegularMap() map[K]V {
	result := make(map[K]V)
	m.Range(func(key K, value V) bool {
		result[key] = value
		return true
	})
	return result
}
