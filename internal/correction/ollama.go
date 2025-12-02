package correction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config holds configuration for the correction engine
type Config struct {
	// OllamaURL is the URL of the Ollama server (default: http://localhost:11434)
	OllamaURL string

	// Model is the name of the Ollama model to use (default: gemma2:2b)
	Model string

	// Timeout for correction requests (default: 30s)
	Timeout time.Duration

	// Temperature for LLM generation (0.0 = deterministic, default: 0.1)
	Temperature float32

	// SystemPrompt is the system prompt for the LLM
	SystemPrompt string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		OllamaURL:   "http://localhost:11434",
		Model:       "gemma2:2b",
		Timeout:     30 * time.Second,
		Temperature: 0.1,
		SystemPrompt: `Você é um corretor de texto em português brasileiro.
Sua tarefa é corrigir o texto fornecido adicionando:
- Pontuação adequada (vírgulas, pontos, etc)
- Capitalização correta (início de frases, nomes próprios)
- Correções gramaticais mínimas necessárias
- Formatação de parágrafos quando apropriado

IMPORTANTE: Mantenha o conteúdo e significado original. Não reescreva ou mude o sentido, apenas corrija.
Retorne APENAS o texto corrigido, sem explicações ou comentários adicionais.`,
	}
}

// Engine handles text correction using Ollama
type Engine struct {
	config     Config
	httpClient *http.Client
}

// New creates a new correction engine
func New(cfg Config) (*Engine, error) {
	if cfg.OllamaURL == "" {
		cfg.OllamaURL = DefaultConfig().OllamaURL
	}
	if cfg.Model == "" {
		cfg.Model = DefaultConfig().Model
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultConfig().Timeout
	}
	if cfg.SystemPrompt == "" {
		cfg.SystemPrompt = DefaultConfig().SystemPrompt
	}

	return &Engine{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

// Correct corrects the given text using the LLM
func (e *Engine) Correct(ctx context.Context, text string) (string, error) {
	if text == "" {
		return "", nil
	}

	// Trim whitespace
	text = strings.TrimSpace(text)
	if text == "" {
		return "", nil
	}

	// Build request
	reqBody := map[string]interface{}{
		"model":  e.config.Model,
		"prompt": text,
		"system": e.config.SystemPrompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": e.config.Temperature,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		e.config.OllamaURL+"/api/generate",
		bytes.NewReader(reqJSON),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Done {
		return "", fmt.Errorf("ollama did not finish generation")
	}

	// Clean up response
	corrected := strings.TrimSpace(result.Response)

	return corrected, nil
}

// HealthCheck verifies that Ollama is running and the model is available
func (e *Engine) HealthCheck(ctx context.Context) error {
	// Check if Ollama is running
	req, err := http.NewRequestWithContext(ctx, "GET", e.config.OllamaURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama is not running or not accessible at %s: %w", e.config.OllamaURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama health check failed with status %d", resp.StatusCode)
	}

	// Check if model exists
	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return fmt.Errorf("failed to decode tags response: %w", err)
	}

	modelFound := false
	for _, model := range tags.Models {
		if strings.HasPrefix(model.Name, e.config.Model) {
			modelFound = true
			break
		}
	}

	if !modelFound {
		return fmt.Errorf("model %s not found in ollama. Run: ollama pull %s", e.config.Model, e.config.Model)
	}

	return nil
}

// Model returns the model name being used
func (e *Engine) Model() string {
	return e.config.Model
}

// Close releases resources (currently no-op, but kept for consistency)
func (e *Engine) Close() error {
	return nil
}
