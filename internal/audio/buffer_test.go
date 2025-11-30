package audio

import (
	"testing"
)

func TestRingBuffer(t *testing.T) {
	t.Run("Write and Read", func(t *testing.T) {
		rb := NewRingBuffer(10)
		input := []int16{1, 2, 3, 4, 5}

		err := rb.Write(input)
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}

		if rb.Available() != 5 {
			t.Errorf("Expected 5 available, got %d", rb.Available())
		}

		output, err := rb.Read(5)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		for i, v := range output {
			if v != input[i] {
				t.Errorf("Index %d: expected %d, got %d", i, input[i], v)
			}
		}

		if rb.Available() != 0 {
			t.Errorf("Expected 0 available, got %d", rb.Available())
		}
	})

	t.Run("Buffer Full", func(t *testing.T) {
		rb := NewRingBuffer(5)
		input := []int16{1, 2, 3, 4, 5}
		rb.Write(input)

		err := rb.Write([]int16{6})
		if err != ErrBufferFull {
			t.Errorf("Expected ErrBufferFull, got %v", err)
		}
	})

	t.Run("WriteForce Overwrite", func(t *testing.T) {
		rb := NewRingBuffer(5)
		rb.WriteForce([]int16{1, 2, 3, 4, 5})

		// Overwrite first 2 elements
		rb.WriteForce([]int16{6, 7})

		if rb.Available() != 5 {
			t.Errorf("Expected 5 available, got %d", rb.Available())
		}

		output, _ := rb.Read(5)
		expected := []int16{3, 4, 5, 6, 7}

		for i, v := range output {
			if v != expected[i] {
				t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	})

	t.Run("Wrap Around", func(t *testing.T) {
		rb := NewRingBuffer(5)
		rb.Write([]int16{1, 2, 3})
		rb.Read(2) // Read 1, 2. Buffer has {3} at index 2. Head at 3, Tail at 2.

		// Write 4, 5, 6, 7. Should wrap around.
		// Buffer: [6, 7, 3, 4, 5]
		err := rb.Write([]int16{4, 5, 6, 7})
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}

		output, _ := rb.Read(5)
		expected := []int16{3, 4, 5, 6, 7}

		for i, v := range output {
			if v != expected[i] {
				t.Errorf("Index %d: expected %d, got %d", i, expected[i], v)
			}
		}
	})
}
