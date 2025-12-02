package correction

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "default config",
			config: Config{},
		},
		{
			name: "custom config",
			config: Config{
				OllamaURL:   "http://custom:11434",
				Model:       "custom-model",
				Timeout:     10 * time.Second,
				Temperature: 0.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := New(tt.config)
			assert.NoError(t, err)
			assert.NotNil(t, engine)

			if tt.config.OllamaURL != "" {
				assert.Equal(t, tt.config.OllamaURL, engine.config.OllamaURL)
			} else {
				assert.Equal(t, DefaultConfig().OllamaURL, engine.config.OllamaURL)
			}
		})
	}
}

func TestCorrect(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		ollamaResponse string
		statusCode     int
		wantErr        bool
		expectedOutput string
	}{
		{
			name:           "simple correction",
			input:          "olá mundo como vai você",
			ollamaResponse: `{"response": "Olá, mundo! Como vai você?", "done": true}`,
			statusCode:     http.StatusOK,
			wantErr:        false,
			expectedOutput: "Olá, mundo! Como vai você?",
		},
		{
			name:           "empty input",
			input:          "",
			expectedOutput: "",
			wantErr:        false,
		},
		{
			name:           "whitespace only",
			input:          "   ",
			expectedOutput: "",
			wantErr:        false,
		},
		{
			name:       "ollama error",
			input:      "test",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:           "incomplete response",
			input:          "test",
			ollamaResponse: `{"response": "Test.", "done": false}`,
			statusCode:     http.StatusOK,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip empty/whitespace tests with mock server
			if tt.input == "" || len(tt.input) == len(strings.TrimSpace(tt.input)) && tt.input != strings.TrimSpace(tt.input) {
				engine, err := New(Config{})
				require.NoError(t, err)

				result, err := engine.Correct(context.Background(), tt.input)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedOutput, result)
				}
				return
			}

			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/generate", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Verify request body
				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.input, reqBody["prompt"])
				assert.Equal(t, false, reqBody["stream"])

				w.WriteHeader(tt.statusCode)
				if tt.ollamaResponse != "" {
					w.Write([]byte(tt.ollamaResponse))
				}
			}))
			defer server.Close()

			engine, err := New(Config{
				OllamaURL: server.URL,
			})
			require.NoError(t, err)

			result, err := engine.Correct(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, result)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		model          string
		tagsResponse   string
		statusCode     int
		wantErr        bool
		expectedErrMsg string
	}{
		{
			name:         "model exists",
			model:        "gemma2:2b",
			tagsResponse: `{"models": [{"name": "gemma2:2b"}]}`,
			statusCode:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "model not found",
			model:          "missing-model",
			tagsResponse:   `{"models": [{"name": "other-model"}]}`,
			statusCode:     http.StatusOK,
			wantErr:        true,
			expectedErrMsg: "model missing-model not found",
		},
		{
			name:           "ollama not running",
			model:          "gemma2:2b",
			statusCode:     http.StatusServiceUnavailable,
			wantErr:        true,
			expectedErrMsg: "health check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/tags", r.URL.Path)
				w.WriteHeader(tt.statusCode)
				if tt.tagsResponse != "" {
					w.Write([]byte(tt.tagsResponse))
				}
			}))
			defer server.Close()

			engine, err := New(Config{
				OllamaURL: server.URL,
				Model:     tt.model,
			})
			require.NoError(t, err)

			err = engine.HealthCheck(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCorrectContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "Corrected.", "done": true}`))
	}))
	defer server.Close()

	engine, err := New(Config{
		OllamaURL: server.URL,
		Timeout:   50 * time.Millisecond,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = engine.Correct(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestClose(t *testing.T) {
	engine, err := New(Config{})
	require.NoError(t, err)

	err = engine.Close()
	assert.NoError(t, err)

	// Should be idempotent
	err = engine.Close()
	assert.NoError(t, err)
}
