// Package main provides the main entry point for the dictate2me CLI application.
//
// dictate2me is an offline voice-to-text transcription and correction tool that:
//   - Captures audio from the microphone
//   - Transcribes speech to text using Whisper (100% offline)
//   - Corrects grammar and punctuation using a local LLM
//   - Integrates with text editors like Obsidian
//
// Usage:
//
//	dictate2me [command] [flags]
//
// Available Commands:
//
//	start       Start capturing and transcribing audio
//	stop        Stop the current recording session
//	transcribe  Transcribe an audio file
//	models      Manage AI models
//	version     Show version information
//	help        Help about any command
//
// For more information, see https://github.com/zandercpzed/dictate2me
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zandercpzed/dictate2me/internal/audio"
	"github.com/zandercpzed/dictate2me/internal/correction"
	"github.com/zandercpzed/dictate2me/internal/transcription"
)

// Version information (set via ldflags during build)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		runStart(os.Args[2:])
	case "version":
		printVersion()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runStart(args []string) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	modelPath := fs.String("model", "models/vosk-model-small-pt-0.3", "Path to Vosk model")
	noCorrection := fs.Bool("no-correction", false, "Disable text correction with LLM")
	ollamaModel := fs.String("ollama-model", "gemma2:2b", "Ollama model for correction")
	fs.Parse(args)

	fmt.Printf("Loading transcription model from: %s\n", *modelPath)

	// Init Audio
	capture, err := audio.New(audio.DefaultConfig())
	if err != nil {
		fmt.Printf("Error initializing audio: %v\n", err)
		os.Exit(1)
	}
	defer capture.Close()

	// Init Transcription
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

	// Init Correction (optional)
	var corrEngine *correction.Engine
	if !*noCorrection {
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
			fmt.Printf("✓ Correction engine ready (model: %s)\n", *ollamaModel)
		}
	}

	// Start Capture
	if err := capture.Start(); err != nil {
		fmt.Printf("Error starting capture: %v\n", err)
		os.Exit(1)
	}
	defer capture.Stop()

	fmt.Println("🎤 Listening... (Press Ctrl+C to stop)")
	if corrEngine != nil {
		fmt.Println("✏️  Text correction enabled")
	}
	fmt.Println()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Processing loop
	go func() {
		for {
			select {
			case samples, ok := <-capture.Stream():
				if !ok {
					return
				}
				segments, err := transEngine.TranscribeStream(samples)
				if err != nil {
					fmt.Printf("Error transcribing: %v\n", err)
					continue
				}
				for _, seg := range segments {
					transcribed := seg.Text
					finalText := transcribed

					// Apply correction if available
					if corrEngine != nil && transcribed != "" {
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						corrected, corrErr := corrEngine.Correct(ctx, transcribed)
						cancel()

						if corrErr != nil {
							fmt.Printf("⚠️  Correction error: %v\n", corrErr)
						} else {
							finalText = corrected
						}
					}

					// Clear line and print final result
					fmt.Printf("\r\033[K")
					if corrEngine != nil && finalText != transcribed {
						fmt.Printf("📝 %s\n", transcribed)
						fmt.Printf("✏️  %s\n", finalText)
					} else {
						fmt.Printf("> %s\n", finalText)
					}
				}

				// Show partial result
				partial, _ := transEngine.PartialResult()
				if partial != "" {
					fmt.Printf("\r💭 %s", partial)
				}
			case err := <-capture.Error():
				fmt.Printf("Audio error: %v\n", err)
			}
		}
	}()

	<-sigCh
	fmt.Println("\n\nStopping...")
}

func printUsage() {
	fmt.Println("Usage: dictate2me [command] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  start    Start capturing and transcribing audio")
	fmt.Println("  version  Show version information")
	fmt.Println("  help     Show this help message")
	fmt.Println()
	fmt.Println("Flags for start:")
	fmt.Println("  --model           Path to Vosk model (default: models/vosk-model-small-pt-0.3)")
	fmt.Println("  --no-correction   Disable text correction with LLM")
	fmt.Println("  --ollama-model    Ollama model for correction (default: gemma2:2b)")
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("dictate2me %s\n", version)
	fmt.Printf("  commit:   %s\n", commit)
	fmt.Printf("  built:    %s\n", date)
	fmt.Printf("  built by: %s\n", builtBy)
}
