// Package audio provides audio capture and processing capabilities.
//
// This package handles:
//   - Real-time audio capture from the microphone via PortAudio
//   - Circular buffer management for efficient audio streaming
//   - Voice Activity Detection (VAD) to detect speech segments
//   - Audio format conversion (to WAV 16kHz mono 16-bit)
//
// # Architecture
//
// The Capture struct manages the PortAudio stream and feeds data into a channel.
// A RingBuffer is provided for scenarios requiring a sliding window of audio history.
// VAD is implemented using a simple energy-based algorithm with hysteresis.
//
// # Usage
//
//	// Configure capture
//	cfg := audio.DefaultConfig()
//	cfg.SampleRate = 16000
//
//	// Create capture instance
//	capture, err := audio.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer capture.Close()
//
//	// Start capturing
//	if err := capture.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Consume audio stream
//	for chunk := range capture.Stream() {
//	    // Process audio chunk (e.g., feed to VAD or Whisper)
//	    if vad.Process(chunk) {
//	        // Speech detected
//	    }
//	}
package audio
