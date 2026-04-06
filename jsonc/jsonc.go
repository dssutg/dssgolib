// Package jsonc removes comments from JSON.
// This allows to parse JSON files with Go-like
// single-line and multi-line comments.
package jsonc

import (
	"bytes"
	"encoding/json"
	"errors"
	"unicode/utf8"
)

// JSONUnmarshal is the JSON parser used by the functions of
// this package. It defaults to the Go standard [json.Unmarshal]
// but you can overwrite it to use another parser.
var JSONUnmarshal = json.Unmarshal

// ErrInvalidUTF8 is returned by Sanitize if the data is not valid UTF-8.
var ErrInvalidUTF8 = errors.New("invalid UTF-8")

// Sanitize removes all comments from JSONC data but does not validate the data.
// It returns [ErrInvalidUTF8] if the data is not valid UTF-8.
func Sanitize(data []byte) ([]byte, error) {
	if !utf8.Valid(data) {
		return nil, ErrInvalidUTF8
	}
	return sanitize(data), nil
}

// State machine flags.
const (
	hasCommentRunes byte = 1 << iota
	isString
	isCommentLine
	isCommentBlock
	checkNext
)

func sanitize(data []byte) []byte {
	var state byte
	return bytes.Map(func(r rune) rune {
		shouldCheckNext := state&checkNext != 0
		state &^= checkNext
		switch r {
		case '\n':
			state &^= isCommentLine
		case '\\':
			if state&isString != 0 {
				state |= checkNext
			}
		case '"':
			if state&isString != 0 {
				if shouldCheckNext { // escaped quote
					break // go write rune
				}
				state &^= isString
			} else if state&(isCommentLine|isCommentBlock) == 0 {
				state |= isString
			}
		case '/':
			if state&isString != 0 {
				break // go write rune
			}
			if state&isCommentBlock != 0 {
				if shouldCheckNext {
					state &^= isCommentBlock
				} else {
					state |= isCommentLine
				}
			} else {
				if shouldCheckNext {
					state |= isCommentLine
				} else {
					state |= checkNext
				}
			}
			return -1 // mark rune for skip
		case '*':
			if state&isString != 0 {
				break // go write rune
			}
			if shouldCheckNext {
				state |= isCommentBlock
			} else if state&isCommentBlock != 0 {
				state |= checkNext
			}
			return -1 // mark rune for skip
		}
		if state&(isCommentLine|isCommentBlock) != 0 {
			return -1 // mark rune for skip
		}
		return r
	}, data)
}

// Unmarshal parses the JSONC-encoded data and stores the result in the value
// pointed by v removing all comments from the data (if any).
//
// It uses [HasCommentRunes] to check whether the data contains any comment.
// Note that this operation is as expensive as the larger the data. On small
// data sets it just adds a small overhead to the unmarshaling process, but
// on large data sets it may have a significant impact on performance. In such
// cases, it may be more efficient to call [Sanitize] and then the standard
// (or any other) library directly.
//
// If the data contains comment runes, it calls [Sanitize] to remove them and
// returns [ErrInvalidUTF8] if the data is not valid UTF-8. Note that if no
// comments are found, it is assumed that the given data is valid JSON-encoded
// and the UTF-8 validity is not checked.
//
// Any error is reported from [JSONUnmarshal] as is.
//
// It uses the standard library for unmarshaling by default, but can be
// configured to use a custom Unmarshal function. For that set JSONUnmarshal
// to the desired function.
func Unmarshal(data []byte, v any) error {
	if HasCommentRunes(data) {
		var err error
		data, err = Sanitize(data)
		if err != nil {
			return err
		}
	}
	return JSONUnmarshal(data, v)
}

// HasCommentRunes returns true if the data contains any comment rune.
// It checks whether the data contains any '/' character, and if so, it looks
// whether the previous one is a '/' or the next one is a '/' or a '*'.
// If not, it returns false.
//
// Caveat: if the data contains a string that looks like a comment as
// '{"url": "//"}', HasCommentRunes returns true.
func HasCommentRunes(data []byte) bool {
	var state byte
	bytes.IndexFunc(data, func(r rune) bool {
		if state&checkNext != 0 {
			if r == '/' || r == '*' {
				state |= hasCommentRunes
				return true
			}
			state &^= checkNext
		}
		if r == '/' {
			state |= checkNext
		}
		return false
	})
	return state&hasCommentRunes != 0
}
