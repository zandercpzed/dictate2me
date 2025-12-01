package transcription

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	vosk "github.com/alphacep/vosk-api/go"
)

// Config holds configuration for the transcription engine
type Config struct {
	ModelPath  string
	SampleRate float64 // 16000 Hz recommended
	Language   string  // "pt", "en", etc.
}

// Engine handles speech-to-text transcription using Vosk
type Engine struct {
	model      *vosk.VoskModel
	recognizer *vosk.VoskRecognizer
	config     Config
	mu         sync.Mutex
	closed     bool
}

// Segment represents a transcribed text segment
type Segment struct {
	Text  string
	Start time.Duration
	End   time.Duration
	Conf  float32 // Confidence score
}

// VoskResult represents the JSON result from Vosk
type VoskResult struct {
	Text   string         `json:"text"`
	Result []VoskWordInfo `json:"result,omitempty"`
}

// VoskWordInfo represents word-level information
type VoskWordInfo struct {
	Conf  float32 `json:"conf"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Word  string  `json:"word"`
}

// New creates a new transcription engine using Vosk
func New(cfg Config) (*Engine, error) {
	if cfg.ModelPath == "" {
		return nil, fmt.Errorf("model path is required")
	}
	if _, err := os.Stat(cfg.ModelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model directory not found: %s", cfg.ModelPath)
	}
	if cfg.Language == "" {
		cfg.Language = "pt"
	}
	if cfg.SampleRate == 0 {
		cfg.SampleRate = 16000.0
	}

	// Set Vosk log level (0 = no logs, -1 = all logs)
	vosk.SetLogLevel(0)

	// Load model
	model, err := vosk.NewModel(cfg.ModelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load vosk model from %s: %w", cfg.ModelPath, err)
	}

	// Create recognizer
	recognizer, err := vosk.NewRecognizer(model, cfg.SampleRate)
	if err != nil {
		model.Free()
		return nil, fmt.Errorf("failed to create vosk recognizer: %w", err)
	}

	// Enable word-level timestamps
	recognizer.SetWords(1)

	return &Engine{
		model:      model,
		recognizer: recognizer,
		config:     cfg,
	}, nil
}

// TranscribeStream processes audio samples in streaming mode
// Audio must be 16kHz mono int16 samples
func (e *Engine) TranscribeStream(samples []int16) ([]Segment, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil, fmt.Errorf("engine is closed")
	}

	if len(samples) == 0 {
		return []Segment{}, nil
	}

	// Convert int16 samples to bytes
	data := int16ToBytes(samples)

	// Accept waveform
	state := e.recognizer.AcceptWaveform(data)

	segments := []Segment{}

	if state > 0 {
		// Final result available
		result := e.recognizer.Result()
		seg, err := parseVoskResult(result)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vosk result: %w", err)
		}
		if seg.Text != "" {
			segments = append(segments, seg)
		}
	}

	return segments, nil
}

// FinalResult gets the final transcription result
// Call this after all audio has been processed
func (e *Engine) FinalResult() (Segment, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return Segment{}, fmt.Errorf("engine is closed")
	}

	result := e.recognizer.FinalResult()
	return parseVoskResult(result)
}

// PartialResult gets partial transcription result
// Useful for real-time feedback
func (e *Engine) PartialResult() (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return "", fmt.Errorf("engine is closed")
	}

	result := e.recognizer.PartialResult()

	var partial struct {
		Partial string `json:"partial"`
	}

	if err := json.Unmarshal([]byte(result), &partial); err != nil {
		return "", fmt.Errorf("failed to parse partial result: %w", err)
	}

	return partial.Partial, nil
}

// Reset resets the recognizer state
func (e *Engine) Reset() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return fmt.Errorf("engine is closed")
	}

	e.recognizer.Reset()
	return nil
}

// Close releases resources
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	if e.recognizer != nil {
		e.recognizer.Free()
		e.recognizer = nil
	}

	if e.model != nil {
		e.model.Free()
		e.model = nil
	}

	e.closed = true
	return nil
}

// parseVoskResult parses Vosk JSON result into Segment
func parseVoskResult(jsonResult string) (Segment, error) {
	var result VoskResult
	if err := json.Unmarshal([]byte(jsonResult), &result); err != nil {
		return Segment{}, fmt.Errorf("failed to unmarshal vosk result: %w", err)
	}

	if result.Text == "" {
		return Segment{}, nil
	}

	seg := Segment{
		Text: result.Text,
	}

	// Calculate timing and confidence from word-level info
	if len(result.Result) > 0 {
		seg.Start = time.Duration(result.Result[0].Start * float64(time.Second))
		seg.End = time.Duration(result.Result[len(result.Result)-1].End * float64(time.Second))

		// Average confidence
		var totalConf float32
		for _, word := range result.Result {
			totalConf += word.Conf
		}
		seg.Conf = totalConf / float32(len(result.Result))
	}

	return seg, nil
}

// int16ToBytes converts int16 samples to byte array
func int16ToBytes(samples []int16) []byte {
	data := make([]byte, len(samples)*2)
	for i, sample := range samples {
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}
	return data
}
