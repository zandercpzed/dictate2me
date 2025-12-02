// Package api provides a local HTTP REST API for dictate2me.
//
// The API server runs on localhost and provides endpoints for:
//   - Text transcription from audio
//   - Text correction using LLM
//   - Real-time streaming via WebSocket
//
// # Authentication
//
// All endpoints (except /health) require Bearer token authentication.
// The token is automatically generated on first run and saved to ~/.dictate2me/api-token.
//
// # Example Usage
//
//	// Create server
//	server, err := api.New(api.Config{
//		TranscriptionEngine: transcriptionEngine,
//		CorrectionEngine:    correctionEngine,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Start server
//	go server.Start()
//
//	// Shutdown gracefully
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	server.Shutdown(ctx)
//
// # Endpoints
//
//	GET  /api/v1/health      - Health check (no auth)
//	POST /api/v1/transcribe  - Transcribe audio (requires auth)
//	POST /api/v1/correct     - Correct text (requires auth)
//	GET  /api/v1/stream      - WebSocket streaming (requires auth)
//
// # Security
//
// The server is designed to run locally only:
//   - Binds to 127.0.0.1 (localhost)
//   - Token-based authentication
//   - CORS restricted to localhost origins
//   - Rate limiting (100 req/min per IP)
//
// For more details, see docs/adr/0006-api-rest-local.md
package api
