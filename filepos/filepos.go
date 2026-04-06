// Package filepos provides the routines to work with file positions.
// It's useful to make parser error reporting richer and human-readable.
package filepos

import "bytes"

// LinePosition stores human-readable position of a file offset.
type LinePosition struct {
	ByteOffset int    // Data offset in bytes (0-based)
	LineNumber int    // Number of \n-separated line of the offset (1-based)
	Column     int    // Byte offset inside the line (1-based)
	Line       string // Line content at the offset
}

// ByteOffsetToLinePosition returns the line number (1-based), column (0-based),
// and the full line content at the given byte offset in data. If offset is out
// of range, the function returns the line position of the minimum (if negative)
// or maximum offset in the provided data. If the provided data is of zero
// length, the function returns line number and column equal to one, and empty
// line.
func ByteOffsetToLinePosition(data []byte, offset int) LinePosition {
	// Handle edge case when data length is zero.
	if len(data) == 0 {
		return LinePosition{
			ByteOffset: offset,
			LineNumber: 1,
			Column:     1,
			Line:       "",
		}
	}

	// Make sure offset stays within the data bounds.
	if offset < 0 {
		offset = 0
	}
	if offset > len(data) {
		offset = len(data) - 1
	}

	// Scan bytes up to offset to count newlines and
	// locate the start of the line.
	lineNum := 1
	lineStart := 0
	for i := range offset {
		if data[i] == '\n' {
			lineNum++
			lineStart = i + 1
		}
	}

	// Column is byte offset minus the start-of-line index.
	column := offset - lineStart

	// Find the end of line (next '\n') after lineStart.
	rest := data[lineStart:]
	nlPos := bytes.IndexByte(rest, '\n')
	var lineEnd int
	if nlPos >= 0 {
		lineEnd = lineStart + nlPos
	} else {
		lineEnd = len(data)
	}

	lineText := string(data[lineStart:lineEnd])

	return LinePosition{
		ByteOffset: offset,
		LineNumber: lineNum,
		Column:     column,
		Line:       lineText,
	}
}
