package utils

import (
	"math"
	"testing"
)

func TestDigitCountUint64(t *testing.T) {
	t.Parallel()

	// Slow but surely correct implementation of the same function.
	digitCountSlow := func(x uint64) int {
		if x == 0 {
			return 1
		}
		n := 0
		for x > 0 {
			x /= 10
			n++
		}
		return n
	}

	assertCount := func(x uint64) {
		got := DigitCountUint64(x)
		want := digitCountSlow(x)
		if got != want {
			t.Fatalf("DigitCountUint64(%d) = %d; want %d", x, got, want)
		}
	}

	for i := range math.MaxUint16 + 1 {
		assertCount(uint64(i))
	}

	assertCount(1e4)
	assertCount(1e5)
	assertCount(1e6)
	assertCount(1e7)
	assertCount(1e8)
	assertCount(1e9)
	assertCount(1e10)
	assertCount(1e11)
	assertCount(1e12)
	assertCount(1e13)
	assertCount(1e14)
	assertCount(1e15)
	assertCount(1e16)
	assertCount(1e17)
	assertCount(1e18)

	assertCount(1<<64 - 1)
}

func TestNearlyEqual(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		a, b, eps float64
		want      bool
	}{
		{a: 0.3, b: 0.30000000000000004, eps: 0.0000001, want: true},
		{a: 0.30000000000000004, b: 0.3, eps: 0.0000001, want: true},
		{a: 4, b: 4, eps: 0, want: true},
		{a: 5, b: 2, eps: 0, want: false},
		{a: 2, b: 5, eps: 0, want: false},
		{a: 2, b: 5, eps: 100, want: true},
	}

	for _, tc := range testCases {
		got := NearlyEqual(tc.a, tc.b, tc.eps)
		if got != tc.want {
			t.Errorf("NearlyEqual(%v, %v, %v) = %v; want = %v", tc.a, tc.b, tc.eps, got, tc.want)
		}
	}
}

func BenchmarkIntSin(b *testing.B) {
	InitIntSinTable()

	for b.Loop() {
		for i := -512; i <= 512; i++ {
			s := IntSin(i * IntSinAngles / 360)
			_ = s
		}
	}
}

func BenchmarkSin(b *testing.B) {
	for b.Loop() {
		for i := -512.0; i <= 512; i++ {
			s := math.Sin(i * math.Pi / 180)
			_ = s
		}
	}
}
