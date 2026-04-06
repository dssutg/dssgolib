package utils

import (
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Strftime is a convenient wrapper around [AppendStrftime].
// It returns the formatted date as string.
func Strftime(format string, t time.Time) string {
	b := make([]byte, 0, len(format)*2)
	b = AppendStrftime(b, format, t)
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// AppendStrftime is the famous strftime() function that provides a very
// powerful, concise and intuitive way of formatting dates.
// It appends the formatted date t to the buffer b according to the
// format string.
//
// Most format specifiers of GNU date are implemented. It aims to be greatly
// compatible with the GNU strftime, although not completely and not all
// features are supported, because it's not a one-to-one port.
//
// See https://man7.org/linux/man-pages/man1/date.1.html for more information.
func AppendStrftime(b []byte, format string, t time.Time) []byte {
	// Our primary goal is to make this function very fast,
	// so we sacrifice readability a lot.

	for len(format) != 0 {
		// Jump to the first possible format specifier
		// ignoring all the rest bytes for analysis,
		// as well as appending them all at once.
		// IndexByte may utilize SIMD instructions
		// so it's more efficient.
		i := strings.IndexByte(format, '%')
		if i < 0 || i == len(format)-1 {
			// No format specifier found or the percent is at the end of the string.
			// Write the rest of the format string and return.
			b = append(b, format...)
			break
		}

		// Write the bytes before the format specifier.
		b = append(b, format[:i]...)

		// Skip everything up to percent (inclusive).
		format = format[i+1:]

		// Dispatch the padding specifier.
		readPadByte := format[0]
		padByte := readPadByte
		customPad := padByte == '0' || padByte == '-' || padByte == '_'
		if customPad {
			if len(format) == 1 { // if nothing after padding
				b = append(b, '%', readPadByte) // append as-is
				break
			}
			format = format[1:] // skip padding character
		} else {
			padByte = '0' // default padding
		}

		// Dispatch the format specifier.
		switch format[0] {
		case 'F': // full date; like %+4Y-%m-%d
			m := int(t.Month())
			d := t.Day()
			m1 := byte(m/10) + '0' //nolint:gosec
			m2 := byte(m%10) + '0' //nolint:gosec
			d1 := byte(d/10) + '0' //nolint:gosec
			d2 := byte(d%10) + '0' //nolint:gosec
			b = strftimeAppendPadInt4(b, t.Year(), '0')
			b = append(b, '-', m1, m2, '-', d1, d2)
		case 'D', 'x': // date; same as %m/%d/%y
			// %x is locale's date representation (e.g., 12/31/99)
			// but we make it the same as %D.
			y := t.Year()
			m := int(t.Month())
			d := t.Day()
			y1 := byte(y/10%10) + '0' //nolint:gosec
			y2 := byte(y%10) + '0'    //nolint:gosec
			m1 := byte(m/10) + '0'    //nolint:gosec
			m2 := byte(m%10) + '0'    //nolint:gosec
			d1 := byte(d/10) + '0'    //nolint:gosec
			d2 := byte(d%10) + '0'    //nolint:gosec
			b = append(b, m1, m2, '/', d1, d2, '/', y1, y2)
		case 'T', 'X': // time; same as %H:%M:%S
			// %X is locale's time representation (e.g., 23:13:48)
			// but we make it the same as %T.
			h := t.Hour()
			m := t.Minute()
			s := t.Second()
			h1 := byte(h/10) + '0' //nolint:gosec
			h2 := byte(h%10) + '0' //nolint:gosec
			m1 := byte(m/10) + '0' //nolint:gosec
			m2 := byte(m%10) + '0' //nolint:gosec
			s1 := byte(s/10) + '0' //nolint:gosec
			s2 := byte(s%10) + '0' //nolint:gosec
			b = append(b, h1, h2, ':', m1, m2, ':', s1, s2)
		case 'R': // 24-hour hour and minute; same as %H:%M
			h := t.Hour()
			m := t.Minute()
			h1 := byte(h/10) + '0' //nolint:gosec
			h2 := byte(h%10) + '0' //nolint:gosec
			m1 := byte(m/10) + '0' //nolint:gosec
			m2 := byte(m%10) + '0' //nolint:gosec
			b = append(b, h1, h2, ':', m1, m2)
		case 'r': // locale's 12-hour clock time (e.g., 11:11:04 PM)
			h := t.Hour()
			m := t.Minute()
			s := t.Second()
			p := byte('a')
			if h >= 12 {
				p = 'p'
				h -= 12
			}
			if h == 0 {
				h = 12
			}
			h1 := byte(h/10) + '0'
			h2 := byte(h%10) + '0'
			m1 := byte(m/10) + '0' //nolint:gosec
			m2 := byte(m%10) + '0' //nolint:gosec
			s1 := byte(s/10) + '0' //nolint:gosec
			s2 := byte(s%10) + '0' //nolint:gosec
			b = append(b, h1, h2, ':', m1, m2, ':', s1, s2, ' ', p, 'm')
		case 'c': // locale's date and time (e.g., Thu Mar  3 23:05:25 2005)
			d := t.Day()
			h := t.Hour()
			m := t.Minute()
			s := t.Second()
			b = append(b, t.Weekday().String()[:3]...)
			b = append(b, ' ')
			b = append(b, t.Month().String()[:3]...)
			b = append(b, ' ')
			d2 := byte(d%10) + '0' //nolint:gosec
			if d >= 10 {
				d1 := byte(d/10) + '0' //nolint:gosec
				b = append(b, d1, d2, ' ')
			} else {
				b = append(b, ' ', d2, ' ')
			}
			h1 := byte(h/10) + '0' //nolint:gosec
			h2 := byte(h%10) + '0' //nolint:gosec
			m1 := byte(m/10) + '0' //nolint:gosec
			m2 := byte(m%10) + '0' //nolint:gosec
			s1 := byte(s/10) + '0' //nolint:gosec
			s2 := byte(s%10) + '0' //nolint:gosec
			b = append(b, h1, h2, ':', m1, m2, ':', s1, s2, ' ')
			b = strftimeAppendPadInt4(b, t.Year(), '0')
		case 'k': // hour, space padded (0..23); same as %_H
			h := t.Hour()
			h2 := byte(h%10) + '0' //nolint:gosec
			if h >= 10 {
				h1 := byte(h/10) + '0' //nolint:gosec
				b = append(b, h1, h2)
			} else {
				b = append(b, ' ', h2)
			}
		case 'l': // hour, space padded (1..12); same as %_I
			h := t.Hour()
			if h >= 12 {
				h -= 12
			}
			if h == 0 {
				h = 12
			}
			h2 := byte(h%10) + '0'
			if h >= 10 {
				h1 := byte(h/10) + '0'
				b = append(b, h1, h2)
			} else {
				b = append(b, ' ', h2)
			}
		case 'e': // day of month, space padded; same as %_d
			d := t.Day()
			d2 := byte(d%10) + '0' //nolint:gosec
			if d >= 10 {
				d1 := byte(d/10) + '0' //nolint:gosec
				b = append(b, d1, d2)
			} else {
				b = append(b, ' ', d2)
			}
		case 'H': // hour (00..23)
			b = strftmeAppendPadInt2(b, t.Hour(), padByte)
		case 'M': // minute (00..59)
			b = strftmeAppendPadInt2(b, t.Minute(), padByte)
		case 'S': // second (00..60)
			b = strftmeAppendPadInt2(b, t.Second(), padByte)
		case 'Y': // year
			b = strftimeAppendPadInt4(b, t.Year(), padByte)
		case 'G': // year of ISO week number (see %V); normally useful only with %V
			b = strftimeAppendPadInt4(b, t.Year(), padByte)
		case 'y': // last two digits of year (00..99)
			b = strftmeAppendPadInt2(b, t.Year()%100, padByte)
		case 'g': // last two digits of year of ISO week number (see %G)
			b = strftmeAppendPadInt2(b, (t.Year()%100+100)%100, padByte)
		case 'C': // century; like %Y, except omit last two digits (e.g., 20)
			b = strftmeAppendPadInt2(b, t.Year()/100, padByte)
		case 'q': // quarter of year (1..4)
			b = append(b, byte((int(t.Month())-1)/3+1)+'0') //nolint:gosec
		case 'd': // day of month (e.g., 01)
			b = strftmeAppendPadInt2(b, t.Day(), padByte)
		case 'j': // day of year (001..366)
			b = strftimeAppendPadInt3(b, t.YearDay(), padByte)
		case 'm': // month (01..12)
			b = strftmeAppendPadInt2(b, int(t.Month()), padByte)
		case 'w': // day of week (0..6); 0 is Sunday
			b = append(b, byte(t.Weekday())+'0') //nolint:gosec
		case 'V': // ISO week number, with Monday as first day of week (01..53)
			_, week := t.ISOWeek()
			b = strftmeAppendPadInt2(b, week, padByte)
		case 'U': // week number of year, with Sunday as first day of week (00..53)
			jan1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()).Weekday()
			firstSunOff := (7 - int(jan1)) % 7
			week := (t.YearDay() - 1 - firstSunOff + 7) / 7
			b = strftmeAppendPadInt2(b, week, padByte)
		case 'W': // week number of year, with Monday as first day of week (00..53)
			jan1 := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()).Weekday()
			firstMonOff := (7 - ((int(jan1) + 6) % 7)) % 7
			week := (t.YearDay() - 1 - firstMonOff + 7) / 7
			b = strftmeAppendPadInt2(b, week, padByte)
		case 'I': // hour (01..12)
			h := t.Hour()
			if h >= 12 {
				h -= 12
			}
			if h == 0 {
				h = 12
			}
			b = strftmeAppendPadInt2(b, h, padByte)
		case 'P': // like %p, but lower case
			if t.Hour() >= 12 {
				b = append(b, 'p', 'm')
			} else {
				b = append(b, 'a', 'm')
			}
		case 'p': // locale's equivalent of either AM or PM; blank if not known
			if t.Hour() >= 12 {
				b = append(b, 'P', 'M')
			} else {
				b = append(b, 'A', 'M')
			}
		case 's': // seconds since the Epoch (1970-01-01 00:00 UTC)
			b = strconv.AppendInt(b, t.Unix(), 10)
		case 'N': // nanoseconds (000000000..999999999)
			b = strconv.AppendInt(b, t.UnixNano(), 10)
		case 'u': // day of week (1..7); 1 is Monday
			wday := byte(t.Weekday()) //nolint:gosec
			if wday == 0 {
				wday = 7
			}
			b = append(b, wday+'0')
		case 'A': // locale's full weekday name (e.g., Sunday)
			b = append(b, t.Weekday().String()...)
		case 'a': // locale's abbreviated weekday name (e.g., Sun)
			b = append(b, t.Weekday().String()[:3]...)
		case 'B': // locale's full month name (e.g., January)
			b = append(b, t.Month().String()...)
		case 'b', 'h': // locale's abbreviated month name (e.g., Jan)
			b = append(b, t.Month().String()[:3]...)
		case 't': // a tab
			b = append(b, '\t')
		case 'n': // a newline
			b = append(b, '\n')
		case 'Z': // alphabetic time zone abbreviation (e.g., EDT)
			name, _ := t.Zone() // offset in seconds east of UTC
			b = append(b, name...)
		case 'z': // +hhmm numeric time zone (e.g., -0400)
			_, offset := t.Zone() // offset in seconds east of UTC
			sign := byte('+')
			if offset < 0 {
				sign = '-'
				offset = -offset
			}
			h := offset / 3600
			m := (offset % 3600) / 60
			h1 := byte(h/10) + '0' //nolint:gosec
			h2 := byte(h%10) + '0' //nolint:gosec
			m1 := byte(m/10) + '0' //nolint:gosec
			m2 := byte(m%10) + '0' //nolint:gosec
			b = append(b, sign, h1, h2, m1, m2)
		case ':':
			n := len(format)
			switch {
			case n >= 4:
				if format[1] == ':' && format[2] == ':' && format[3] == 'z' {
					// %:::z - numeric time zone with : to necessary precision (e.g., -04, +05:30)
					_, offset := t.Zone() // offset in seconds east of UTC
					sign := byte('+')
					if offset < 0 {
						sign = '-'
						offset = -offset
					}
					s := offset
					m := s / 60
					h := m / 60
					s %= 60
					m %= 60
					h1 := byte(h/10) + '0' //nolint:gosec
					h2 := byte(h%10) + '0' //nolint:gosec
					switch {
					case s != 0: // if non-zero seconds, minutes surely appended as well
						m1 := byte(m/10) + '0' //nolint:gosec
						m2 := byte(m%10) + '0' //nolint:gosec
						s1 := byte(s/10) + '0' //nolint:gosec
						s2 := byte(s%10) + '0' //nolint:gosec
						b = append(b, sign, h1, h2, ':', m1, m2, ':', s1, s2)
					case m != 0: // if zero seconds, then only append minutes if non-zero
						m1 := byte(m/10) + '0' //nolint:gosec
						m2 := byte(m%10) + '0' //nolint:gosec
						b = append(b, sign, h1, h2, ':', m1, m2)
					default:
						b = append(b, sign, h1, h2)
					}
					format = format[4:] // consume :::z
					continue            // skip consuming any byte because we've already done it
				}
			case n >= 3:
				if format[1] == ':' && format[2] == 'z' {
					// %::z - +hh:mm:ss numeric time zone (e.g., -04:00:00)
					_, offset := t.Zone() // offset in seconds east of UTC
					sign := byte('+')
					if offset < 0 {
						sign = '-'
						offset = -offset
					}
					s := offset
					m := s / 60
					h := m / 60
					s %= 60
					m %= 60
					h1 := byte(h/10) + '0' //nolint:gosec
					h2 := byte(h%10) + '0' //nolint:gosec
					m1 := byte(m/10) + '0' //nolint:gosec
					m2 := byte(m%10) + '0' //nolint:gosec
					s1 := byte(s/10) + '0' //nolint:gosec
					s2 := byte(s%10) + '0' //nolint:gosec
					b = append(b, sign, h1, h2, ':', m1, m2, ':', s1, s2)
					format = format[3:] // consume ::z
					continue            // skip consuming any byte because we've already done it
				}
			case n >= 2:
				if format[1] == 'z' {
					// %:z - +hh:mm numeric time zone (e.g., -04:00)
					_, offset := t.Zone() // offset in seconds east of UTC
					sign := byte('+')
					if offset < 0 {
						sign = '-'
						offset = -offset
					}
					h := offset / 3600
					m := (offset % 3600) / 60
					h1 := byte(h/10) + '0' //nolint:gosec
					h2 := byte(h%10) + '0' //nolint:gosec
					m1 := byte(m/10) + '0' //nolint:gosec
					m2 := byte(m%10) + '0' //nolint:gosec
					b = append(b, sign, h1, h2, ':', m1, m2)
					format = format[2:] // consume :z
					continue            // skip consuming any byte because we've already done it
				}
			}
			b = append(b, '%', ':')
		case '%': // a literal %
			b = append(b, '%')
		default: // append as-is any other byte
			if customPad {
				b = append(b, '%', readPadByte, format[0])
			} else {
				b = append(b, '%', format[0])
			}
		}

		format = format[1:] // skip the format specifier
	}

	return b
}

