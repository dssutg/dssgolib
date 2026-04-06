// Package leb128 provides both encoder and decoder of LEB128 format, which is
// variable-integer binary representation, similar to UTF-8.
//
// More information about LEB128: https://en.wikipedia.org/wiki/LEB128
package leb128

import (
	"errors"
	"io"
)

var (
	ErrUnexpectedEOF    = errors.New("unexpected EOF: incomplete LEB128 sequence")
	ErrFailedToReadByte = errors.New("failed to read byte")
	ErrDecodingOverflow = errors.New("LEB128 decoding overflow")
)

// WriteULEB128 encodes a uint64 value into unsigned LEB128 format and writes it
// to the given io.Writer.
func WriteULEB128(w io.Writer, value uint64) error {
	for {
		byteValue := byte(value & 0x7F) // Get the last 7 bits
		value >>= 7                     // Shift right by 7 bits
		if value != 0 {
			byteValue |= 0x80 // Set the continuation bit
		}
		if _, err := w.Write([]byte{byteValue}); err != nil {
			return err // Return any write error
		}
		if value == 0 {
			break
		}
	}

	return nil
}

// ReadULEB128 decodes unsigned LEB128 from the provided io.Reader and returns
// the uint64 value.
func ReadULEB128(r io.Reader) (uint64, error) {
	var result uint64
	var shift uint
	buf := make([]byte, 1)

	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return 0, ErrUnexpectedEOF
			}
			return 0, err
		}
		if n != 1 {
			return 0, ErrFailedToReadByte
		}

		b := buf[0]
		result |= (uint64(b&0x7F) << shift)

		// If the continuation bit is not set, we are done.
		if b&0x80 == 0 {
			break
		}

		shift += 7
		if shift >= 64 {
			return 0, ErrDecodingOverflow
		}
	}

	return result, nil
}
