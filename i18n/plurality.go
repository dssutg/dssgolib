package i18n

func PluralIdxByLang(lang Lang, count int64) int {
	return pluralIdxByLangMap[lang].pluralIdxByLang(count)
}
