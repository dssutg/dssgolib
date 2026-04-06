package utils

import (
	"slices"
	"testing"
)

func TestNaturalCmp(t *testing.T) {
	t.Parallel()

	want := [14]string{
		"a1", "a2", "a9", "a10", "a11", "a23", "a143",
		"b1", "b2", "b9", "b10", "b11", "b23", "b143",
	}

	v := want
	slices.SortFunc(v[:], NaturalCmp)
	got := v

	if v != want {
		t.Errorf("got %v; want %v", got, want)
	}
}
