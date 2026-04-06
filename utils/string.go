// String manipulation routines.
package utils

import (
	"bytes"
	"strings"
	"unicode"
	"unicode/utf8"
)

type RuneOrByte interface {
	~rune | ~byte
}

// CapitalizeString returns the string with the first character in uppercase.
// If the string is empty, it returns the empty string.
func CapitalizeString(s string) string {
	if s == "" {
		return s
	}

	// Convert string to a slice of runes to properly handle Unicode characters.
	runes := []rune(s)

	// Convert the first rune to uppercase.
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

// UncapitalizeString returns the string with the first character in lowercase.
// If the string is empty, it returns the empty string.
func UncapitalizeString(s string) string {
	if s == "" {
		return s
	}

	// Convert string to a slice of runes to properly handle Unicode characters.
	runes := []rune(s)

	// Convert the first rune to uppercase.
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

// StringHasLeadingWhitespace reports whether s begins with a Unicode whitespace rune.
func StringHasLeadingWhitespace(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsSpace(r)
}

// StringHasTrailingWhitespace reports whether s ends with a Unicode whitespace rune.
func StringHasTrailingWhitespace(s string) bool {
	if s == "" {
		return false
	}
	r, _ := utf8.DecodeLastRuneInString(s)
	return unicode.IsSpace(r)
}

// StringHasWhitespace reports whether s has a Unicode whitespace rune.
func StringHasWhitespace(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// StringContainsFold reports whether substr is within s ignoring
// case-sensitivity. Both s and substr are interpreted as UTF-8 strings,
// Characters are equal under simple Unicode case-folding, which is a
// more general form of case-insensitivity.
func StringContainsFold(s, substr string) bool {
	if substr == "" {
		return true
	}
	n, m := len(s), len(substr)
	if n < m {
		return false
	}
	for i := 0; i <= n-m; i++ {
		if strings.EqualFold(s[i:i+m], substr) {
			return true
		}
	}
	return false
}

var unicodeTo1251Partial = [0x2122 + 1]byte{
	0x00A0: 160, 0x00A4: 164, 0x00A6: 166, 0x00A7: 167, 0x00A9: 169, 0x00AB: 171,
	0x00AC: 172, 0x00AD: 173, 0x00AE: 174, 0x00B0: 176, 0x00B1: 177, 0x00B5: 181,
	0x00B6: 182, 0x00B7: 183, 0x00BB: 187, 0x0401: 168, 0x0402: 128, 0x0403: 129,
	0x0404: 170, 0x0405: 189, 0x0406: 178, 0x0407: 175, 0x0408: 163, 0x0409: 138,
	0x040A: 140, 0x040B: 142, 0x040C: 141, 0x040E: 161, 0x040F: 143, 0x0410: 192,
	0x0411: 193, 0x0412: 194, 0x0413: 195, 0x0414: 196, 0x0415: 197, 0x0416: 198,
	0x0417: 199, 0x0418: 200, 0x0419: 201, 0x041A: 202, 0x041B: 203, 0x041C: 204,
	0x041D: 205, 0x041E: 206, 0x041F: 207, 0x0420: 208, 0x0421: 209, 0x0422: 210,
	0x0423: 211, 0x0424: 212, 0x0425: 213, 0x0426: 214, 0x0427: 215, 0x0428: 216,
	0x0429: 217, 0x042A: 218, 0x042B: 219, 0x042C: 220, 0x042D: 221, 0x042E: 222,
	0x042F: 223, 0x0430: 224, 0x0431: 225, 0x0432: 226, 0x0433: 227, 0x0434: 228,
	0x0435: 229, 0x0436: 230, 0x0437: 231, 0x0438: 232, 0x0439: 233, 0x043A: 234,
	0x043B: 235, 0x043C: 236, 0x043D: 237, 0x043E: 238, 0x043F: 239, 0x0440: 240,
	0x0441: 241, 0x0442: 242, 0x0443: 243, 0x0444: 244, 0x0445: 245, 0x0446: 246,
	0x0447: 247, 0x0448: 248, 0x0449: 249, 0x044A: 250, 0x044B: 251, 0x044C: 252,
	0x044D: 253, 0x044E: 254, 0x044F: 255, 0x0451: 184, 0x0452: 144, 0x0453: 131,
	0x0454: 186, 0x0455: 190, 0x0456: 179, 0x0457: 191, 0x0458: 188, 0x0459: 154,
	0x045A: 156, 0x045B: 158, 0x045C: 157, 0x045E: 162, 0x045F: 159, 0x0490: 165,
	0x0491: 180, 0x2013: 150, 0x2014: 151, 0x2018: 145, 0x2019: 146, 0x201A: 130,
	0x201C: 147, 0x201D: 148, 0x201E: 132, 0x2020: 134, 0x2021: 135, 0x2022: 149,
	0x2026: 133, 0x2030: 137, 0x2039: 139, 0x203A: 155, 0x20AC: 136, 0x2116: 185,
	0x2122: 153,
}

var win1251ToUnicode = [128]uint16{
	0x0402, 0x0403, 0x201A, 0x0453, 0x201E, 0x2026, 0x2020, 0x2021, 0x20AC,
	0x2030, 0x0409, 0x2039, 0x040A, 0x040C, 0x040B, 0x040F, 0x0452, 0x2018,
	0x2019, 0x201C, 0x201D, 0x2022, 0x2013, 0x2014, 0x0000, 0x2122, 0x0459,
	0x203A, 0x045A, 0x045C, 0x045B, 0x045F, 0x00A0, 0x040E, 0x045E, 0x0408,
	0x00A4, 0x0490, 0x00A6, 0x00A7, 0x0401, 0x00A9, 0x0404, 0x00AB, 0x00AC,
	0x00AD, 0x00AE, 0x0407, 0x00B0, 0x00B1, 0x0406, 0x0456, 0x0491, 0x00B5,
	0x00B6, 0x00B7, 0x0451, 0x2116, 0x0454, 0x00BB, 0x0458, 0x0405, 0x0455,
	0x0457, 0x0410, 0x0411, 0x0412, 0x0413, 0x0414, 0x0415, 0x0416, 0x0417,
	0x0418, 0x0419, 0x041A, 0x041B, 0x041C, 0x041D, 0x041E, 0x041F, 0x0420,
	0x0421, 0x0422, 0x0423, 0x0424, 0x0425, 0x0426, 0x0427, 0x0428, 0x0429,
	0x042A, 0x042B, 0x042C, 0x042D, 0x042E, 0x042F, 0x0430, 0x0431, 0x0432,
	0x0433, 0x0434, 0x0435, 0x0436, 0x0437, 0x0438, 0x0439, 0x043A, 0x043B,
	0x043C, 0x043D, 0x043E, 0x043F, 0x0440, 0x0441, 0x0442, 0x0443, 0x0444,
	0x0445, 0x0446, 0x0447, 0x0448, 0x0449, 0x044A, 0x044B, 0x044C, 0x044D,
	0x044E, 0x044F,
}

// UTF8ToWin1251Soft returns the slice of bytes in Windows-1251 encoding
// converted from UTF-8 string. Invalid characters are ignored and not
// added to the result bytes.
func UTF8ToWin1251Soft(s string) []byte {
	var b bytes.Buffer

	// Windows-1251 always takes less bytes than UTF-8 so
	// we can allocate the buffer just once.
	b.Grow(len(s))

	for _, c := range s {
		// Ignore invalid runes.
		if c < 0 {
			continue
		}

		// Fast path for ASCII.
		if c <= 127 {
			b.WriteByte(byte(c))
			continue
		}

		// If greater than the maximum character in the table.
		if c >= rune(len(unicodeTo1251Partial)) {
			continue
		}

		// Look up the value in the table.
		if c1251 := unicodeTo1251Partial[c]; c1251 != 0 {
			b.WriteByte(c1251)
		}
	}

	return b.Bytes()
}

// Win1251ToStringSoft returns the UTF-8 encoded string converted from
// NUL-terminated slice of bytes in Windows-1251 encoding. Invalid
// characters are ignored and not added to the result bytes.
func Win1251ToStringSoft(b []byte) string {
	var sb strings.Builder

	// Preallocate for UCS-2 (BMP characters).
	sb.Grow(len(b) * 2)

	for _, c := range b {
		if c == 0 { // NUL terminator
			break
		}
		switch {
		case c < 128:
			sb.WriteByte(c) // ASCII
		case c == 152:
			// Ignore non-existent codepoint in Windows-1251
		default:
			sb.WriteRune(rune(win1251ToUnicode[c-128])) // Look up in table
		}
	}

	return sb.String()
}

// TrimStringToRuneCount returns s limited to at most maxRunes.
// If the maxRunes is less than or equal to zero, an empty string
// is returned. String does not have to be valid UTF-8.
// This function does not allocate.
func TrimStringToRuneCount(s string, maxRunes int) string {
	// Fast path: zero or negative rune count.
	if maxRunes <= 0 {
		return ""
	}

	// Fast path: if length in runes is <= maxRunes,
	// return original string.
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}

	// Find byte index of the max-th rune.
	i := 0
	for r := 0; r < maxRunes && i < len(s); r++ {
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
	}

	// Slice to the max-th rune.
	return s[:i]
}

// StringCollapseSpaces removes both leading and trailing whitespace from s,
// and replaces any consecutive whitespace with a single ASCII space (0x20). It
// preserves the original byte sequence (including invalid UTF-8). If the
// string already contains only non-consecutive ASCII spaces, the original
// string without leading and trailing whitespace is returned.
// A new allocation is done only if cleaning is needed and cannot be done
// by slicing the string.
func StringCollapseSpaces(s string) string {
	// Remove leading and trailing whitespace.
	// No allocation, only slicing.
	s = strings.TrimSpace(s)

	// Check if this string has already the proper format.
	// If so, no point to build it - return immediately.
	if StringHasOnlyNonConsecutiveASCIISpaces(s) {
		return s
	}

	var b strings.Builder

	b.Grow(len(s)) // preallocate to the max possible length

	first := true
	for part := range strings.FieldsSeq(s) {
		if !first {
			b.WriteByte(' ')
		}
		first = false
		b.WriteString(part)
	}

	return b.String()
}

// StringHasOnlyNonConsecutiveASCIISpaces reports whether s contains only
// ASCII non-consecutive spaces (0x20). Returns false if any other whitespace
// (ASCII or non-ASCII) or consecutive spaces are present.
// Returns true if s does not have whitespace.
func StringHasOnlyNonConsecutiveASCIISpaces(s string) bool {
	n := len(s)
	i := 0
	for i < n { // iterate by bytes to avoid UTF-8 decoding
		c := s[i]
		switch c {
		case ' ':
			if i != n-1 && s[i+1] == ' ' {
				return false // consecutive ASCII space
			}
			i++
		case '\t', '\n', '\v', '\f', '\r':
			return false // other ASCII whitespace
		default:
			if c < utf8.RuneSelf { // ASCII non-space byte
				i++
				continue
			}
			r, size := utf8.DecodeRuneInString(s[i:])
			if unicode.IsSpace(r) {
				return false // non-ASCII whitespace rune
			}
			i += size // skip rune
		}
	}
	return true
}

// TruncateStringBytes returns the string of at most maxBytes bytes long.
// If the rune at the end of the string becomes invalid after slicing
// the string bytes, the trailing invalid bytes are droped.
func TruncateStringBytes(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return "" // limit allows for only empty string
	}

	if len(s) <= maxBytes {
		return s // already within the limit
	}

	// Take first maxBytes bytes, then back up to valid UTF-8 boundary.
	s = s[:maxBytes]
	for len(s) > 0 {
		if r, n := utf8.DecodeLastRuneInString(s); r != utf8.RuneError || n != 1 {
			break // rune is valid, quit
		}
		s = s[:len(s)-1] // drop invalid byte, until can decode the rune
	}

	return s
}

// IsASCIIControl reports whether the provided character
// is ASCII control character.
func IsASCIIControl[T RuneOrByte](c T) bool {
	return c <= 0x1F || c == 0x7F
}

// SplitToSlice splits s to substrings separated by sep to parts slice.
// This function does not allocate memory because the result is put to the
// provided parts slice. The parts that cannot fit in the slice are
// discarded. The total number of split parts including discarded ones
// is returned.
func SplitToSlice(parts []string, s, sep string) int {
	n := 0
	for part := range strings.SplitSeq(s, sep) {
		if n < len(parts) {
			parts[n] = part
		}
		n++
	}
	return n
}

// GetSplitPart splits s to substrings separated by sep.
// It returns the first part whose index is the partIdx and
// whether the part is found. If part not found, empty string is returned.
// This function does not allocate memory.
func GetSplitPart(s, sep string, partIdx int) (string, bool) {
	i := 0
	for part := range strings.SplitSeq(s, sep) {
		if partIdx == i {
			return part, true
		}
		i++
	}
	return "", false
}

// StringHasPrefixFold reports whether the string s begins with prefix.
// Strings are considered equal under simple Unicode case-folding, which is a
// more general form of case-insensitivity.
func StringHasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

// StringHasSuffixFold reports whether the string s ends with suffix.
// Strings are considered equal under simple Unicode case-folding, which is a
// more general form of case-insensitivity.
func StringHasSuffixFold(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}

// StringTrimPrefixFold returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
// Strings are considered equal under simple Unicode case-folding, which is a
// more general form of case-insensitivity.
func StringTrimPrefixFold(s, prefix string) string {
	if StringHasPrefixFold(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// StringTrimSuffixFold returns s without the provided trailing suffix string
// If s doesn't end with suffix, s is returned unchanged.
// Strings are considered equal under simple Unicode case-folding, which is a
// more general form of case-insensitivity.
func StringTrimSuffixFold(s, suffix string) string {
	if StringHasSuffixFold(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}
