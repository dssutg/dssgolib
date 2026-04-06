package kvcfg

import "fmt"

// ParseError is the error returned by the kv-parser.
type ParseError struct {
	LineNo int    // line number in the key-value input with the error
	Msg    string // error message
	Err    error  // Wrapped err or nil of none
}

// Error implements the error interface for [ParseError].
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d: %s", e.LineNo, e.Msg)
}

// Unwrap unwraps the original error or nil of none.
// It implements the interface used by the [errors.Is].
func (e *ParseError) Unwrap() error {
	return e.Err
}
