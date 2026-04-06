package utils

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	// DayDuration is the duration of one day.
	DayDuration = 24 * time.Hour

	DaysPerWeek = 7

	SecondsPerMinute = 60
	SecondsPerHour   = 60 * SecondsPerMinute
	SecondsPerDay    = 24 * SecondsPerHour
	SecondsPerWeek   = DaysPerWeek * SecondsPerDay

	DaysPer400Years = 365*400 + 97
)

// Hour12 describes 12-hour format.
type Hour12 struct {
	Hour int  // 12-hour format (1..12)
	IsPM bool // Is this hour PM - true, AM - false
}

// String implements the [fmt.Stringer] interface.
func (hour12 Hour12) String() string {
	if hour12.IsPM {
		return fmt.Sprintf("%d PM", hour12.Hour)
	}
	return fmt.Sprintf("%d AM", hour12.Hour)
}

// To24Format returns the 24-hour format of the 12-format time.
// The 12-hour format is assumed to be valid.
func (hour12 Hour12) To24Format() int {
	// Handle special case.
	if hour12.Hour == 12 {
		if hour12.IsPM {
			return 12
		}
		return 0
	}
	// General conversion.
	if hour12.IsPM {
		return hour12.Hour + 12
	}
	return hour12.Hour
}

// MakeHour12From24 returns the 12-hour format of the 24-format time.
// If the provided hour is out of range, it is clamped to be within
// 0...23 hours.
func MakeHour12From24(hour24 int) Hour12 {
	// Clamp the hour to be valid.
	hour24 = min(max(0, hour24), 23)

	// General conversion.
	hour12 := Hour12{IsPM: hour24/12 != 0}
	hour24 %= 12

	// Handle special case.
	if hour24 == 0 {
		hour12.Hour = 12
	} else {
		hour12.Hour = hour24
	}

	return hour12
}

// IsWeekend reports whether the week day of the provided time
// is either Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

// GetTimeMidnight returns the same calendar date as the provided
// time but at the midnight time, i.e., the very start of the day
// described by the time point.
func GetTimeMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetDayDuration returns the duration of the day described by the provided time.
func GetDayDuration(t time.Time) time.Duration {
	return t.Sub(GetTimeMidnight(t))
}

// DurationBetweenDayTimes wraps subtraction of two day-based durations and
// returns the positive duration from start to end within a 24-hour day
// (wraps around if negative).
func DurationBetweenDayTimes(start, end time.Duration) time.Duration {
	diff := end - start

	if diff < 0 {
		return diff + DayDuration
	}

	return diff
}

// HMSToDuration returns the duration of the provided 24-hour
// format HH:MM:SS time.
func HMSToDuration(hour, minute, second int) time.Duration {
	return time.Duration(hour*60*60+minute*60+second) * time.Second
}

// SleepUntilHMS sleeps until the specified day time.
// If the current time is already later the provided one,
// then sleep until the same time on the next day.
func SleepUntilHMS(hour, minute, second int) {
	start := GetDayDuration(time.Now())
	end := HMSToDuration(hour, minute, second)

	duration := DurationBetweenDayTimes(start, end)

	time.Sleep(duration)
}

// IsLeapYear reports whether the provided year is leap.
//
// NOTE: Extracted from the standard time package.
// For some stupid reason, such a useful function is private there.
func IsLeapYear(year int) bool {
	// year%4 == 0 && (year%100 != 0 || year%400 == 0)
	// Bottom 2 bits must be clear.
	// For multiples of 25, bottom 4 bits must be clear.
	// Thanks to Cassio Neri for this trick.
	mask := 0xf
	if year%25 != 0 {
		mask = 3
	}

	return year&mask == 0
}

