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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/zandercpzed/dictate2me/internal/audio"
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
	fs.Parse(args)

	fmt.Printf("Loading model from: %s\n", *modelPath)

	// Init Audio
	capture, err := audio.New(audio.DefaultConfig())
	if err != nil {
		fmt.Printf("Error initializing audio: %v\n", err)
		os.Exit(1)
	}
	defer capture.Close()

	// Init Transcription
	engine, err := transcription.New(transcription.Config{
		ModelPath:  *modelPath,
		SampleRate: 16000,
		Language:   "pt",
	})
	if err != nil {
		fmt.Printf("Error initializing transcription: %v\n", err)
		os.Exit(1)
	}
	defer engine.Close()

	// Start Capture
	if err := capture.Start(); err != nil {
		fmt.Printf("Error starting capture: %v\n", err)
		os.Exit(1)
	}
	defer capture.Stop()

	fmt.Println("🎤 Listening... (Press Ctrl+C to stop)")

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
				segments, err := engine.TranscribeStream(samples)
				if err != nil {
					fmt.Printf("Error transcribing: %v\n", err)
					continue
				}
				for _, seg := range segments {
					// Clear line and print final result
					fmt.Printf("\r\033[K> %s\n", seg.Text)
				}

				// Show partial result
				partial, _ := engine.PartialResult()
				if partial != "" {
					fmt.Printf("\r... %s", partial)
				}
			case err := <-capture.Error():
				fmt.Printf("Audio error: %v\n", err)
			}
		}
	}()

	<-sigCh
	fmt.Println("\nStopping...")
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
	fmt.Println("  --model  Path to Vosk model (default: models/vosk-model-small-pt-0.3)")
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("dictate2me %s\n", version)
	fmt.Printf("  commit:   %s\n", commit)
	fmt.Printf("  built:    %s\n", date)
	fmt.Printf("  built by: %s\n", builtBy)
}
