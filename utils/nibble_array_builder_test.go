package utils

import "testing"

func TestWriteNibbleBE(t *testing.T) {
	var buf [4]byte

	b := MakeNibbleArrayBuilder(buf[:])
	b.WriteNibbleBE(0xAA)
	b.WriteNibbleBE(0xBB)
	b.WriteNibbleBE(0xCC)
	b.WriteNibbleBE(0xDD)

	want := [4]byte{0xab, 0xcd}

	if buf != want {
		t.Errorf("buf: %#v; want: %#v", buf, want)
	}
}

func TestWriteNibbleLE(t *testing.T) {
	var buf [4]byte

	b := MakeNibbleArrayBuilder(buf[:])
	b.WriteNibbleLE(0xAA)
	b.WriteNibbleLE(0xBB)
	b.WriteNibbleLE(0xCC)
	b.WriteNibbleLE(0xDD)

	want := [4]byte{0xba, 0xdc}

	if buf != want {
		t.Errorf("buf: %#v; want: %#v", buf, want)
	}
}
