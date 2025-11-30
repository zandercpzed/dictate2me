package audio

import (
	"fmt"
	"sync"

	"github.com/gordonklaus/portaudio"
)

// Default constants for audio capture
const (
	DefaultSampleRate = 16000
	DefaultChannels   = 1
	DefaultFrameSize  = 1024 // ~64ms at 16kHz
)

// Config holds configuration for audio capture
type Config struct {
	SampleRate int
	Channels   int
	FrameSize  int
	DeviceID   int // -1 for default device
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		SampleRate: DefaultSampleRate,
		Channels:   DefaultChannels,
		FrameSize:  DefaultFrameSize,
		DeviceID:   -1,
	}
}

// Capture handles audio capture from microphone
type Capture struct {
	config Config
	stream *portaudio.Stream
	buffer []int16
	outCh  chan []int16
	errCh  chan error

	// State management
	mu        sync.Mutex
	running   bool
	closed    bool
	initOnce  sync.Once
	closeOnce sync.Once
}

// New creates a new Capture instance
func New(cfg Config) (*Capture, error) {
	// Initialize PortAudio (safe to call multiple times if managed globally,
	// but here we assume one capture instance manages the lifecycle for simplicity
	// or we rely on a global init/terminate in main. For now, let's init here)
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize portaudio: %w", err)
	}

	if cfg.SampleRate == 0 {
		cfg.SampleRate = DefaultSampleRate
	}
	if cfg.Channels == 0 {
		cfg.Channels = DefaultChannels
	}
	if cfg.FrameSize == 0 {
		cfg.FrameSize = DefaultFrameSize
	}

	return &Capture{
		config: cfg,
		outCh:  make(chan []int16, 100), // Buffer for ~6 seconds of audio chunks
		errCh:  make(chan error, 1),
		buffer: make([]int16, cfg.FrameSize),
	}, nil
}

// Start begins audio capture
func (c *Capture) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}
	if c.closed {
		return fmt.Errorf("capture is closed")
	}

	var err error
	// Open default stream for now
	// TODO: Support selecting specific device via c.config.DeviceID
	c.stream, err = portaudio.OpenDefaultStream(
		c.config.Channels,
		0, // no output
		float64(c.config.SampleRate),
		c.config.FrameSize,
		c.processAudio,
	)
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}

	if err := c.stream.Start(); err != nil {
		c.stream.Close()
		return fmt.Errorf("failed to start stream: %w", err)
	}

	c.running = true
	return nil
}

// Stop pauses audio capture
func (c *Capture) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running || c.stream == nil {
		return nil
	}

	if err := c.stream.Stop(); err != nil {
		return fmt.Errorf("failed to stop stream: %w", err)
	}

	if err := c.stream.Close(); err != nil {
		return fmt.Errorf("failed to close stream: %w", err)
	}

	c.stream = nil
	c.running = false
	return nil
}

// Close releases resources
func (c *Capture) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	// Stop if running
	if c.running && c.stream != nil {
		_ = c.stream.Stop()
		_ = c.stream.Close()
	}

	// Terminate PortAudio
	// Note: In a real app with multiple streams, we'd need reference counting
	portaudio.Terminate()

	close(c.outCh)
	close(c.errCh)
	c.closed = true
	c.running = false

	return nil
}

// Stream returns the channel for reading audio chunks
func (c *Capture) Stream() <-chan []int16 {
	return c.outCh
}

// Error returns the channel for reading errors
func (c *Capture) Error() <-chan error {
	return c.errCh
}

// processAudio is the callback called by PortAudio
// It runs in a separate goroutine/thread managed by PortAudio
func (c *Capture) processAudio(in []int16) {
	// Create a copy of the data to send over channel
	// We must copy because 'in' is reused by PortAudio
	chunk := make([]int16, len(in))
	copy(chunk, in)

	// Non-blocking send to avoid blocking the audio callback
	select {
	case c.outCh <- chunk:
		// Success
	default:
		// Buffer full, drop chunk and notify error
		select {
		case c.errCh <- fmt.Errorf("audio buffer overflow, dropping chunk"):
		default:
			// Error channel also full, just drop
		}
	}
}
