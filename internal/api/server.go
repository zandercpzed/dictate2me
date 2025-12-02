package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zandercpzed/dictate2me/internal/correction"
	"github.com/zandercpzed/dictate2me/internal/transcription"
)

// Config holds configuration for the API server
type Config struct {
	Host  string
	Port  int
	Token string // API token for authentication

	// TranscriptionEngine is required
	TranscriptionEngine *transcription.Engine

	// CorrectionEngine is optional
	CorrectionEngine *correction.Engine
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Host:  "127.0.0.1",
		Port:  8765,
		Token: "", // Will be auto-generated
	}
}

// Server is the HTTP API server
type Server struct {
	config Config
	server *http.Server
	logger *slog.Logger

	// State
	mu       sync.RWMutex
	started  time.Time
	requests int64
}

// New creates a new API server
func New(cfg Config) (*Server, error) {
	if cfg.Host == "" {
		cfg.Host = DefaultConfig().Host
	}
	if cfg.Port == 0 {
		cfg.Port = DefaultConfig().Port
	}
	if cfg.TranscriptionEngine == nil {
		return nil, fmt.Errorf("transcription engine is required")
	}

	// Generate or load API token
	if cfg.Token == "" {
		token, err := loadOrGenerateToken()
		if err != nil {
			return nil, fmt.Errorf("failed to load/generate token: %w", err)
		}
		cfg.Token = token
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s := &Server{
		config:  cfg,
		logger:  logger,
		started: time.Now(),
	}

	// Setup routes
	mux := http.NewServeMux()

	// Health check (no auth required)
	mux.HandleFunc("GET /api/v1/health", s.handleHealth)

	// API endpoints (auth required)
	mux.HandleFunc("POST /api/v1/transcribe", s.authMiddleware(s.handleTranscribe))
	mux.HandleFunc("POST /api/v1/correct", s.authMiddleware(s.handleCorrect))
	mux.HandleFunc("GET /api/v1/stream", s.authMiddleware(s.handleStream))

	// Apply global middleware
	handler := s.loggingMiddleware(
		s.corsMiddleware(
			s.rateLimitMiddleware(mux),
		),
	)

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("starting API server",
		"addr", s.server.Addr,
		"token", s.config.Token[:8]+"...",
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down API server")
	return s.server.Shutdown(ctx)
}

// Addr returns the server address
func (s *Server) Addr() string {
	return s.server.Addr
}

// Token returns the API token
func (s *Server) Token() string {
	return s.config.Token
}

// loadOrGenerateToken loads existing token or generates a new one
func loadOrGenerateToken() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".dictate2me")
	tokenFile := filepath.Join(configDir, "api-token")

	// Try to load existing token
	if data, err := os.ReadFile(tokenFile); err == nil {
		token := string(data)
		if len(token) > 0 {
			return token, nil
		}
	}

	// Generate new token
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	// Save token
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	if err := os.WriteFile(tokenFile, []byte(token), 0600); err != nil {
		return "", err
	}

	return token, nil
}

// generateToken generates a random API token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
