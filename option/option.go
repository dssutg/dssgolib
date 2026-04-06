// Package option provides a generic optional-value type. It distinguishes
// between a "zero" value and "no value set", avoiding ambiguity when T's zero
// value is a valid value.
package option

// Option[T] represents a value of type T that may or may not be present. Use
// Some(v) to construct an instance that holds v, and None[T]() for an empty
// instance. You can test presence with OK(), and retrieve the contained value
// with Unwrap(), which panics if no value is set.
// Zero value represents the option without any assigned value.
type Option[T any] struct {
	hasValue bool // True if value is present
	value    T    // The contained value, if hasValue is true
}

// Some returns an Option[T] containing the given value The result .OK() will v.
// be true, and Unwrap() will return v.
func Some[T any](v T) Option[T] {
	return Option[T]{hasValue: true, value: v}
}

// None returns an empty Option[T]. The result .OK() will be false, and Unwrap()
// will panic if called.
func None[T any]() Option[T] {
	return Option[T]{}
}

// OK reports whether this Option contains a value. It returns true for a value
// constructed via Some, and false for None or zero-initialized Option[T].
func (opt Option[T]) OK() bool {
	return opt.hasValue
}

// Unwrap returns the contained value if OK() is true. If OK() is false (the
// Option is empty), Unwrap panics Use OK() to guard calls to Unwrap.
func (opt Option[T]) Unwrap() T {
	if !opt.hasValue {
		panic("option does not have value")
	}
	return opt.value
}

// Coalesce returns the contained value if it is present, otherwise
// (if the option is empty), the provided default value is returned.
func (opt Option[T]) Coalesce(defaultValue T) T {
	if !opt.hasValue {
		return defaultValue
	}
	return opt.value
}

// Get returns the pair of the value itself and whether it is present.
func (opt Option[T]) Get() (value T, ok bool) {
	return opt.value, opt.hasValue
}

// Cast returns a new option with the conversion done by f callback, but
// ff the option is empty, the new option is also empty.
func Cast[Src, Dst any](opt Option[Src], f func(v Src) Dst) Option[Dst] {
	if !opt.hasValue {
		return Option[Dst]{}
	}
	return Option[Dst]{hasValue: true, value: f(opt.value)}
}
