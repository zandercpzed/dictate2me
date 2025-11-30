// Package audio provides audio capture and processing capabilities.
//
// This package handles:
//   - Real-time audio capture from the microphone via PortAudio
//   - Circular buffer management for efficient audio streaming
//   - Voice Activity Detection (VAD) to detect speech segments
//   - Audio format conversion (to WAV 16kHz mono 16-bit)
//
// The audio module is designed to be cross-platform, with platform-specific
// implementations in platform_*.go files.
//
// Example usage:
//
//	capture, err := audio.New(
//	    audio.WithSampleRate(16000),
//	    audio.WithChannels(1),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer capture.Close()
//
//	if err := capture.Start(); err != nil {
//	    log.Fatal(err)
//	}
//
//	for segment := range capture.Segments() {
//	    // Process audio segment
//	    fmt.Printf("Got audio segment: %d samples\n", len(segment.Data))
//	}
package audio
