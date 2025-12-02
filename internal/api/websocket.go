package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow localhost origins only
		origin := r.Header.Get("Origin")
		return origin == "" || origin == "http://localhost:"+r.Host
	},
}

// StreamMessage represents a WebSocket message
type StreamMessage struct {
	Type string          `json:"type"` // "start", "audio", "stop", "partial", "final", "error"
	Data json.RawMessage `json:"data,omitempty"`
}

// StreamStartConfig represents stream start configuration
type StreamStartConfig struct {
	Language         string `json:"language"`
	EnableCorrection bool   `json:"enableCorrection"`
}

// StreamPartialData represents partial result data
type StreamPartialData struct {
	Text string `json:"text"`
}

// StreamFinalData represents final result data
type StreamFinalData struct {
	Transcript string  `json:"transcript"`
	Corrected  string  `json:"corrected,omitempty"`
	Confidence float32 `json:"confidence"`
}

// StreamErrorData represents error data
type StreamErrorData struct {
	Message string `json:"message"`
}

// handleStream handles WebSocket streaming
func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("websocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	s.logger.Info("websocket connection established", "remote", r.RemoteAddr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream state
	var config StreamStartConfig
	started := false

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		var msg StreamMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Error("websocket read error", "error", err)
			}
			break
		}

		// Reset read deadline on each message
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		switch msg.Type {
		case "start":
			if err := json.Unmarshal(msg.Data, &config); err != nil {
				s.sendError(conn, fmt.Sprintf("invalid start config: %v", err))
				continue
			}

			started = true
			s.logger.Info("stream started", "language", config.Language, "correction", config.EnableCorrection)

			// Reset transcription engine
			if err := s.config.TranscriptionEngine.Reset(); err != nil {
				s.sendError(conn, fmt.Sprintf("failed to reset engine: %v", err))
			}

		case "audio":
			if !started {
				s.sendError(conn, "stream not started")
				continue
			}

			// Decode audio data
			var audioData struct {
				Data string `json:"data"` // base64 encoded
			}
			if err := json.Unmarshal(msg.Data, &audioData); err != nil {
				s.sendError(conn, "invalid audio data")
				continue
			}

			decoded, err := base64.StdEncoding.DecodeString(audioData.Data)
			if err != nil {
				s.sendError(conn, "invalid base64 audio")
				continue
			}

			samples, err := bytesToInt16(decoded)
			if err != nil {
				s.sendError(conn, err.Error())
				continue
			}

			// Transcribe
			segments, err := s.config.TranscriptionEngine.TranscribeStream(samples)
			if err != nil {
				s.sendError(conn, fmt.Sprintf("transcription error: %v", err))
				continue
			}

			// Send final results
			for _, seg := range segments {
				finalText := seg.Text

				// Apply correction if enabled
				if config.EnableCorrection && s.config.CorrectionEngine != nil && finalText != "" {
					corrCtx, corrCancel := context.WithTimeout(ctx, 30*time.Second)
					corrected, corrErr := s.config.CorrectionEngine.Correct(corrCtx, finalText)
					corrCancel()

					if corrErr == nil {
						finalText = corrected
					}
				}

				s.sendFinal(conn, seg.Text, finalText, seg.Conf)
			}

			// Send partial result
			partial, err := s.config.TranscriptionEngine.PartialResult()
			if err == nil && partial != "" {
				s.sendPartial(conn, partial)
			}

		case "stop":
			s.logger.Info("stream stopped")

			// Get final result
			if started {
				finalSeg, err := s.config.TranscriptionEngine.FinalResult()
				if err == nil && finalSeg.Text != "" {
					finalText := finalSeg.Text

					// Apply correction if enabled
					if config.EnableCorrection && s.config.CorrectionEngine != nil {
						corrCtx, corrCancel := context.WithTimeout(ctx, 30*time.Second)
						corrected, corrErr := s.config.CorrectionEngine.Correct(corrCtx, finalText)
						corrCancel()

						if corrErr == nil {
							finalText = corrected
						}
					}

					s.sendFinal(conn, finalSeg.Text, finalText, finalSeg.Conf)
				}
			}

			// Close gracefully
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			return

		default:
			s.sendError(conn, "unknown message type: "+msg.Type)
		}
	}
}

// sendPartial sends a partial result message
func (s *Server) sendPartial(conn *websocket.Conn, text string) {
	data, _ := json.Marshal(StreamPartialData{Text: text})
	msg := StreamMessage{
		Type: "partial",
		Data: data,
	}
	if err := conn.WriteJSON(msg); err != nil {
		s.logger.Error("failed to send partial", "error", err)
	}
}

// sendFinal sends a final result message
func (s *Server) sendFinal(conn *websocket.Conn, transcript, corrected string, confidence float32) {
	finalData := StreamFinalData{
		Transcript: transcript,
		Corrected:  corrected,
		Confidence: confidence,
	}

	data, _ := json.Marshal(finalData)
	msg := StreamMessage{
		Type: "final",
		Data: data,
	}
	if err := conn.WriteJSON(msg); err != nil {
		s.logger.Error("failed to send final", "error", err)
	}
}

// sendError sends an error message
func (s *Server) sendError(conn *websocket.Conn, message string) {
	data, _ := json.Marshal(StreamErrorData{Message: message})
	msg := StreamMessage{
		Type: "error",
		Data: data,
	}
	if err := conn.WriteJSON(msg); err != nil {
		s.logger.Error("failed to send error", "error", err)
	}
}
