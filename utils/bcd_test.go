package utils

import (
	"bytes"
	"testing"
)

func TestBCDRoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		end  byte
		want string
	}{
		{"", 0xF, ""}, // empty input -> terminated immediately
		{"0", 0xF, "0"},
		{"1", 0xF, "1"},
		{"12", 0xF, "12"},
		{"123", 0xF, "123"},
		{"0123456789", 0xF, "0123456789"},
		{"1a2b3", 0xF, "123"},
	}

	for _, tc := range cases {
		{
			buf := make([]byte, (len(tc.in)+1)/2+1)
			StringToBCDSliceBE(buf, tc.in, tc.end)
			got := BCDToStringBE(buf)
			if got != tc.want {
				t.Fatalf("BE round-trip: input=%q end=0x%X => got %q want %q", tc.in, tc.end, got, tc.want)
			}
		}

		{
			buf := make([]byte, (len(tc.in)+1)/2+1)
			StringToBCDSliceLE(buf, tc.in, tc.end)
			got := BCDToStringLE(buf)
			if got != tc.want {
				t.Fatalf("LE round-trip: input=%q end=0x%X => got %q want %q", tc.in, tc.end, got, tc.want)
			}
		}
	}
}

func TestBCDBEvsLEEncodingDifferences(t *testing.T) {
	t.Parallel()

	in := "1234"
	be := make([]byte, (len(in)+1)/2+1)
	le := make([]byte, (len(in)+1)/2+1)
	StringToBCDSliceBE(be, in, 0xF)
	StringToBCDSliceLE(le, in, 0xF)

	// BCDToString* should return same textual digits for both
	if BCDToStringBE(be) != BCDToStringLE(le) {
		t.Fatalf("BE and LE round-trip mismatch: BE=%q LE=%q", BCDToStringBE(be), BCDToStringLE(le))
	}

	// but raw bytes should differ (nibble ordering)
	if bytes.Equal(be, le) {
		t.Fatalf("got equal bytes for BE and LE; want different")
	}
}

func TestBCDToStringStopsOnNibbleGreaterThan9(t *testing.T) {
	t.Parallel()

	// craft a byte slice containing a digit and then a termination nibble > 9
	// For BE: high nibble valid, low nibble > 9 should stop after high digit.
	be := []byte{0x1A} // high=1, low=10 -> should produce "1"
	if BCDToStringBE(be) != "1" {
		t.Fatalf("BCDToStringBE did not stop at invalid low nibble: got %q want %q", BCDToStringBE(be), "1")
	}

	// For LE: low nibble is stored in low position, but BCDToStringLE reads low/high reversed in bit ops.
	// Using 0xA1 (low=0x1, high=0xA) should yield "1" then stop.
	le := []byte{0xA1}
	if BCDToStringLE(le) != "1" {
		t.Fatalf("BCDToStringLE did not stop at invalid nibble: got %q want %q", BCDToStringLE(le), "1")
	}
}

func BenchmarkStringToBCDSliceBE(b *testing.B) {
	in := "0123456789012345678901234567890123456789"
	buf := make([]byte, (len(in)+1)/2+1)

	for b.Loop() {
		StringToBCDSliceBE(buf, in, 0xF)
	}
}

func BenchmarkStringToBCDSliceLE(b *testing.B) {
	in := "0123456789012345678901234567890123456789"
	buf := make([]byte, (len(in)+1)/2+1)

	for b.Loop() {
		StringToBCDSliceLE(buf, in, 0xF)
	}
}

func BenchmarkBCDToStringBE(b *testing.B) {
	in := make([]byte, 20)
	// prepare a valid BCD byte slice for BE: 0..19 -> bytes 0x01 0x23 ...
	for i := range 10 {
		in[i] = byte(((2*i)&0xF)<<4 | ((2*i + 1) & 0xF))
	}

	for b.Loop() {
		_ = BCDToStringBE(in)
	}
}

func BenchmarkBCDToStringLE(b *testing.B) {
	in := make([]byte, 20)
	for i := range 10 {
		// LE expects low nibble first in storage; create equivalent pattern
		in[i] = byte(((2*i+1)&0xF)<<4 | ((2 * i) & 0xF))
	}

	for b.Loop() {
		_ = BCDToStringLE(in)
	}
}
