package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zandercpzed/dictate2me/internal/transcription"
)

func setupTestServer(t *testing.T) (*Server, string) {
	// Create mock transcription engine
	transEngine, err := transcription.New(transcription.Config{
		ModelPath:  "../../models/vosk-model-small-pt-0.3",
		SampleRate: 16000,
		Language:   "pt",
	})
	if err != nil {
		t.Skip("Skipping test, Vosk model not available:", err)
	}

	server, err := New(Config{
		Host:                "127.0.0.1",
		Port:                0, // Use random available port
		Token:               "test-token-123456",
		TranscriptionEngine: transEngine,
	})
	require.NoError(t, err)
	require.NotNil(t, server)

	return server, "Bearer test-token-123456"
}

func TestHandleHealth(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp HealthResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, "healthy", resp.Status)
	assert.Equal(t, "ready", resp.Services["transcription"])
	assert.Greater(t, resp.Uptime, int64(0))
}

func TestMiddlewareAuth(t *testing.T) {
	server, validToken := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid token",
			authHeader:     validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "wrong token",
			authHeader:     "Bearer wrong-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			handler := server.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHandleCorrect(t *testing.T) {
	server, token := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	tests := []struct {
		name           string
		request        CorrectRequest
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "valid request but no correction engine",
			request:        CorrectRequest{Text: "test"},
			expectedStatus: http.StatusServiceUnavailable,
			checkResponse:  false,
		},
		{
			name:           "empty text",
			request:        CorrectRequest{Text: ""},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/v1/correct", bytes.NewReader(body))
			req.Header.Set("Authorization", token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleCorrect(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	handler := server.rateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Send 100 requests (should pass)
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 101st request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestCORSMiddleware(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	handler := server.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name         string
		origin       string
		method       string
		expectHeader bool
	}{
		{
			name:         "localhost origin",
			origin:       "http://localhost:3000",
			method:       "GET",
			expectHeader: true,
		},
		{
			name:         "non-localhost origin",
			origin:       "http://example.com",
			method:       "GET",
			expectHeader: false,
		},
		{
			name:         "preflight request",
			origin:       "http://localhost:3000",
			method:       "OPTIONS",
			expectHeader: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if tt.expectHeader {
				assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

func TestServerStartShutdown(t *testing.T) {
	server, _ := setupTestServer(t)
	defer server.config.TranscriptionEngine.Close()

	// Start server in goroutine
	go func() {
		server.Start()
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestBytesToInt16(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []int16
		wantErr bool
	}{
		{
			name:    "valid data",
			input:   []byte{0x00, 0x01, 0xFF, 0xFF},
			want:    []int16{256, -1},
			wantErr: false,
		},
		{
			name:    "odd length",
			input:   []byte{0x00},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			input:   []byte{},
			want:    []int16{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bytesToInt16(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
