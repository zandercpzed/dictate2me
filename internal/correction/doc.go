// Package correction provides intelligent text correction using local LLMs.
//
// This package handles:
//   - Grammar and syntax correction for Portuguese (PT-BR)
//   - Punctuation normalization
//   - LLM model management (Phi-3, Gemma)
//   - Prompt engineering for correction tasks
//   - Correction caching for common patterns
//
// The package uses llama.cpp via CGO bindings for local LLM inference.
//
// Example usage:
//
//	corrector, err := correction.New(
//	    correction.WithModel("models/phi-3-mini-4k-q4.gguf"),
//	    correction.WithPromptTemplate(correction.PTBRGrammarPrompt),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer corrector.Close()
//
//	corrected, err := corrector.Correct("eu fui na casa do joao ontem")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Corrected:", corrected)
//	// Output: Eu fui à casa do João ontem.
package correction
