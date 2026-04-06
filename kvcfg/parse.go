// package kvcfg implements the parser for the kv-format.
//
// The kv-format is a simple key-value file format like this:
//
//	# Player properties
//	name = John
//	age  = 23
//
//	# Navigation
//	KEY_UP    = moveUp
//	KEY_DOWN  = moveDown
//	KEY_LEFT  = moveLeft
//	KEY_RIGHT = moveRight
//	KEY_SPACE = shoot
//
//	...
//
// Each line contains a key-value pair separated by the equals (=) sign. Leading
// and trailing spaces in each line are ignored, as well as the spaces around
// the equals sign. Blank lines and comments are ignored. Comments start with
// the hash (#) sign, and can have any characters until the end of line.
// Comments can only be alone on a separate line.
//
// The semantics of both the key and value are not defined because it is meant
// to be up to the user of this package.
//
// The origin and one of the applications of this format is to store
// the localization strings. This is so common that the handy [ParseToMap]
// function is written exactly for that.
package kvcfg

import "strings"

// ParseToMap returns the parsed key-value map parsed from the provided data in
// kv-format, or an error, if any. The semantics of the both key and value are
// not defined, i.e., the function simply fills the result map with they
// key-value pairs. This can be useful for localization strings.
func ParseToMap(data []byte) (map[string]string, error) {
	pairs := make(map[string]string)

	err := ParseWithCallback(data, func(key, value string) error {
		pairs[key] = value
		return nil
	})

	return pairs, err
}

// ParseWithCallback parses the data in the kv-format, and calls handlePair to
// handle the parsed pair. If the handler callback returns an error, the parser
// stops and wraps this error within [ParseError] along with the line number in
// the data where the error occurred.
func ParseWithCallback(data []byte, handlePair func(key, value string) error) error {
	lineNo := 1

	for line := range strings.Lines(string(data)) {
		// Ignore leading and trailing whitespace.
		line = strings.TrimSpace(line)

		// Ignore empty line or comments.
		if line == "" || strings.HasPrefix(line, "#") {
			lineNo++
			continue
		}

		// Parse at most 2 =-separated parts in this line.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Only key value pairs expected.
			return &ParseError{
				LineNo: lineNo,
				Msg:    "expected key-value pair separated by =",
			}
		}

		// Ignore spaces around the equal sign in both key and value.
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Call the handler callback with the parsed key-value pair.
		if err := handlePair(key, value); err != nil {
			return &ParseError{
				LineNo: lineNo,
				Msg:    err.Error(),
				Err:    err,
			}
		}

		lineNo++
	}

	return nil
}
