package option

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// MarshalJSON implements json.Marshaler interface.
func (opt Option[T]) MarshalJSON() ([]byte, error) {
	if !opt.hasValue {
		return []byte("null"), nil
	}
	return json.Marshal(opt.value)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (opt *Option[T]) UnmarshalJSON(p []byte) error {
	// If JSON is "null", set to None.
	if string(p) == "null" {
		*opt = None[T]()
		return nil
	}

	// Ensure opt is non-nil receiver.
	if opt == nil {
		return nil
	}

	// Allocate a temp variable of type T and unmarshal into it.
	var v T
	if err := json.Unmarshal(p, &v); err != nil {
		// On error, leave opt unchanged.
		return err
	}
	*opt = Some(v)
	return nil
}

// IsZero reports whether the Option is empty.
// The encoding/json package uses this to support
// the omitzero struct tag.
func (opt Option[T]) IsZero() bool {
	return !opt.hasValue
}

// MarshalBinary implements encoding.BinaryMarshaler for gob.
func (opt Option[T]) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	// Encode presence first.
	if err := enc.Encode(opt.hasValue); err != nil {
		return nil, err
	}
	if !opt.hasValue {
		return buf.Bytes(), nil
	}

	// Encode the value.
	if err := enc.Encode(opt.value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for gob.
func (opt *Option[T]) UnmarshalBinary(data []byte) error {
	// Reset receiver.
	*opt = None[T]()

	if len(data) == 0 {
		// Empty input -> None.
		return nil
	}

	dec := gob.NewDecoder(bytes.NewReader(data))

	var hasValue bool
	if err := dec.Decode(&hasValue); err != nil {
		return err
	}

	if !hasValue {
		// Keep assigned None.
		return nil
	}

	var v T
	if err := dec.Decode(&v); err != nil {
		return err
	}

	*opt = Some(v)

	return nil
}
