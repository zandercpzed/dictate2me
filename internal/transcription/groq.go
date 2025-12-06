package transcription

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// GroqConfig holds Groq API configuration
type GroqConfig struct {
	APIKey   string
	Model    string // whisper-large-v3 (recommended)
	Language string
	Timeout  time.Duration
}

// GroqEngine implements transcription using Groq's Whisper API
type GroqEngine struct {
	config     GroqConfig
	httpClient *http.Client
}

// DefaultGroqConfig returns default configuration
func DefaultGroqConfig() GroqConfig {
	apiKey := os.Getenv("GROQ_API_KEY")
	return GroqConfig{
		APIKey:   apiKey,
		Model:    "whisper-large-v3",
		Language: "pt",
		Timeout:  30 * time.Second,
	}
}

// NewGroq creates a new Groq transcription engine
func NewGroq(cfg GroqConfig) (*GroqEngine, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable is required. Get one at: https://console.groq.com")
	}
	if cfg.Model == "" {
		cfg.Model = DefaultGroqConfig().Model
	}
	if cfg.Language == "" {
		cfg.Language = DefaultGroqConfig().Language
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultGroqConfig().Timeout
	}

	return &GroqEngine{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// TranscribeStream transcribes audio samples (streaming mode)
func (g *GroqEngine) TranscribeStream(samples []int16) ([]Segment, error) {
	text, err := g.transcribeInt16(samples)
	if err != nil {
		return nil, err
	}

	if text == "" {
		return []Segment{}, nil
	}

	return []Segment{{
		Text: text,
		Conf: 0.95, // Groq/Whisper is highly accurate
	}}, nil
}

// PartialResult returns partial transcription (not supported by Groq/Whisper)
func (g *GroqEngine) PartialResult() (string, error) {
	// Whisper doesn't support streaming partial results
	return "", nil
}

// FinalResult gets final transcription
func (g *GroqEngine) FinalResult() (Segment, error) {
	// In streaming mode, this is called when user stops
	// We don't buffer, so return empty
	return Segment{Text: "", Conf: 0.95}, nil
}

// Reset resets the transcription state
func (g *GroqEngine) Reset() error {
	// Stateless, nothing to reset
	return nil
}

// Close releases resources
func (g *GroqEngine) Close() error {
	// Nothing to clean up
	return nil
}

// transcribeInt16 sends int16 samples to Groq for transcription
func (g *GroqEngine) transcribeInt16(samples []int16) (string, error) {
	// Convert int16 to WAV bytes
	wavBytes, err := int16ToWAV(samples, 16000)
	if err != nil {
		return "", fmt.Errorf("failed to convert to WAV: %w", err)
	}

	// Encode to base64 for JSON transport (alternative: multipart/form-data)
	encodedAudio := base64.StdEncoding.EncodeToString(wavBytes)

	// Prepare request
	reqBody := map[string]interface{}{
		"model":    g.config.Model,
		"file":     encodedAudio,
		"language": g.config.Language,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/audio/transcriptions", bytes.NewReader(reqJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Groq: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Text, nil
}

// int16ToWAV converts int16 samples to WAV format bytes
func int16ToWAV(samples []int16, sampleRate int) ([]byte, error) {
	buf := new(bytes.Buffer)

	// WAV header (44 bytes)
	numSamples := len(samples)
	dataSize := numSamples * 2 // 16-bit = 2 bytes per sample
	fileSize := 36 + dataSize

	// RIFF header
	buf.WriteString("RIFF")
	writeInt32(buf, uint32(fileSize))
	buf.WriteString("WAVE")

	// fmt chunk
	buf.WriteString("fmt ")
	writeInt32(buf, 16) // fmt chunk size
	writeInt16(buf, 1)  // audio format (PCM)
	writeInt16(buf, 1)  // num channels (mono)
	writeInt32(buf, uint32(sampleRate))
	writeInt32(buf, uint32(sampleRate*2)) // byte rate
	writeInt16(buf, 2)                    // block align
	writeInt16(buf, 16)                   // bits per sample

	// data chunk
	buf.WriteString("data")
	writeInt32(buf, uint32(dataSize))

	// Write samples
	for _, sample := range samples {
		writeInt16(buf, uint16(sample))
	}

	return buf.Bytes(), nil
}

func writeInt32(buf *bytes.Buffer, val uint32) {
	buf.WriteByte(byte(val))
	buf.WriteByte(byte(val >> 8))
	buf.WriteByte(byte(val >> 16))
	buf.WriteByte(byte(val >> 24))
}

func writeInt16(buf *bytes.Buffer, val uint16) {
	buf.WriteByte(byte(val))
	buf.WriteByte(byte(val >> 8))
}
