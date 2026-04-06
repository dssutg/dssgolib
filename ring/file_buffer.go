package ring

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/dssutg/dssgolib/option"
)

// FileBuffer is a concurrent ring buffer stored in file.
// Suitable for large persistent ring buffers.
// For in-memory, refer to [Buffer].
type FileBuffer struct {
	mu       sync.RWMutex         // for concurrent access
	file     *os.File             // file handle for the buffer
	elemSize int                  // size of each element in bytes in the buffer
	maxElems int                  // maximum amount of elements in the buffer
	bias     int                  // byte offset of the start of the buffer in the file
	next     int                  // next element absolute index (does not wrap)
	nextOff  option.Option[int64] // byte offset of next element index in file
}

// OpenFileBuffer open a file name for reading and writing. The file
// is created if it does not exist. The file is never truncated.
// It returns an initialized buffer ready to be used, and an error
// if cannot open the file.
func OpenFileBuffer(name string, elemSize, maxElems, bias int) (*FileBuffer, error) {
	flag := os.O_CREATE | os.O_RDWR // DO NOT truncate

	logFile, err := os.OpenFile(name, flag, 0o644)
	if err != nil {
		return nil, err
	}

	buf := NewFileBuffer(logFile, elemSize, maxElems, bias)

	return buf, nil
}

// Close implements [io.Closer]. It closes the underlying
// file and reset the buffer to zero value.
func (b *FileBuffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.file.Close(); err != nil {
		return err
	}

	*b = FileBuffer{}

	return nil
}

// NewFileBuffer returns a new FileBuffer instance associated with the file.
// The file is assumed to be owned by the buffer, otherwise concurrent reads or writes
// are undefined. The buffer starts at the bias file byte offset with the
// total size of maxElems per elemSize bytes. This function panics if file is nil,
// elemSize or maxElems are not positive, or bias is negative.
func NewFileBuffer(file *os.File, elemSize, maxElems, bias int) *FileBuffer {
	// Validate properties.
	switch {
	case file == nil:
		panic("file is nil")
	case elemSize <= 0:
		panic("elemSize is not positive")
	case maxElems <= 0:
		panic("maxElems is not positive")
	case bias < 0:
		panic("bias is negative")
	}

	// Instantiate the buffer.
	return &FileBuffer{
		file:     file,
		elemSize: elemSize,
		maxElems: maxElems,
		bias:     bias,
	}
}

// Reset initializes the buffer with zero value.
// It must be called after all resources taken by
// the buffer are freed.
func (b *FileBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	*b = FileBuffer{}
}

// SetIndexOffset set the byte offset in file of the next index.
// When a new record is written, the next index is also
// written to the location. It panics if the offset is negative.
func (b *FileBuffer) SetIndexOffset(off int64) {
	if off < 0 {
		panic("negative next index offset")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextOff = option.Some(off)
}

// NextIndex returns the index of the next element to be added.
// The index is not wrapped along the buffer, so it can be used
// as an element ID, although it is zero-based.
func (b *FileBuffer) NextIndex() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.next
}

// SetNextIndex sets the index of the next element to the one provided.
// It panics if index is negative.
// It is intended to be used when restoring buffer state
// after initializing it.
func (b *FileBuffer) SetNextIndex(next int) {
	if next < 0 {
		panic("negative next index")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.next = next
}

// ModifyFile allows to concurrently modify the buffer's file.
// The cb is called with the pointer to the file to access.
// The error cb is returned.
func (b *FileBuffer) ModifyFile(cb func(f *os.File) error) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return cb(b.file)
}

// Sync calles file Sync method.
func (b *FileBuffer) Sync() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.file.Sync()
}

// Write implements the [io.Writer] interface. It appends
// a new record to the ring buffer, and advances the index
// of the next element only if the write has been successful.
// It panics if p length is not the exact size of element.
func (b *FileBuffer) Write(p []byte) (n int, err error) {
	if len(p) != b.elemSize {
		panic("bad p length")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// Compute file byte offset.
	off := int64((b.next%b.maxElems)*b.elemSize + b.bias)

	// Write the element bytes at the offset.
	n, err = b.file.WriteAt(p, off)
	if err != nil {
		return n, err
	}

	// If the byte offset of the next index is set, also write the next index.
	if nextOff, ok := b.nextOff.Get(); ok {
		var nextIdxBytes [8]byte
		binary.LittleEndian.PutUint64(nextIdxBytes[:], uint64(b.next+1)) // #nosec G115: no overflow
		if _, err = b.file.WriteAt(nextIdxBytes[:], nextOff); err != nil {
			return n, fmt.Errorf("cannot write next index: %w", err)
		}
	}

	// Advance the index of the next element on success.
	// Do not wrap to allow absolute IDs.
	b.next++

	return n, nil
}

// ReadLastN returns a slice of at most limit last written elements.
// It returns an error if the read has been unsuccessful.
// The result slice length is the read element count multiplied by the
// element size. This is to avoid extra allocations per elements as all
// elements are read in a single system call or two if the elements are wrapped.
func (b *FileBuffer) ReadLastN(limit int) (listBytes []byte, nextIndex int, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Determine the actual element count to read.
	nextIndex = b.next
	written := nextIndex
	count := min(limit, b.maxElems, written)

	if count == 0 {
		// No elements to read, return an empty slice.
		return
	}

	// Determine the wrapped index range for read.
	endIdx := b.next % b.maxElems             // exclusive
	startIdx := (b.next - count) % b.maxElems // inclusive

	// Allocate the result list bytes for all read elements.
	listBytes = make([]byte, count*b.elemSize)

	// The index range does not wrap.
	if startIdx < endIdx {
		// Read all records in one go.
		startOff := int64(startIdx*b.elemSize + b.bias)
		if _, err = b.file.ReadAt(listBytes, startOff); err != nil {
			listBytes = nil
		}
		return
	}

	// The index range is wrapped, so have to read in two chunks.
	// First one is read from start index to ring buffer end,
	// and second is the rest records from the start of the buffer.
	startSize := (b.maxElems - startIdx) * b.elemSize
	startOff := int64(startIdx*b.elemSize + b.bias)
	if _, err = b.file.ReadAt(listBytes[:startSize], startOff); err != nil {
		listBytes = nil
		return
	}
	endOff := int64(0*b.elemSize + b.bias)
	if _, err = b.file.ReadAt(listBytes[startSize:], endOff); err != nil {
		listBytes = nil
		return
	}

	return
}
