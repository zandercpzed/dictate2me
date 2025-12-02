// Package main provides the daemon process for dictate2me.
//
// The daemon runs in the background and handles:
//   - HTTP REST API for editor integrations
//   - WebSocket streaming for real-time transcription
//   - Token-based authentication
//
// Usage:
//
//	dictate2me-daemon [options]
//
// Options:
//
//	--port <port>         API server port (default: 8765)
//	--host <host>         API server host (default: 127.0.0.1)
//	--model <path>        Vosk model path
//	--ollama-model <name> Ollama model for correction
//	--no-correction       Disable text correction
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zandercpzed/dictate2me/internal/api"
	"github.com/zandercpzed/dictate2me/internal/correction"
	"github.com/zandercpzed/dictate2me/internal/transcription"
)

func main() {
	// Parse flags
	port := flag.Int("port", 8765, "API server port")
	host := flag.String("host", "127.0.0.1", "API server host")
	modelPath := flag.String("model", "models/vosk-model-small-pt-0.3", "Vosk model path")
	ollamaModel := flag.String("ollama-model", "gemma2:2b", "Ollama model for correction")
	noCorrection := flag.Bool("no-correction", false, "Disable text correction")
	flag.Parse()

	fmt.Println("🎤 dictate2me daemon")
	fmt.Printf("Starting API server at %s:%d\n", *host, *port)
	fmt.Println()

	// Initialize transcription engine
	fmt.Printf("Loading transcription model from: %s\n", *modelPath)
	transEngine, err := transcription.New(transcription.Config{
		ModelPath:  *modelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	if err != nil {
		fmt.Printf("Error initializing transcription: %v\n", err)
		os.Exit(1)
	}
	defer transEngine.Close()

	// Initialize correction engine (optional)
	var corrEngine *correction.Engine
	if !*noCorrection {
		fmt.Printf("Initializing correction engine (model: %s)\n", *ollamaModel)
		corrEngine, err = correction.New(correction.Config{
			Model: *ollamaModel,
		})
		if err != nil {
			fmt.Printf("Error initializing correction: %v\n", err)
			os.Exit(1)
		}
		defer corrEngine.Close()

		// Health check
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := corrEngine.HealthCheck(ctx); err != nil {
			cancel()
			fmt.Printf("⚠️  Ollama health check failed: %v\n", err)
			fmt.Println("   Running without correction. Install Ollama or use --no-correction flag.")
			corrEngine = nil
		} else {
			cancel()
			fmt.Printf("✓ Correction engine ready\n")
		}
	}

	// Create API server
	server, err := api.New(api.Config{
		Host:                *host,
		Port:                *port,
		TranscriptionEngine: transEngine,
		CorrectionEngine:    corrEngine,
	})
	if err != nil {
		fmt.Printf("Error creating API server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("✓ Daemon ready\n")
	fmt.Printf("  API endpoint: http://%s:%d\n", *host, *port)
	fmt.Printf("  Token: %s\n", server.Token())
	fmt.Printf("  Token saved to: ~/.dictate2me/api-token\n")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop...")
	fmt.Println()

	// Start API server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			fmt.Printf("API server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println()
	fmt.Println("Shutting down gracefully...")

	// Shutdown API server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}

	fmt.Println("✓ Daemon stopped")
}
