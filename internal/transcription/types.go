package transcription

import "context"

// Segment represents a transcribed text segment with confidence score
type Segment struct {
	Text string  // Transcribed text
	Conf float32 // Confidence score (0.0 to 1.0)
}

// Transcriber is the interface for all transcription engines
type Transcriber interface {
	// TranscribeStream processes audio samples and returns segments
	TranscribeStream(samples []int16) ([]Segment, error)

	// PartialResult returns partial/interim transcription (if supported)
	PartialResult() (string, error)

	// FinalResult gets the final transcription result
	FinalResult() (Segment, error)

	// Reset resets the transcription state
	Reset() error

	// Close releases resources
	Close() error
}

// TranscriberWithContext extends Transcriber with context support
type TranscriberWithContext interface {
	Transcriber
	TranscribeWithContext(ctx context.Context, samples []int16) (string, error)
}
