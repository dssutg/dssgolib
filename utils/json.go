package utils

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// UnmarshalJSON is a type-safe wrapper around [json.Unmarshal] (jsonv2).
// It unmarshals JSON bytes into a value of type T.
func UnmarshalJSON[T any](in []byte, out *T, opts ...json.Options) error {
	return json.Unmarshal(in, out, opts...)
}

// UnmarshalDecodeJSON is a type-safe wrapper around [json.UnmarshalDecode].
// It unmarshals JSON from a decoder into a value of type T
func UnmarshalDecodeJSON[T any](in *jsontext.Decoder, out *T, opts ...json.Options) error {
	return json.UnmarshalDecode(in, out, opts...)
}
