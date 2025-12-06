package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// authMiddleware validates the API token
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// 1. Check Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "" {
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 2. Check Query parameter (for WebSockets)
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		// LOGGING FOR DEBUGGING
		s.logger.Info("auth attempt",
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
			"has_auth_header", auth != "",
			"token_source", func() string {
				if auth != "" {
					return "header"
				}
				if r.URL.Query().Get("token") != "" {
					return "query"
				}
				return "none"
			}(),
			"received_token_len", len(token),
			"expected_token_len", len(s.config.Token),
			"match", token == s.config.Token,
		)

		if token == "" {
			s.respondError(w, http.StatusUnauthorized, "missing Authorization header or token query param")
			return
		}

		if token != s.config.Token {
			s.respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		next(w, r)
	}
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow localhost and Obsidian app origins
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost",
			"app://obsidian.md",
			"capacitor://localhost",
		}

		allowed := false
		for _, prefix := range allowedOrigins {
			if origin != "" && strings.HasPrefix(origin, prefix) {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment request counter
		s.mu.Lock()
		s.requests++
		s.mu.Unlock()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		s.logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start),
			"ip", r.RemoteAddr,
		)
	})
}

// rateLimitMiddleware implements simple rate limiting
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	type client struct {
		requests int
		window   time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Cleanup old entries periodically
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for ip, c := range clients {
				if now.Sub(c.window) > 1*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		c, exists := clients[ip]
		if !exists {
			c = &client{window: time.Now()}
			clients[ip] = c
		}

		// Reset window if expired
		if time.Since(c.window) > 1*time.Minute {
			c.requests = 0
			c.window = time.Now()
		}

		c.requests++
		requests := c.requests
		mu.Unlock()

		// Rate limit: 100 requests per minute
		if requests > 100 {
			s.respondError(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Hijack implements http.Hijacker to allow WebSockets to work
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not support Hijacker")
}

// Helper functions for responses

func (s *Server) respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to encode JSON response", "error", err)
	}
}

func (s *Server) respondError(w http.ResponseWriter, code int, message string) {
	s.respondJSON(w, code, map[string]string{
		"error": message,
	})
}

func (s *Server) respondSuccess(w http.ResponseWriter, data interface{}) {
	s.respondJSON(w, http.StatusOK, data)
}

// decodeJSON decodes JSON request body
func (s *Server) decodeJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("empty request body")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}
