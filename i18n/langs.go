// Package i18n contains routines for internationalization.
package i18n

type Lang string

const (
	LangEnglish Lang = "en"
	LangRussian Lang = "ru"
)

type langTableEntry struct {
	pluralIdxByLang func(count int64) int
	relDateLocale   RelDateLocale
}

var pluralIdxByLangMap = map[Lang]langTableEntry{
	LangEnglish: {pluralIdxByLang: pluralIdxEnglish, relDateLocale: EnglishDateLocale},
	LangRussian: {pluralIdxByLang: pluralIdxRussian, relDateLocale: RussianDateLocale},
}
