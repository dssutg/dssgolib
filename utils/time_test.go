package utils

import "testing"

func TestHour12(t *testing.T) {
	t.Parallel()

	for hour24 := range 24 {
		if MakeHour12From24(hour24).To24Format() != hour24 {
			t.Errorf("bad conversion for hour %d", hour24)
		}
	}
}
