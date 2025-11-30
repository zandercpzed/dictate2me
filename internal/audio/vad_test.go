package audio

import (
	"testing"
	"time"
)

func TestVAD(t *testing.T) {
	sampleRate := 16000
	vad := NewVAD(sampleRate)

	// Configure for fast testing
	vad.minSpeechDur = 10 * time.Millisecond
	vad.minSilenceDur = 10 * time.Millisecond
	vad.SetThreshold(100.0)

	t.Run("Silence Detection", func(t *testing.T) {
		vad.Reset()
		silence := make([]int16, 1024) // All zeros

		if vad.Process(silence) {
			t.Error("Expected silence, got speech")
		}
	})

	t.Run("Speech Detection", func(t *testing.T) {
		vad.Reset()
		// Generate sine wave (loud)
		speech := make([]int16, 1024)
		for i := range speech {
			speech[i] = 1000 // Above threshold 100
		}

		// Feed enough chunks to trigger minSpeechDur
		// 1024 samples @ 16kHz is ~64ms. minSpeechDur is 10ms.
		// One chunk should be enough.

		if !vad.Process(speech) {
			t.Error("Expected speech, got silence")
		}
	})

	t.Run("Hysteresis", func(t *testing.T) {
		vad.Reset()
		vad.minSilenceDur = 100 * time.Millisecond // Require 100ms silence to stop

		// 1. Start speech
		speech := make([]int16, 1024)
		for i := range speech {
			speech[i] = 1000
		}
		vad.Process(speech)
		if !vad.isSpeech {
			t.Fatal("Failed to detect initial speech")
		}

		// 2. Brief silence (64ms) < minSilenceDur (100ms)
		silence := make([]int16, 1024)
		if !vad.Process(silence) {
			t.Error("VAD dropped speech too early (hysteresis failed)")
		}

		// 3. More silence (another 64ms) -> Total 128ms > 100ms
		if vad.Process(silence) {
			t.Error("VAD failed to detect end of speech")
		}
	})
}
