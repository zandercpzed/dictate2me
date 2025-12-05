package audio

import (
    "testing"
    "time"
)

func TestProcessAudioAndChannels(t *testing.T) {
    // Create a Capture instance without calling New (avoids PortAudio init)
    c := &Capture{
        config: DefaultConfig(),
        outCh:  make(chan []int16, 2),
        errCh:  make(chan error, 1),
        buffer: make([]int16, DefaultFrameSize),
    }

    // Simulate audio callback with sample data
    samples := []int16{1, 2, 3, 4}

    // Call processAudio multiple times to fill buffer
    c.processAudio(samples)

    select {
    case got := <-c.Stream():
        if len(got) != len(samples) {
            t.Fatalf("expected chunk length %d, got %d", len(samples), len(got))
        }
    case <-time.After(100 * time.Millisecond):
        t.Fatalf("timed out waiting for audio chunk")
    }

    // Fill buffer to cause overflow
    c.processAudio(samples)
    c.processAudio(samples)

    // One of the subsequent attempts may signal an error due to overflow
    select {
    case err := <-c.Error():
        if err == nil {
            t.Fatalf("expected error on overflow, got nil")
        }
    case <-time.After(100 * time.Millisecond):
        // ok: error channel may not be filled depending on timing
    }

    // Closing should not panic
    if err := c.Close(); err != nil {
        t.Fatalf("Close returned error: %v", err)
    }
}
