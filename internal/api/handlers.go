package api

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"
)

// TranscribeRequest represents a transcription request
type TranscribeRequest struct {
	Audio      string  `json:"audio"`      // Base64 encoded WAV
	SampleRate float64 `json:"sampleRate"` // Default: 16000
	Language   string  `json:"language"`   // Default: "pt"
}

// TranscribeResponse represents a transcription response
type TranscribeResponse struct {
	Text       string    `json:"text"`
	Confidence float32   `json:"confidence"`
	Segments   []Segment `json:"segments"`
}

// Segment represents a transcription segment
type Segment struct {
	Text       string  `json:"text"`
	Start      float64 `json:"start"` // seconds
	End        float64 `json:"end"`   // seconds
	Confidence float32 `json:"confidence"`
}

// CorrectRequest represents a correction request
type CorrectRequest struct {
	Text string `json:"text"`
}

// CorrectResponse represents a correction response
type CorrectResponse struct {
	Original  string `json:"original"`
	Corrected string `json:"corrected"`
	Model     string `json:"model"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
	Uptime   int64             `json:"uptime"` // seconds
}

// handleHealth handles GET /api/v1/health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	services := map[string]string{
		"transcription": "ready",
	}

	// Check correction service
	if s.config.CorrectionEngine != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := s.config.CorrectionEngine.HealthCheck(ctx); err != nil {
			services["correction"] = fmt.Sprintf("unavailable: %v", err)
		} else {
			services["correction"] = "ready"
		}
	} else {
		services["correction"] = "disabled"
	}

	uptime := time.Since(s.started).Seconds()

	s.respondSuccess(w, HealthResponse{
		Status:   "healthy",
		Services: services,
		Uptime:   int64(uptime),
	})
}

// handleTranscribe handles POST /api/v1/transcribe
func (s *Server) handleTranscribe(w http.ResponseWriter, r *http.Request) {
	var req TranscribeRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate
	if req.Audio == "" {
		s.respondError(w, http.StatusBadRequest, "audio field is required")
		return
	}

	// Decode base64 audio
	audioData, err := base64.StdEncoding.DecodeString(req.Audio)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid base64 audio data")
		return
	}

	// Convert to int16 samples
	samples, err := bytesToInt16(audioData)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(samples) == 0 {
		s.respondError(w, http.StatusBadRequest, "empty audio data")
		return
	}

	// Transcribe
	segments, err := s.config.TranscriptionEngine.TranscribeStream(samples)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("transcription failed: %v", err))
		return
	}

	// Get final result if no segments
	if len(segments) == 0 {
		finalSeg, err := s.config.TranscriptionEngine.FinalResult()
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get final result: %v", err))
			return
		}
		if finalSeg.Text != "" {
			segments = append(segments, finalSeg)
		}
	}

	// Build response
	var text string
	var avgConf float32
	respSegments := make([]Segment, 0, len(segments))

	for _, seg := range segments {
		text += seg.Text + " "
		avgConf += seg.Conf
		respSegments = append(respSegments, Segment{
			Text:       seg.Text,
			Start:      seg.Start.Seconds(),
			End:        seg.End.Seconds(),
			Confidence: seg.Conf,
		})
	}

	if len(segments) > 0 {
		avgConf /= float32(len(segments))
	}

	s.respondSuccess(w, TranscribeResponse{
		Text:       text,
		Confidence: avgConf,
		Segments:   respSegments,
	})
}

// handleCorrect handles POST /api/v1/correct
func (s *Server) handleCorrect(w http.ResponseWriter, r *http.Request) {
	if s.config.CorrectionEngine == nil {
		s.respondError(w, http.StatusServiceUnavailable, "correction service is not available")
		return
	}

	var req CorrectRequest
	if err := s.decodeJSON(r, &req); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Text == "" {
		s.respondError(w, http.StatusBadRequest, "text field is required")
		return
	}

	// Correct with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	corrected, err := s.config.CorrectionEngine.Correct(ctx, req.Text)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("correction failed: %v", err))
		return
	}

	s.respondSuccess(w, CorrectResponse{
		Original:  req.Text,
		Corrected: corrected,
		Model:     s.config.CorrectionEngine.Model(),
	})
}

// bytesToInt16 converts byte array to int16 samples
func bytesToInt16(data []byte) ([]int16, error) {
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("audio data length must be even")
	}

	samples := make([]int16, len(data)/2)
	for i := 0; i < len(samples); i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(data[i*2:]))
	}

	return samples, nil
}
