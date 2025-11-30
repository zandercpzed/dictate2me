// Package transcription provides speech-to-text transcription using Whisper.
//
// This package handles:
//   - Loading and managing Whisper models (GGML/GGUF format)
//   - Audio transcription with support for multiple languages
//   - Model quantization for efficiency (Q5_K_M, Q4_K_M)
//   - Batch processing for optimal performance
//
// The package uses whisper.cpp via CGO bindings for efficient inference.
//
// Example usage:
//
//	engine, err := transcription.New(
//	    transcription.WithModel("models/whisper-small.bin"),
//	    transcription.WithLanguage("pt"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer engine.Close()
//
//	text, err := engine.Transcribe(audioData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Transcribed:", text)
package transcription
