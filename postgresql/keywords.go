package postgresql

import "github.com/dssutg/dssgolib/postgresql/internal"

// Keywords is a list of PostgreSQL keywords.
// This array is sorted lexicographically in ascending
// order, all keywords are ASCII lowercase.
//
// The optimized for speed [IsKeyword] function allows to
// case-insensitively determine whether a string is in this list.
var Keywords = internal.Keywords

// KeywordSet is a set of lowercase PostgreSQL keywords.
var KeywordSet = make(map[string]struct{}, len(Keywords))

func init() {
	for _, word := range Keywords {
		KeywordSet[word] = struct{}{}
	}
}
