package kvcfg

import "slices"

// MarshalStringMap returns the kv-config bytes from the given string map.
// The map keys in the config are sorted lexicographically.
func MarshalStringMap(m map[string]string) []byte {
	// Collect keys and estimate the result slice capacity.
	est := 0
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
		est += len(key) + 1 + len(m[key]) + 1 // key + '=' + value + '\n'
	}
	slices.Sort(keys)

	buf := make([]byte, 0, est)

	// Marshal the map.
	for _, k := range keys {
		buf = append(buf, k...)
		buf = append(buf, '=')
		buf = append(buf, m[k]...)
		buf = append(buf, '\n')
	}

	return buf
}
