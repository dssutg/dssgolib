// Package ring implements simple Ring Buffer collection.
package ring

// Buffer is a generic ring buffer (circular buffer) that holds elements of any type.
// It supports adding elements and retrieving the stored elements in order.
type Buffer[T any] struct {
	buffer []T // buffer holds the underlying slice of elements.
	head   int // head points to the position of the oldest element in the buffer.
	tail   int // tail points to the next position where a new element will be inserted.

	// size tracks the current number of elements in the buffer.
	// (Note: in the current implementation, this field is defined but never updated.)
	size int

	capacity int  // capacity is the size of the buffer.
	isFull   bool // isFull indicates whether the buffer is full.
}

// New creates and returns a pointer to a new Buffer with the specified capacity.
func New[T any](capacity int) *Buffer[T] {
	return &Buffer[T]{
		buffer:   make([]T, capacity),
		capacity: capacity,
	}
}

// Add inserts an element into the ring buffer.
// When the buffer is full, it will overwrite the oldest value.
func (b *Buffer[T]) Add(element T) {
	// Insert the element at the tail position.
	b.buffer[b.tail] = element

	// If the buffer is full, move the head pointer forward to overwrite the oldest
	// element.
	if b.isFull {
		b.head = (b.head + 1) % b.capacity
	}

	// Move the tail pointer forward.
	b.tail = (b.tail + 1) % b.capacity

	// Check if the tail has wrapped around to the head.
	// If so, it means the buffer has reached its capacity.
	if b.tail == b.head {
		b.isFull = true
	}
}

// GetAll returns a slice containing all the elements in the buffer in proper.
// order If the buffer is not full, it returns the slice between head and tail.
// If the buffer is full, it returns the slice starting from head to the end.
// and then from the beginning up to tail .
func (b *Buffer[T]) GetAll() []T {
	// Initialize a result slice with a capacity of b.size. Note: b.size is not
	// being updated in the Add method; consider calculating the size dynamically.
	result := make([]T, 0, b.size)

	if !b.isFull {
		// If the buffer is not full, simply return the slice from head to tail.
		result = append(result, b.buffer[b.head:b.tail]...)
	} else {
		// If the buffer is full, the stored data is wrapped.
		// First, append elements from head to the end.
		result = append(result, b.buffer[b.head:]...)
		// Then, append elements from the beginning to the tail.
		result = append(result, b.buffer[:b.tail]...)
	}

	return result
}
