package ring

import (
	"slices"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	t.Parallel()

	capacity := 5
	buf := New[int](capacity)

	if buf.capacity != capacity {
		t.Errorf("buf.capacity = %d, want %d", buf.capacity, capacity)
	}
	if len(buf.buffer) != capacity {
		t.Errorf("len(buf.buffer) %d, want %d", len(buf.buffer), capacity)
	}
}

func TestAddAndGetAll(t *testing.T) {
	t.Parallel()

	buf := New[int](3)

	for i := 1; i <= 3; i++ {
		buf.Add(i)
	}

	got := buf.GetAll()
	want := []int{1, 2, 3}

	if !slices.Equal(got, want) {
		t.Errorf("buf.GetAll() = %v, want %v", got, want)
	}
}

func TestAddBeyondCapacity(t *testing.T) {
	t.Parallel()

	buf := New[int](3)

	for i := 1; i <= 4; i++ {
		buf.Add(i)
	}

	got := buf.GetAll()
	want := []int{2, 3, 4}

	if !slices.Equal(got, want) {
		t.Errorf("buf.GetAll() = %v, want %v", got, want)
	}
}

func TestWrapAround(t *testing.T) {
	t.Parallel()

	buf := New[int](3)

	for i := 1; i <= 5; i++ {
		buf.Add(i)
	}

	got := buf.GetAll()
	want := []int{3, 4, 5}

	if !slices.Equal(got, want) {
		t.Errorf("buf.GetAll() = %v, want %v", got, want)
	}
}

func TestWrapAround2(t *testing.T) {
	t.Parallel()

	buf := New[int](3)

	for i := 1; i <= 12; i++ {
		buf.Add(i)
	}

	got := buf.GetAll()
	want := []int{10, 11, 12}

	if !slices.Equal(got, want) {
		t.Errorf("buf.GetAll() = %v, want %v", got, want)
	}
}

func TestGetAllWhenEmpty(t *testing.T) {
	t.Parallel()

	buf := New[int](3)

	got := buf.GetAll()
	want := []int{}

	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
