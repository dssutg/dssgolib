package tcpreader

import (
	"bytes"
	"errors"
	"io"
	"net"
)

type MessageReader struct {
	conn     net.Conn     // Connection
	delim    []byte       // Delimeter bytes
	buf      bytes.Buffer // Buffer to accumulate chunks
	chunkBuf []byte       // Buffer that holds one chunk for each read call
}

// MakeMessageReader returns a new MessageReader with the provided connection
// and delimeter bytes.
func MakeMessageReader(conn net.Conn, delim []byte) MessageReader {
	return MessageReader{
		conn:     conn,
		delim:    delim,
		chunkBuf: make([]byte, 1024),
	}
}

var ErrIncompleteMessage = errors.New("incomplete message")

// ReadMessage reads from the connection until it encounters the delimiter.
// It returns the message (not including the delimiter) or an error.
func (mr *MessageReader) ReadMessage() ([]byte, error) {
	for {
		// Check if the buffer contains the delimiter.
		data := mr.buf.Bytes()
		if index := bytes.Index(data, mr.delim); index != -1 {
			// Copy the message before the delimiter.
			message := make([]byte, index)
			copy(message, data[:index])
			// Remove the processed message and the delimiter from the buffer.
			mr.buf.Next(index + len(mr.delim))
			return message, nil
		}

		// Delimeter not found, read more data from the connection.
		n, err := mr.conn.Read(mr.chunkBuf)
		if n > 0 {
			mr.buf.Write(mr.chunkBuf[:n])
		}
		if err == nil {
			continue
		}
		length := mr.buf.Len()
		mr.buf.Reset()
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
		if length != 0 {
			return nil, ErrIncompleteMessage
		}
		return nil, io.EOF
	}
}
