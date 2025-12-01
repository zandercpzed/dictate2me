package transcription

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testModelPath = "../../models/vosk-model-small-pt-0.3"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				ModelPath:  testModelPath,
				SampleRate: 16000,
				Language:   "pt",
			},
			wantErr: false,
		},
		{
			name: "missing model path",
			config: Config{
				SampleRate: 16000,
				Language:   "pt",
			},
			wantErr: true,
		},
		{
			name: "non-existent model",
			config: Config{
				ModelPath:  "/nonexistent/path",
				SampleRate: 16000,
				Language:   "pt",
			},
			wantErr: true,
		},
		{
			name: "default values",
			config: Config{
				ModelPath: testModelPath,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid config" || tt.name == "default values" {
				if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
					t.Skip("Skipping test, model not found:", testModelPath)
				}
			}

			engine, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, engine)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, engine)
				if engine != nil {
					defer engine.Close()
					assert.Equal(t, tt.config.ModelPath, engine.config.ModelPath)
					if tt.config.SampleRate == 0 {
						assert.Equal(t, 16000.0, engine.config.SampleRate)
					}
					if tt.config.Language == "" {
						assert.Equal(t, "pt", engine.config.Language)
					}
				}
			}
		})
	}
}

func TestTranscribeStream(t *testing.T) {
	if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test, model not found:", testModelPath)
	}

	engine, err := New(Config{
		ModelPath:  testModelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	require.NoError(t, err)
	defer engine.Close()

	tests := []struct {
		name    string
		samples []int16
		wantErr bool
	}{
		{
			name:    "empty samples",
			samples: []int16{},
			wantErr: false,
		},
		{
			name:    "silence samples",
			samples: make([]int16, 16000), // 1 second of silence
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments, err := engine.TranscribeStream(tt.samples)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, segments)
			}
		})
	}
}

func TestPartialResult(t *testing.T) {
	if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test, model not found:", testModelPath)
	}

	engine, err := New(Config{
		ModelPath:  testModelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	require.NoError(t, err)
	defer engine.Close()

	// Send some silence
	samples := make([]int16, 8000) // 0.5 seconds
	_, err = engine.TranscribeStream(samples)
	require.NoError(t, err)

	// Get partial result
	partial, err := engine.PartialResult()
	assert.NoError(t, err)
	assert.NotNil(t, partial)
}

func TestFinalResult(t *testing.T) {
	if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test, model not found:", testModelPath)
	}

	engine, err := New(Config{
		ModelPath:  testModelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	require.NoError(t, err)
	defer engine.Close()

	// Send some samples
	samples := make([]int16, 16000) // 1 second
	_, err = engine.TranscribeStream(samples)
	require.NoError(t, err)

	// Get final result
	segment, err := engine.FinalResult()
	assert.NoError(t, err)
	assert.NotNil(t, segment)
}

func TestReset(t *testing.T) {
	if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test, model not found:", testModelPath)
	}

	engine, err := New(Config{
		ModelPath:  testModelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	require.NoError(t, err)
	defer engine.Close()

	// Send some samples
	samples := make([]int16, 8000)
	_, err = engine.TranscribeStream(samples)
	require.NoError(t, err)

	// Reset
	err = engine.Reset()
	assert.NoError(t, err)

	// Should be able to transcribe again
	_, err = engine.TranscribeStream(samples)
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {
	if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
		t.Skip("Skipping test, model not found:", testModelPath)
	}

	engine, err := New(Config{
		ModelPath:  testModelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	require.NoError(t, err)

	// Close once
	err = engine.Close()
	assert.NoError(t, err)

	// Close again should not error
	err = engine.Close()
	assert.NoError(t, err)

	// Operations after close should error
	_, err = engine.TranscribeStream([]int16{1, 2, 3})
	assert.Error(t, err)

	_, err = engine.PartialResult()
	assert.Error(t, err)

	err = engine.Reset()
	assert.Error(t, err)
}

func TestInt16ToBytes(t *testing.T) {
	tests := []struct {
		name     string
		samples  []int16
		expected []byte
	}{
		{
			name:     "empty",
			samples:  []int16{},
			expected: []byte{},
		},
		{
			name:     "single sample",
			samples:  []int16{256},
			expected: []byte{0, 1},
		},
		{
			name:     "multiple samples",
			samples:  []int16{0, 256, -1},
			expected: []byte{0, 0, 0, 1, 255, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int16ToBytes(tt.samples)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func BenchmarkTranscribeStream(b *testing.B) {
	if _, err := os.Stat(testModelPath); os.IsNotExist(err) {
		b.Skip("Skipping benchmark, model not found:", testModelPath)
	}

	engine, err := New(Config{
		ModelPath:  testModelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	require.NoError(b, err)
	defer engine.Close()

	samples := make([]int16, 8000) // 0.5 seconds

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.TranscribeStream(samples)
		if err != nil {
			b.Fatal(err)
		}
	}
}
