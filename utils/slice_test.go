package utils

import (
	"fmt"
	"slices"
	"testing"
)

func TestSqueezeSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input       []string
		wantLen     int
		wantBacking []string
	}{
		{
			input:       []string{"", "", ""},
			wantLen:     0,
			wantBacking: []string{"", "", ""},
		},
		{
			input:       []string{},
			wantLen:     0,
			wantBacking: []string{},
		},
		{
			input:       []string{"", "", "a"},
			wantLen:     1,
			wantBacking: []string{"a", "", ""},
		},
		{
			input:       []string{"", "a", "b"},
			wantLen:     2,
			wantBacking: []string{"a", "b", ""},
		},
		{
			input:       []string{"a", "b", "c"},
			wantLen:     3,
			wantBacking: []string{"a", "b", "c"},
		},
		{
			input:       []string{"a", "", "", "b", "c"},
			wantLen:     3,
			wantBacking: []string{"a", "b", "c", "", ""},
		},
		{
			input:       []string{"a", "", "", "b", "", "c"},
			wantLen:     3,
			wantBacking: []string{"a", "b", "c", "", "", ""},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%#v", tc.input), func(t *testing.T) {
			t.Parallel()
			orig := slices.Clone(tc.input)
			got := SqueezeSlice(orig)
			if len(got) != tc.wantLen {
				t.Errorf("len(got) = %d, want %d", len(got), tc.wantLen)
			}
			if !slices.Equal(orig, tc.wantBacking) {
				t.Errorf("orig = %v, want %v", orig, tc.wantBacking)
			}
		})
	}

	if SliceEvery([]int{11, 22, 33}, func(x int) bool {
		return x < 0
	}) {
		t.Errorf("SliceEvery must return false (< 0)")
	}

	if !SliceEvery([]int{11, 22, 33}, func(x int) bool {
		return x >= 0
	}) {
		t.Errorf("SliceEvery must return true (>= 0)")
	}

	if !SliceEvery([]int{}, func(x int) bool {
		return x >= 0
	}) {
		t.Errorf("SliceEvery must return true for empty arrays")
	}
}
