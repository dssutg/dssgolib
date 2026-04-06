package utils

import (
	"fmt"
	"strconv"
)

// BoolToString converts boolean to its lowercase string representation.
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// StringToBool converts common string representations to bool.
// Returns true for: "1", "t", "true", "y", "yes", "on" (case-insensitive).
// Any other value returns false.
func StringToBool(s string) bool {
	switch len(s) {
	case 1: // "t", "y", or "1"
		switch s[0] {
		case 'T', 't', 'Y', 'y', '1':
			return true
		default:
			return false
		}
	case 2: // "on"
		return (s[0] == 'O' || s[0] == 'o') && (s[1] == 'N' || s[1] == 'n')
	case 3: // "yes"
		return (s[0] == 'Y' || s[0] == 'y') &&
			(s[1] == 'E' || s[1] == 'e') &&
			(s[2] == 'S' || s[2] == 's')
	case 4: // "true"
		return (s[0] == 'T' || s[0] == 't') && (s[1] == 'R' || s[1] == 'r') &&
			(s[2] == 'U' || s[2] == 'u') && (s[3] == 'E' || s[3] == 'e')
	default:
		return false
	}
}

// Itoa converts any integer type to its base-10 string representation.
func Itoa[T Integer](v T) string {
	switch v := any(v).(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uintptr:
		return strconv.FormatUint(uint64(v), 10)
	default:
		panic("unreachable")
	}
}

// ParseDecimalUint8 is a type-safe wrapper around strconv.ParseUint(s, 10, 8).
func ParseDecimalUint8(s string) (uint8, error) {
	x, err := strconv.ParseUint(s, 10, 8)
	return uint8(x), err
}

// ParseDecimalInt8 is a type-safe wrapper around strconv.ParseInt(s, 10, 8).
func ParseDecimalInt8(s string) (int8, error) {
	x, err := strconv.ParseInt(s, 10, 8)
	return int8(x), err
}

// ParseDecimalUint16 is a type-safe wrapper around strconv.ParseUint(s, 10, 16).
func ParseDecimalUint16(s string) (uint16, error) {
	x, err := strconv.ParseUint(s, 10, 16)
	return uint16(x), err
}

// ParseDecimalInt16 is a type-safe wrapper around strconv.ParseInt(s, 10, 16).
func ParseDecimalInt16(s string) (int16, error) {
	x, err := strconv.ParseInt(s, 10, 16)
	return int16(x), err
}

// ParseDecimalUint32 is a type-safe wrapper around strconv.ParseUint(s, 10, 32).
func ParseDecimalUint32(s string) (uint32, error) {
	x, err := strconv.ParseUint(s, 10, 32)
	return uint32(x), err
}

// ParseDecimalInt32 is a type-safe wrapper around strconv.ParseInt(s, 10, 32).
func ParseDecimalInt32(s string) (int32, error) {
	x, err := strconv.ParseInt(s, 10, 32)
	return int32(x), err
}

// ParseDecimalUint64 is a type-safe wrapper around strconv.ParseUint(s, 10, 64).
func ParseDecimalUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// ParseDecimalInt64 is a type-safe wrapper around strconv.ParseInt(s, 10, 64).
func ParseDecimalInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseFloat32 is a type-safe wrapper around strconv.ParseFloat(s, 32).
func ParseFloat32(s string) (float32, error) {
	x, err := strconv.ParseFloat(s, 32)
	return float32(x), err
}

// ParseFloat64 is a type-safe wrapper around strconv.ParseFloat(s, 64).
func ParseFloat64(s string) (float64, error) {
	x, err := strconv.ParseFloat(s, 64)
	return float64(x), err
}

// FormatFloat32 is a type-safe wrapper around
// strconv.FormatFloat(x, 'f', -1, 32).
func FormatFloat32(x float32) string {
	return strconv.FormatFloat(float64(x), 'f', -1, 32)
}

// FormatFloat64 is a type-safe wrapper around
// strconv.FormatFloat(x, 'f', -1, 64).
func FormatFloat64(x float64) string {
	return strconv.FormatFloat(x, 'f', -1, 64)
}

// ByteToBool returns the boolean representation of the provided byte value.
func ByteToBool(b byte) bool {
	return b != 0
}