func strftmeAppendPadInt2(b []byte, x int, padByte byte) []byte {
	x2 := byte(x%10) + '0' //nolint:gosec
	switch padByte {
	case '0':
		x1 := byte(x/10) + '0' //nolint:gosec
		b = append(b, x1, x2)
	case '-':
		if x >= 10 {
			x1 := byte(x/10) + '0' //nolint:gosec
			b = append(b, x1, x2)
		} else {
			b = append(b, x2)
		}
	case '_':
		if x >= 10 {
			x1 := byte(x/10) + '0' //nolint:gosec
			b = append(b, x1, x2)
		} else {
			b = append(b, ' ', x2)
		}
	}
	return b
}

func strftimeAppendPadInt3(b []byte, x int, padByte byte) []byte {
	x3 := byte(x%10) + '0' //nolint:gosec
	switch padByte {
	case '0':
		x1 := byte(x/100%10) + '0' //nolint:gosec
		x2 := byte(x/10%10) + '0'  //nolint:gosec
		b = append(b, x1, x2, x3)
	case '-':
		switch {
		case x >= 100:
			x1 := byte(x/100%10) + '0'
			x2 := byte(x/10%10) + '0'
			b = append(b, x1, x2, x3)
		case x >= 10:
			x2 := byte(x/10%10) + '0'
			b = append(b, x2, x3)
		default:
			b = append(b, x3)
		}
	case '_':
		switch {
		case x >= 100:
			x1 := byte(x/100%10) + '0'
			x2 := byte(x/10%10) + '0'
			b = append(b, x1, x2, x3)
		case x >= 10:
			x2 := byte(x/10%10) + '0'
			b = append(b, ' ', x2, x3)
		default:
			b = append(b, ' ', ' ', x3)
		}
	}
	return b
}

