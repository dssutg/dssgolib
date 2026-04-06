// Package semver provides routines for manipulating semantic versions.
package semver

import (
	"strconv"
	"strings"
)

// Compare returns -1, if version a < version b (older), 1, if a > b (newer), 0, if same.
func Compare(a, b []int) int {
	la, lb := len(a), len(b)
	n := max(lb, la)

	for i := range n {
		va, vb := 0, 0

		if i < la {
			va = a[i]
		}
		if i < lb {
			vb = b[i]
		}

		if va != vb {
			if va < vb {
				return -1
			}
			return 1
		}
	}

	return 0
}

// ParseStringForSorting attempts to parse a version string like "14", "14.3", "12beta1", etc.
// It returns a slice of ints for numeric components (non-numeric suffixes are ignored).
// Non-parsable components are treated as -1 so they sort lower.
func ParseStringForSorting(v string) []int {
	// Keep only prefix until first non-digit/dot/underscore char
	trimmed := v
	for i, r := range v {
		if r != '.' && r != '_' && (r < '0' || r > '9') {
			trimmed = v[:i]
			break
		}
	}
	if trimmed == "" {
		return []int{0}
	}
	parts := strings.FieldsFunc(trimmed, func(r rune) bool { return r == '.' || r == '_' })
	out := make([]int, len(parts))
	for i, p := range parts {
		if n, err := strconv.Atoi(p); err == nil {
			out[i] = n
		} else {
			out[i] = -1
		}
	}
	return out
}
