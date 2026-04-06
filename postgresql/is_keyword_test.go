package postgresql_test

import (
	"strings"
	"testing"

	"github.com/dssutg/dssgolib/postgresql"
)

func TestIsKeywordAllKeywords(t *testing.T) {
	t.Parallel()

	for _, word := range postgresql.Keywords {
		if !postgresql.IsKeyword(word) {
			t.Errorf("IsKeyword(%q) = false; want true", word)
		}
		if !postgresql.IsKeyword(strings.ToUpper(word)) {
			t.Errorf("IsKeyword(%q) = false; want true", strings.ToUpper(word))
		}
	}
}

func BenchmarkIsKeyword(b *testing.B) {
	for b.Loop() {
		for _, word := range postgresql.Keywords {
			ok := postgresql.IsKeyword(word)
			_ = ok
		}
	}
}

// Implementation with map for speed comparison.
// It does not even consider case-insensitivity
// unlike the optimized function. That should make
// map even faster but it is still slower!
func BenchmarkIsKeywordWithMap(b *testing.B) {
	for b.Loop() {
		for _, word := range postgresql.Keywords {
			_, ok := postgresql.KeywordSet[word]
			_ = ok
		}
	}
}

func TestIsKeywordAllNonKeywords(t *testing.T) {
	t.Parallel()

	nonKeywords := []string{
		"aaaaaaaa",
		"ifx",
		"package",
		"arrays",
		"arra",
	}

	for _, word := range nonKeywords {
		if postgresql.IsKeyword(word) {
			t.Errorf("IsKeyword(%q) = true; want false", word)
		}
		if postgresql.IsKeyword(strings.ToUpper(word)) {
			t.Errorf("IsKeyword(%q) = true; want false", strings.ToUpper(word))
		}
	}
}