func strftimeAppendPadInt4(b []byte, x int, padByte byte) []byte {
	if x >= 10000 {
		b = strconv.AppendInt(b, int64(x), 10)
		return b
	}
	x4 := byte(x%10) + '0' //nolint:gosec
	switch padByte {
	case '0':
		x1 := byte(x/1000) + '0'   //nolint:gosec
		x2 := byte(x/100%10) + '0' //nolint:gosec
		x3 := byte(x/10%10) + '0'  //nolint:gosec
		b = append(b, x1, x2, x3, x4)
	case '-':
		switch {
		case x >= 1000:
			x1 := byte(x/1000) + '0'
			x2 := byte(x/100%10) + '0'
			x3 := byte(x/10%10) + '0'
			b = append(b, x1, x2, x3, x4)
		case x >= 100:
			x2 := byte(x/100%10) + '0'
			x3 := byte(x/10%10) + '0'
			b = append(b, x2, x3, x4)
		case x >= 10:
			x3 := byte(x/10%10) + '0'
			b = append(b, x3, x4)
		default:
			b = append(b, x4)
		}
	case '_':
		switch {
		case x >= 1000:
			x1 := byte(x/1000) + '0'
			x2 := byte(x/100%10) + '0'
			x3 := byte(x/10%10) + '0'
			b = append(b, x1, x2, x3, x4)
		case x >= 100:
			x2 := byte(x/100%10) + '0'
			x3 := byte(x/10%10) + '0'
			b = append(b, ' ', x2, x3, x4)
		case x >= 10:
			x3 := byte(x/10%10) + '0'
			b = append(b, ' ', ' ', x3, x4)
		default:
			b = append(b, ' ', ' ', ' ', x4)
		}
	}
	return b
}
