package i18n

import (
	"time"
)

type RelDateLocale interface { //nolint:interfacebloat
	DaysAgo(days int64) string
	Yesterday() string
	HoursAgo(hours int64) string
	OneHourAgo() string
	MinutesAgo(minutes int64) string
	OneMinuteAgo() string
	SecondsAgo(seconds int64) string
	JustThen() string
	InSeconds(seconds int64) string
	InOneMinute() string
	InMinutes(minutes int64) string
	InOneHour() string
	InHours(hours int64) string
	Tomorrow() string
	InDays(days int64) string
}

const (
	minuteSecs = 60
	hourSecs   = minuteSecs * 60
	daySecs    = hourSecs * 24
)

// FormatRelDate formats relative date. Relative dates are the ones that are
// relative to today. For example, if today is May 10, then May 9 is formatted
// as simply Yesterday.
func FormatRelDate(someTime, now time.Time, dateLocale RelDateLocale) string {
	nowUnix := now.Unix()
	someTimeUnix := someTime.Unix()

	delta := nowUnix - someTimeUnix

	if delta < 0 {
		var absDelta int64
		if nowUnix > someTimeUnix {
			absDelta = nowUnix - someTimeUnix
		} else {
			absDelta = someTimeUnix - nowUnix
		}

		switch {
		case absDelta < 30:
			return dateLocale.JustThen()
		case absDelta < minuteSecs:
			return dateLocale.InSeconds(absDelta)
		case absDelta < 2*minuteSecs:
			return dateLocale.InOneMinute()
		case absDelta < hourSecs:
			return dateLocale.InMinutes(absDelta / minuteSecs)
		case absDelta/hourSecs == 1:
			return dateLocale.InOneHour()
		case absDelta < daySecs:
			return dateLocale.InHours(absDelta / hourSecs)
		case absDelta < daySecs*2:
			return dateLocale.Tomorrow()
		default:
			return dateLocale.InDays(absDelta / daySecs)
		}
	}

	switch {
	case delta < 30:
		return dateLocale.JustThen()
	case delta < minuteSecs:
		return dateLocale.SecondsAgo(delta)
	case delta < 2*minuteSecs:
		return dateLocale.OneMinuteAgo()
	case delta < hourSecs:
		return dateLocale.MinutesAgo(delta / minuteSecs)
	case delta/hourSecs == 1:
		return dateLocale.OneHourAgo()
	case delta < daySecs:
		return dateLocale.HoursAgo(delta / hourSecs)
	case delta < daySecs*2:
		return dateLocale.Yesterday()
	default:
		return dateLocale.DaysAgo(delta / daySecs)
	}
}
