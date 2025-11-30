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
	"fmt"
	"os"
)

// Version information (set via ldflags during build)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {
	// TODO: Implement CLI with cobra
	// For now, just print a placeholder message

	if len(os.Args) > 1 && os.Args[1] == "version" {
		printVersion()
		return
	}

	fmt.Println("🎤 dictate2me - Offline Voice Transcription & Correction")
	fmt.Println()
	fmt.Println("Status: 🚧 In Development (Phase 0: Bootstrap)")
	fmt.Println()
	fmt.Println("This is a placeholder. The CLI will be implemented in Phase 4.")
	fmt.Println()
	fmt.Println("For more information, see: https://github.com/zandercpzed/dictate2me")
}

// printVersion prints version information
func printVersion() {
	fmt.Printf("dictate2me %s\n", version)
	fmt.Printf("  commit:   %s\n", commit)
	fmt.Printf("  built:    %s\n", date)
	fmt.Printf("  built by: %s\n", builtBy)
}