// DaysInMonth returns the number of days in the provided month.
// Leap years are considered.
//
// NOTE: Extracted from the standard time package.
// For some stupid reason, such a useful function is private there.
func DaysInMonth(month time.Month, year int) int {
	if month == time.February {
		if IsLeapYear(year) {
			return 29
		}
		return 28
	}

	// With the special case of February eliminated, the pattern is
	//	31 30 31 30 31 30 31 31 30 31 30 31
	// Adding m&1 produces the basic alternation;
	// adding (m>>3)&1 inverts the alternation starting in August.
	return 30 + int((month+month>>3)&1)
}

// The running sum of the number of days before each month.
// This is only for non-leap years.
var DaysBeforeMonthNonLeap = [...]int{
	0,                                    // before January
	0 + 31,                               // before February
	0 + 31 + 28,                          // before March
	0 + 31 + 28 + 31,                     // before April
	0 + 31 + 28 + 31 + 30,                // before May
	0 + 31 + 28 + 31 + 30 + 31,           // before June
	0 + 31 + 28 + 31 + 30 + 31 + 30,      // before July
	0 + 31 + 28 + 31 + 30 + 31 + 30 + 31, // before August
	0 + 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,                // before September
	0 + 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,           // before October
	0 + 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,      // before November
	0 + 31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30, // before December
}

// The running sum of the number of days before each month.
// This is only for leap years.
var DaysBeforeMonthLeap = [...]int{
	0,                                    // before January
	0 + 31,                               // before February
	0 + 31 + 29,                          // before March
	0 + 31 + 29 + 31,                     // before April
	0 + 31 + 29 + 31 + 30,                // before May
	0 + 31 + 29 + 31 + 30 + 31,           // before June
	0 + 31 + 29 + 31 + 30 + 31 + 30,      // before July
	0 + 31 + 29 + 31 + 30 + 31 + 30 + 31, // before August
	0 + 31 + 29 + 31 + 30 + 31 + 30 + 31 + 31,                // before September
	0 + 31 + 29 + 31 + 30 + 31 + 30 + 31 + 31 + 30,           // before October
	0 + 31 + 29 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,      // before November
	0 + 31 + 29 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30, // before December
}

// DaysBeforeMonth returns the number of days in the provided month.
// Leap years are considered.
func DaysBeforeMonth(month time.Month, year int) int {
	if IsLeapYear(year) {
		return DaysBeforeMonthLeap[month]
	}
	return DaysBeforeMonthNonLeap[month]
}

// MonthStartWeekday returns the weekday of the first day in this month.
// This function along with [DaysInMonth] is a must-have to build
// a month calendar.
func MonthStartWeekday(month time.Month, year int, loc *time.Location) time.Weekday {
	return time.Date(year, month, 1, 0, 0, 0, 0, loc).Weekday()
}

// PrintStandardMonthCalendar prints the standard calendar
// starting with Sunday.
func PrintStandardMonthCalendar(w io.Writer, month time.Month, year int, loc *time.Location) {
	header1Len := len(month.String()) + 1 + DigitCountUint64(AbsToUint64(int64(year)))
	if year < 0 {
		header1Len++ // for minus sign
	}

	header2 := "Su Mo Tu We Th Fr Sa"

	// Print space to center header 1.
	fmt.Fprint(w, strings.Repeat(" ", (len(header2)-header1Len)/2))
	// Print header 1.
	fmt.Fprintf(w, "%s %d\n", month, year)

	// Print header 2.
	fmt.Fprintln(w, header2)

	// Compute calendar properties.
	startDay := int(MonthStartWeekday(month, year, loc))
	days := DaysInMonth(month, year)

	// Print space to skip previous month's weekdays.
	fmt.Fprint(w, strings.Repeat(" ", startDay*3))

	// Print the calendar.
	col := startDay
	for day := 1; day <= days; day++ {
		fmt.Fprintf(w, "%2d ", day)
		col++
		if col%7 == 0 {
			fmt.Fprintln(w)
		}
	}
	if col%7 != 0 {
		fmt.Fprintln(w)
	}
}
