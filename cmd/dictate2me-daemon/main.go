// Package main provides the daemon process for dictate2me.
//
// The daemon runs in the background and handles:
//   - Continuous audio capture
//   - Real-time transcription
//   - Text correction
//   - API server for editor integrations
//
// The daemon is started automatically by the CLI when needed.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("🎤 dictate2me-daemon")
	fmt.Println()
	fmt.Println("Status: 🚧 In Development (Phase 0: Bootstrap)")
	fmt.Println()
	fmt.Println("The daemon will be implemented in Phase 4.")
	fmt.Println()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Press Ctrl+C to exit...")

	<-sigChan
	fmt.Println("\nShutting down gracefully...")
}
