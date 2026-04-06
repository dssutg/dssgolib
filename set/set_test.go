package set

import (
	"testing"
)

func TestMake(t *testing.T) {
	t.Parallel()

	s := Make[int]()
	if len(s.Array()) != 0 {
		t.Errorf("len(s.Array()) = %d, want 0", len(s.Array()))
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	s := New[int]()
	if len(s.Array()) != 0 {
		t.Errorf("len(s.Array()) = %d, want 0", len(s.Array()))
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	s := New[int]()
	s.Add(1)
	if !s.Has(1) {
		t.Error("s.Has(1) = false, want true")
	}
}

func TestRemove(t *testing.T) {
	t.Parallel()

	s := New[int]()
	s.Add(1)
	s.Remove(1)
	if s.Has(1) {
		t.Error("s.Has(1) = true, want false")
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	s := New[int]()
	s.Add(1)
	if !s.Has(1) {
		t.Error("s.Has(1) = false, want true")
	}
	if s.Has(2) {
		t.Error("s.Has(2) = true, want false")
	}
}

func TestIter(t *testing.T) {
	t.Parallel()

	s := New[int]()
	s.Add(1)
	s.Add(2)

	count := 0
	s.Iter()(func(int) bool {
		count++

		return true
	})

	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

func TestArray(t *testing.T) {
	t.Parallel()

	s := New[int]()
	s.Add(1)
	s.Add(2)
	arr := s.Array()

	if len(arr) != 2 {
		t.Errorf("len(arr) = %d, want 2", len(arr))
	}

	// Check if the elements are correct
	want := map[int]struct{}{1: {}, 2: {}}
	for _, elem := range arr {
		if _, ok := want[elem]; !ok {
			t.Errorf("unwanted number %d in array", elem)
		}
		delete(want, elem)
	}

	if len(want) != 0 {
		t.Errorf("len(want) = %d, want 0", len(want))
	}
}