// BoolToByte returns the byte representation of the provided boolean value.
func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// FormatHMSTime returns the textual HH:MM:SS representation of the
// provided time in seconds.
func FormatHMSTime(seconds uint64) string {
	// Buffer enough to hold the worst case when
	// seconds is 2^64-1, which is "5124095576030431:00:15".
	var buf [22]byte

	// Convert total seconds to hours (hh), minutes (mm), and seconds (ss).
	// Both mm and ss are always within 0..59 range, but hh maybe be larger.
	ss := seconds
	mm := ss / 60
	ss %= 60
	hh := mm / 60
	mm %= 60

	// Fast path for the two-digit hour.
	if hh <= 99 {
		buf[0] = byte(hh/10) + '0'
		buf[1] = byte(hh%10) + '0'
		buf[2] = ':'
		buf[3] = byte(mm/10) + '0'
		buf[4] = byte(mm%10) + '0'
		buf[5] = ':'
		buf[6] = byte(ss/10) + '0'
		buf[7] = byte(ss%10) + '0'
		return string(buf[:8])
	}

	// Slower path for larger hours.
	// Build the string in reverse.
	i := len(buf)

	// Format MM:SS:.
	buf[i-1] = byte(ss%10) + '0'
	buf[i-2] = byte(ss/10) + '0'
	buf[i-3] = ':'
	buf[i-4] = byte(mm%10) + '0'
	buf[i-5] = byte(mm/10) + '0'
	buf[i-6] = ':'

	// Format HH.
	i -= 7
	for hh != 0 {
		buf[i] = byte(hh%10) + '0'
		hh /= 10
		i--
	}

	// i is behind the first hour digit in the string, so add 1 to it.
	return string(buf[i+1:])
}

// FormatHMTime returns the textual HH:MM representation of the
// provided time in seconds.
func FormatHMTime(seconds uint64) string {
	// Buffer enough to hold the worst case when
	// seconds is 2^64-1, which is "5124095576030431:00".
	var buf [19]byte

	// Convert total seconds to hours (hh), minutes (mm), and seconds (ss).
	// Both mm and ss are always within 0..59 range, but hh maybe be larger.
	mm := seconds / 60
	hh := mm / 60
	mm %= 60

	// Fast path for the two-digit hour.
	if hh <= 99 {
		buf[0] = byte(hh/10) + '0'
		buf[1] = byte(hh%10) + '0'
		buf[2] = ':'
		buf[3] = byte(mm/10) + '0'
		buf[4] = byte(mm%10) + '0'
		return string(buf[:5])
	}

	// Slower path for larger hours.
	// Build the string in reverse.
	i := len(buf)

	// Format MM:.
	buf[i-1] = byte(mm%10) + '0'
	buf[i-2] = byte(mm/10) + '0'
	buf[i-3] = ':'

	// Format HH.
	i -= 4
	for hh != 0 {
		buf[i] = byte(hh%10) + '0'
		hh /= 10
		i--
	}

	// i is behind the first hour digit in the string, so add 1 to it.
	return string(buf[i+1:])
}

// FormatIPv4Octets returns the string representation of an IP address
// provided as 4 octets. The format is a.b.c.d (e.g., 127.0.0.1).
func FormatIPv4Octets(a, b, c, d byte) string {
	var buf [15]byte // enough to hold xxx.xxx.xxx.xxx

	s := buf[:0]

	s = strconv.AppendUint(s, uint64(a), 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, uint64(b), 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, uint64(c), 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, uint64(d), 10)

	return string(s)
}

// FormatIPv4OctetSlice returns the string representation of an IP address
// provided as octets in the byte slice. The format is a.b.c.d (e.g., 127.0.0.1).
// The minimum slice length is assumed to be 4, or this function will panic.
func FormatIPv4OctetSlice(octets []byte) string {
	var buf [15]byte // enough to hold xxx.xxx.xxx.xxx

	s := buf[:0]

	s = strconv.AppendUint(s, uint64(octets[0]), 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, uint64(octets[1]), 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, uint64(octets[2]), 10)
	s = append(s, '.')
	s = strconv.AppendUint(s, uint64(octets[3]), 10)

	return string(s)
}

// ToString converts v to string as defined by [fmt.Sprint].
func ToString[T any](v T) string {
	return fmt.Sprint(v)
}
