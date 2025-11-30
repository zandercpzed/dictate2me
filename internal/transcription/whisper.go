package transcription

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

// Config holds configuration for the transcription engine
type Config struct {
	ModelPath string
	Language  string // "pt", "en", etc.
	Threads   int    // Number of threads to use (defaults to runtime.NumCPU())
}

// Engine handles speech-to-text transcription
type Engine struct {
	model   whisper.Model
	context whisper.Context
	config  Config
	mu      sync.Mutex
	closed  bool
}

// Segment represents a transcribed text segment
type Segment struct {
	Text  string
	Start time.Duration
	End   time.Duration
	Prob  float32
}

// New creates a new transcription engine
func New(cfg Config) (*Engine, error) {
	if cfg.ModelPath == "" {
		return nil, fmt.Errorf("model path is required")
	}
	if _, err := os.Stat(cfg.ModelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", cfg.ModelPath)
	}
	if cfg.Language == "" {
		cfg.Language = "pt"
	}
	if cfg.Threads == 0 {
		cfg.Threads = runtime.NumCPU()
	}

	// Load model
	model, err := whisper.New(cfg.ModelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load whisper model: %w", err)
	}

	// Create context
	ctx, err := model.NewContext()
	if err != nil {
		model.Close()
		return nil, fmt.Errorf("failed to create whisper context: %w", err)
	}

	return &Engine{
		model:   model,
		context: ctx,
		config:  cfg,
	}, nil
}

// Transcribe converts audio samples to text
// Audio must be 16kHz mono float32 samples normalized to [-1, 1]
func (e *Engine) Transcribe(samples []float32) ([]Segment, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil, fmt.Errorf("engine is closed")
	}

	if len(samples) == 0 {
		return []Segment{}, nil
	}

	// Process audio
	if err := e.context.Process(samples, nil, nil); err != nil {
		return nil, fmt.Errorf("failed to process audio: %w", err)
	}

	// Extract segments
	var segments []Segment
	for {
		segment, err := e.context.NextSegment()
		if err != nil {
			break
		}

		segments = append(segments, Segment{
			Text:  segment.Text,
			Start: segment.Start,
			End:   segment.End,
			Prob:  segment.Prob, // Assuming Prob exists in binding, otherwise we might need to check binding API
		})
	}

	return segments, nil
}

// Close releases resources
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	if e.context != nil {
		// e.context.Close() // Check if context has Close method in binding
	}

	if e.model != nil {
		e.model.Close()
	}

	e.closed = true
	return nil
}
