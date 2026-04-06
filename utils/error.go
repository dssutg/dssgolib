// Helper functions for error handling and runtime assertions that simplify
// common boilerplate patterns. It is encouraged to stick to common Go patterns
// to handle errors gracefully but these functions are useful for quick
// throwaway code and debugging, or internal code that just needs to be checked
// for errors and is not intended to provide pretty error handling for end
// users.
package utils

import "errors"

// Must takes an error. If err is non-nil, Must panics with that error.
// Otherwise it returns. This is useful for situations
// where you know an error cannot occur (or you want the program to abort on
// error) and don't want to write explicit error checks at each call site.
//
// Example:
//
//	Must(initAssets("assets.pak"))
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 takes a value of any type and an error. If err is non-nil, Must2 panics
// with that error. Otherwise it returns value. This is useful for situations
// where you know an error cannot occur (or you want the program to abort on
// error) and don't want to write explicit error checks at each call site.
//
// Example:
//
//	config := Must2(loadConfig("config.json"))
func Must2[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Assert panics if cond is false. If msg is non-empty, the panic message will
// include that string. Otherwise it panics with "assertion failed". Use Assert
// to enforce invariants that should never occur in correct code.
//
// Example:
//
//	Assert(len(items) > 0, "items must not be empty")
func Assert(cond bool, msg string) {
	if cond {
		return
	}

	if msg == "" {
		panic("assertion failed")
	}

	panic("assertion failed: " + msg)
}

// ErrorInSet reports whether the provided error
// is one of the errors in the set. Each individual
// target error is checked with [errors.Is].
func ErrorInSet(errs []error, err error) bool {
	for _, target := range errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
