package audio

import (
	"errors"
	"sync"
)

var (
	ErrBufferFull  = errors.New("ring buffer is full")
	ErrBufferEmpty = errors.New("ring buffer is empty")
)

// RingBuffer implements a thread-safe circular buffer for int16 audio samples
type RingBuffer struct {
	data     []int16
	size     int
	head     int // Write index
	tail     int // Read index
	count    int // Number of available samples
	capacity int
	mu       sync.Mutex
}

// NewRingBuffer creates a new ring buffer with specified capacity
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data:     make([]int16, capacity),
		capacity: capacity,
	}
}

// Write writes samples to the buffer
// If buffer is full, it returns ErrBufferFull and writes nothing
func (rb *RingBuffer) Write(samples []int16) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	n := len(samples)
	if rb.count+n > rb.capacity {
		return ErrBufferFull
	}

	// Write data
	for i := 0; i < n; i++ {
		rb.data[rb.head] = samples[i]
		rb.head = (rb.head + 1) % rb.capacity
		rb.count++
	}

	return nil
}

// WriteForce writes samples to the buffer, overwriting oldest data if full
func (rb *RingBuffer) WriteForce(samples []int16) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	n := len(samples)

	// If input is larger than capacity, only take the last capacity samples
	if n > rb.capacity {
		samples = samples[n-rb.capacity:]
		n = rb.capacity
	}

	for i := 0; i < n; i++ {
		rb.data[rb.head] = samples[i]

		// If full, advance tail (overwrite)
		if rb.count == rb.capacity {
			rb.tail = (rb.tail + 1) % rb.capacity
		} else {
			rb.count++
		}

		rb.head = (rb.head + 1) % rb.capacity
	}
}

// Read reads n samples from the buffer
func (rb *RingBuffer) Read(n int) ([]int16, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.count < n {
		return nil, ErrBufferEmpty
	}

	out := make([]int16, n)
	for i := 0; i < n; i++ {
		out[i] = rb.data[rb.tail]
		rb.tail = (rb.tail + 1) % rb.capacity
		rb.count--
	}

	return out, nil
}

// Available returns the number of samples available to read
func (rb *RingBuffer) Available() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.count
}

// Capacity returns the total capacity of the buffer
func (rb *RingBuffer) Capacity() int {
	return rb.capacity
}

// Reset clears the buffer
func (rb *RingBuffer) Reset() {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.head = 0
	rb.tail = 0
	rb.count = 0
}
