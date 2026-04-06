package i18n

import "fmt"

var RussianDateLocale russianDateLocale

type russianDateLocale struct{}

func (l russianDateLocale) DaysAgo(days int64) string {
	switch PluralIdxByLang(LangRussian, days) {
	case 0:
		return fmt.Sprintf("%d день назад", days)
	case 1:
		return fmt.Sprintf("%d дня назад", days)
	default:
		return fmt.Sprintf("%d дней назад", days)
	}
}

func (l russianDateLocale) Yesterday() string {
	return "вчера"
}

func (l russianDateLocale) HoursAgo(hours int64) string {
	switch PluralIdxByLang(LangRussian, hours) {
	case 0:
		return fmt.Sprintf("%d час назад", hours)
	case 1:
		return fmt.Sprintf("%d часа назад", hours)
	default:
		return fmt.Sprintf("%d часов назад", hours)
	}
}

func (l russianDateLocale) OneHourAgo() string {
	return "час назад"
}

func (l russianDateLocale) MinutesAgo(minutes int64) string {
	switch PluralIdxByLang(LangRussian, minutes) {
	case 0:
		return fmt.Sprintf("%d минута назад", minutes)
	case 1:
		return fmt.Sprintf("%d минуты назад", minutes)
	default:
		return fmt.Sprintf("%d минут назад", minutes)
	}
}

func (l russianDateLocale) OneMinuteAgo() string {
	return "одну минуту назад"
}

func (l russianDateLocale) SecondsAgo(seconds int64) string {
	switch PluralIdxByLang(LangRussian, seconds) {
	case 0:
		return fmt.Sprintf("%d секунда назад", seconds)
	case 1:
		return fmt.Sprintf("%d секунды назад", seconds)
	default:
		return fmt.Sprintf("%d секунд назад", seconds)
	}
}

func (l russianDateLocale) JustThen() string {
	return "только что"
}

func (l russianDateLocale) InSeconds(seconds int64) string {
	switch PluralIdxByLang(LangRussian, seconds) {
	case 0:
		return fmt.Sprintf("через %d секунду", seconds)
	case 1:
		return fmt.Sprintf("через %d секунды", seconds)
	default:
		return fmt.Sprintf("через %d секунд", seconds)
	}
}

func (l russianDateLocale) InOneMinute() string {
	return "через минуту"
}

func (l russianDateLocale) InMinutes(minutes int64) string {
	switch PluralIdxByLang(LangRussian, minutes) {
	case 0:
		return fmt.Sprintf("через %d минуту", minutes)
	case 1:
		return fmt.Sprintf("через %d минуты", minutes)
	default:
		return fmt.Sprintf("через %d минут", minutes)
	}
}

func (l russianDateLocale) InOneHour() string {
	return "через час"
}

func (l russianDateLocale) InHours(hours int64) string {
	switch PluralIdxByLang(LangRussian, hours) {
	case 0:
		return fmt.Sprintf("через %d час", hours)
	case 1:
		return fmt.Sprintf("через %d часа", hours)
	default:
		return fmt.Sprintf("через %d часов", hours)
	}
}

func (l russianDateLocale) Tomorrow() string {
	return "завтра"
}

func (l russianDateLocale) InDays(days int64) string {
	switch PluralIdxByLang(LangRussian, days) {
	case 0:
		return fmt.Sprintf("через %d день", days)
	case 1:
		return fmt.Sprintf("через %d дня", days)
	default:
		return fmt.Sprintf("через %d дней", days)
	}
}

func pluralIdxRussian(count int64) int {
	// 1 apple
	if count%10 == 1 && count%100 != 11 {
		return 0
	}
	// 2 apples
	if count%10 >= 2 && count%10 <= 4 && (count%100 < 10 || count%100 >= 20) {
		return 1
	}
	// 0 apples
	return 2
}
