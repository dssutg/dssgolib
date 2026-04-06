package fuzzy

import (
	"cmp"
	"math"
	"slices"
	"strings"
)

// FuzzySearch performs a fuzzy search over array using provided helpers.
// - query: the search string.
// - array: slice of items to search.
// - cleanString: normalizes strings (e.g., lowercasing, removing diacritics).
// - getItemName: returns the searchable name for each item.
// Returns a slice of matching items ordered by descending match quality.
func FuzzySearch[T any](
	query string,
	array []T,
	cleanString func(string) string,
	getItemName func(T) string,
) []T {
	cleanQuery := cleanString(query)

	if cleanQuery == "" {
		return slices.Clone(array)
	}

	type scoredItem struct {
		item  T
		score float64
	}

	var results []scoredItem

	lowestScore := -1.0
	highestScore := math.Inf(1)

	for _, item := range array {
		itemName := cleanString(getItemName(item))

		var score float64

		if strings.Contains(itemName, cleanQuery) {
			// Give the highest score to the exact match.
			score = highestScore
		} else {
			// Fuzzy match.
			matchIndex := 0
			score = 0.0

			for _, r := range cleanQuery {
				// Look for rune in itemName starting at matchIndex.
				idx := strings.IndexRune(itemName[matchIndex:], r)
				if idx == -1 {
					// If any character is not found, set the lowest score.
					score = lowestScore
					break
				}
				// Actual index relative to full string.
				actualIdx := matchIndex + idx
				score += 1.0 - float64(actualIdx)/float64(len(itemName))
				matchIndex = actualIdx + 1
			}
		}

		if score != lowestScore {
			results = append(results, scoredItem{item: item, score: score})
		}
	}

	// Sort by descending score. Items with highest score remain at top.
	slices.SortStableFunc(results, func(a, b scoredItem) int {
		return -cmp.Compare(a.score, b.score)
	})

	out := make([]T, len(results))
	for i, r := range results {
		out[i] = r.item
	}

	return out
}
