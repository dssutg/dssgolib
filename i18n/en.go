package i18n

import "fmt"

var EnglishDateLocale englishDateLocale

type englishDateLocale struct{}

func (l englishDateLocale) DaysAgo(days int64) string {
	return fmt.Sprintf("%d days ago", days)
}

func (l englishDateLocale) Yesterday() string {
	return "yesterday"
}

func (l englishDateLocale) HoursAgo(hours int64) string {
	return fmt.Sprintf("%d hours ago", hours)
}

func (l englishDateLocale) OneHourAgo() string {
	return "1 hour ago"
}

func (l englishDateLocale) MinutesAgo(minutes int64) string {
	return fmt.Sprintf("%d minutes ago", minutes)
}

func (l englishDateLocale) OneMinuteAgo() string {
	return "a minute ago"
}

func (l englishDateLocale) SecondsAgo(seconds int64) string {
	return fmt.Sprintf("%d seconds ago", seconds)
}

func (l englishDateLocale) JustThen() string {
	return "just then"
}

func (l englishDateLocale) InSeconds(seconds int64) string {
	return fmt.Sprintf("%d seconds", seconds)
}

func (l englishDateLocale) InOneMinute() string {
	return "in a minute"
}

func (l englishDateLocale) InMinutes(minutes int64) string {
	return fmt.Sprintf("in %d minutes", minutes)
}

func (l englishDateLocale) InOneHour() string {
	return "in 1 hour"
}

func (l englishDateLocale) InHours(hours int64) string {
	return fmt.Sprintf("in %d hours", hours)
}

func (l englishDateLocale) Tomorrow() string {
	return "tomorrow"
}

func (l englishDateLocale) InDays(days int64) string {
	return fmt.Sprintf("in %d days", days)
}

func pluralIdxEnglish(count int64) int {
	if count != 1 {
		return 1
	}
	return 0
}
