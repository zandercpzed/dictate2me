package transcription

// Transcriber defines the methods required by consumers of a transcription
// engine. This allows using the real Vosk-backed Engine or a lightweight
// noop implementation when the native dependency is unavailable.
type Transcriber interface {
    Reset() error
    TranscribeStream([]int16) ([]Segment, error)
    PartialResult() (string, error)
    FinalResult() (Segment, error)
    Close() error
}

// noopEngine is a minimal transcription engine used as a degraded-mode
// fallback when the real native engine cannot be initialized.
type noopEngine struct{}

func (n *noopEngine) Reset() error {
    return nil
}

func (n *noopEngine) TranscribeStream(samples []int16) ([]Segment, error) {
    return []Segment{}, nil
}

func (n *noopEngine) PartialResult() (string, error) {
    return "", nil
}

func (n *noopEngine) FinalResult() (Segment, error) {
    return Segment{}, nil
}

func (n *noopEngine) Close() error {
    return nil
}

// NewNoopEngine returns a Transcriber that performs no actual
// transcription. It's safe to use when Vosk is not installed or fails
// to initialize; the rest of the daemon will function in degraded mode.
func NewNoopEngine() Transcriber {
    return &noopEngine{}
}
