package utils

import (
	"cmp"
	"slices"
)

// NumRange represents a numeric interval [Start, End].
type NumRange[T cmp.Ordered] struct {
	Start T // always 1st field
	End   T // always 2nd field
}

// MergeNumRanges merges the provided ranges in-place (no allocations)
// and returns the merged, sorted slice with the same backing array.
func MergeNumRanges[T cmp.Ordered](ranges []NumRange[T]) []NumRange[T] {
	n := len(ranges)
	if n == 0 {
		return ranges[:0]
	}

	// Sort by Start.
	slices.SortFunc(ranges, func(a, b NumRange[T]) int {
		return cmp.Compare(a.Start, b.Start)
	})

	write := 0 // write index for merged ranges (the Squeeze algorithm)
	for _, cur := range ranges {
		if write == 0 {
			ranges[write] = cur // 1st element is always left as-is
			write++
			continue
		}
		last := &ranges[write-1]
		switch {
		case last.End < cur.Start: // non-overlapping: move cur into next write pos
			ranges[write] = cur
			write++
		case cur.End > last.End: // overlapping/adjacent: extend the end if needed
			last.End = cur.End
		}
	}

	return ranges[:write] // merged prefix
}

// MergeNumsToRanges merges the array of numbers to ranges.
func MergeNumsToRanges[T int](v []T) []NumRange[T] {
	n := len(v)
	if n == 0 {
		return nil
	}

	// Determine how many merges there are going to be
	// to preallocate the result array of ranges.
	mergedCap := 0
	for i := 1; i < n; i++ {
		if v[i] != v[i-1]+1 {
			mergedCap++
		}
	}
	mergedCap++ // account for the final append

	merged := make([]NumRange[T], 0, mergedCap)

	// Fill the merged array.
	start := v[0]
	end := start
	for _, val := range v[1:] {
		if val != end+1 { // if differs from previous value by more than one
			merged = append(merged, NumRange[T]{Start: start, End: end})
			start = val
		}
		end = val
	}
	merged = append(merged, NumRange[T]{Start: start, End: end})

	return merged
}
