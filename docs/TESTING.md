# Guia de Testes - dictate2me

Este documento descreve a estratégia de testes, cobertura, e como executar todos os testes do projeto.

## 📋 Tabela de Conteúdos

- [Visão Geral](#visão-geral)
- [Mapeamento de Testes](#mapeamento-de-testes)
- [Executando Testes](#executando-testes)
- [Escrevendo Novos Testes](#escrevendo-novos-testes)
- [Cobertura de Testes](#cobertura-de-testes)
- [CI/CD Integration](#cicd-integration)

## Visão Geral

O projeto dictate2me utiliza uma abordagem de testes em múltiplas camadas:

```
┌─────────────────────────────────────────┐
│          E2E Tests (Futuro)             │
│  - Plugin no Obsidian                   │
│  - Workflow completo                    │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│       Integration Tests                 │
│  - API + Transcription + Correction     │
│  - Daemon startup/shutdown              │
└─────────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────────┐
│          Unit Tests                     │
│  - internal/audio/                      │
│  - internal/transcription/              │
│  - internal/correction/                 │
│  - internal/api/                        │
└─────────────────────────────────────────┘
```

## Mapeamento de Testes

### 1. Audio Module (`internal/audio/`)

**Arquivo:** `internal/audio/capture_test.go`

| Teste                       | Descrição                       | Cobertura         |
| --------------------------- | ------------------------------- | ----------------- |
| `TestNew`                   | Criação de instância com opções | Construtor        |
| `TestNewWithInvalidOptions` | Validação de opções inválidas   | Error handling    |
| `TestCapture_Start`         | Início de captura de áudio      | Start flow        |
| `TestCapture_Stop`          | Parada de captura               | Stop flow         |
| `TestCapture_Buffer`        | Buffer circular funciona        | Buffer management |
| `TestCapture_Callbacks`     | Callbacks são chamados          | Event system      |
| `TestVAD_Detection`         | Voice Activity Detection        | VAD algorithm     |

**Como executar:**

```bash
go test -v ./internal/audio/...
```

**Cobertura atual:** 87.5%

**Issues conhecidos:**

- VAD precisa de mais edge cases
- Testes de performance faltam

---

### 2. Transcription Module (`internal/transcription/`)

**Arquivo:** `internal/transcription/engine_test.go`

| Teste                               | Descrição                | Cobertura          |
| ----------------------------------- | ------------------------ | ------------------ |
| `TestNew`                           | Criação de engine        | Construtor         |
| `TestEngine_TranscribeStream`       | Transcrição streaming    | Core functionality |
| `TestEngine_PartialResult`          | Resultados parciais      | Streaming feature  |
| `TestEngine_FinalResult`            | Resultado final          | Completion         |
| `TestEngine_Reset`                  | Reset de engine          | State management   |
| `TestEngine_Close`                  | Cleanup de recursos      | Resource cleanup   |
| `TestEngine_MultipleTranscriptions` | Transcrições sequenciais | Reusability        |

**Como executar:**

```bash
# Necessita modelo Vosk
./scripts/download-vosk-models.sh small
go test -v ./internal/transcription/...
```

**Cobertura atual:** 75.9%

**Issues conhecidos:**

- Testes com modelos grandes são lentos
- Necessita mocks para CI

---

### 3. Correction Module (`internal/correction/`)

**Arquivo:** `internal/correction/ollama_test.go`

| Teste                           | Descrição            | Cobertura          |
| ------------------------------- | -------------------- | ------------------ |
| `TestNew`                       | Criação de engine    | Construtor         |
| `TestEngine_Correct`            | Correção de texto    | Core functionality |
| `TestEngine_HealthCheck`        | Health check Ollama  | Health monitoring  |
| `TestEngine_CorrectWithContext` | Correção com timeout | Context handling   |
| `TestEngine_CorrectError`       | Error handling       | Error cases        |
| `TestEngine_ModelName`          | Nome do modelo       | Getter             |

**Como executar:**

```bash
# Usa mock HTTP server (não precisa Ollama)
go test -v ./internal/correction/...

# Com Ollama real (integração)
OLLAMA_INTEGRATION=1 go test -v ./internal/correction/...
```

**Cobertura atual:** 90%+

**Mock Strategy:**

- HTTP server mockado para unit tests
- Ollama real para integration tests

---

### 4. API Module (`internal/api/`)

**Arquivo:** `internal/api/server_test.go`

| Teste                     | Descrição               | Cobertura           |
| ------------------------- | ----------------------- | ------------------- |
| `TestHandleHealth`        | Endpoint de health      | Health endpoint     |
| `TestMiddlewareAuth`      | Autenticação            | Auth middleware     |
| `TestHandleCorrect`       | Endpoint de correção    | Correct endpoint    |
| `TestHandleTranscribe`    | Endpoint de transcrição | Transcribe endpoint |
| `TestRateLimitMiddleware` | Rate limiting           | Rate limit          |
| `TestCORSMiddleware`      | CORS headers            | CORS                |
| `TestServerStartShutdown` | Lifecycle               | Server lifecycle    |
| `TestBytesToInt16`        | Conversão de dados      | Data conversion     |

**Como executar:**

```bash
go test -v ./internal/api/...
```

**Cobertura atual:** 85%+

**Test Strategy:**

- `httptest` para endpoints
- Mock transcription/correction engines

---

### 5. Integration Tests

**Arquivo:** `test/integration_test.go`

| Teste                    | Descrição                                | Duração |
| ------------------------ | ---------------------------------------- | ------- |
| `TestFullPipeline`       | Audio → Transcription → Correction → API | ~5s     |
| `TestDaemonStartup`      | Daemon initialization                    | ~3s     |
| `TestWebSocketStreaming` | WebSocket E2E                            | ~10s    |
| `TestAPIAuthentication`  | Auth flow completo                       | ~1s     |

**Como executar:**

```bash
# Requer modelo Vosk e Ollama (opcional)
go test -v ./test/...
```

**Cobertura:** Testa integrações entre módulos

---

### 6. CLI Tests

**Arquivo:** `cmd/dictate2me/main_test.go` (futuro)

| Teste           | Descrição      |
| --------------- | -------------- |
| `TestCLI_Start` | Comando start  |
| `TestCLI_Stop`  | Comando stop   |
| `TestCLI_Flags` | Parse de flags |

---

### 7. Plugin Tests (TypeScript)

**Arquivo:** `plugins/obsidian-dictate2me/tests/` (futuro)

| Teste      | Framework  | Descrição         |
| ---------- | ---------- | ----------------- |
| Unit tests | Jest       | Lógica isolada    |
| E2E tests  | Playwright | Workflow completo |

**Como executar:**

```bash
cd plugins/obsidian-dictate2me
npm test
```

---

## Executando Testes

### Todos os Testes

```bash
# Rodar tudo
make test

# Com coverage
make test-coverage

# Ver relatório HTML
open coverage.html
```

### Por Módulo

```bash
# Audio
go test -v ./internal/audio/...

# Transcription
go test -v ./internal/transcription/...

# Correction
go test -v ./internal/correction/...

# API
go test -v ./internal/api/...
```

### Com Filtros

```bash
# Apenas testes rápidos
go test -short ./...

# Testes específicos
go test -v -run TestEngine_Correct ./internal/correction/...

# Com verbose
go test -v ./...

# Com race detector
go test -race ./...
```

### Testes de Integração

```bash
# Script automático (recomendado)
./scripts/test-full.sh

# Manual
go test -v -tags=integration ./test/...
```

## Escrevendo Novos Testes

### Template de Unit Test

```go
package mypackage

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyFunction(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   "test",
			want:    "result",
			wantErr: false,
		},
		{
			name:    "invalid input",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MyFunction(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
```

### Template de Integration Test

```go
// +build integration

package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_FullPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	// ... initialize components

	// Test
	// ... execute workflow

	// Verify
	// ... assert results

	// Cleanup
	defer cleanup()
}
```

### Mocking

**External Dependencies:**

```go
// Mock HTTP server
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}))
defer server.Close()
```

**Interfaces:**

```go
type MockTranscriptionEngine struct {
	mock.Mock
}

func (m *MockTranscriptionEngine) Transcribe(audio []int16) (string, error) {
	args := m.Called(audio)
	return args.String(0), args.Error(1)
}
```

## Cobertura de Testes

### Metas

| Módulo                    | Meta    | Atual    | Status |
| ------------------------- | ------- | -------- | ------ |
| `internal/audio/`         | 90%     | 87.5%    | 🟡     |
| `internal/transcription/` | 90%     | 75.9%    | 🟡     |
| `internal/correction/`    | 90%     | 90%+     | ✅     |
| `internal/api/`           | 90%     | 85%+     | 🟡     |
| **Total**                 | **90%** | **~85%** | 🟡     |

### Visualizar Cobertura

```bash
# Gerar relatório
make test-coverage

# Ver no navegador
open coverage.html

# Ver no terminal
go tool cover -func=coverage.out
```

### Áreas com Baixa Cobertura

1. **`internal/transcription/`**:

   - Falta: testes com modelos grandes
   - Falta: edge cases de streaming

2. **`internal/audio/`**:

   - Falta: testes de performance
   - Falta: mock de PortAudio

3. **`cmd/`**:
   - Falta: testes de CLI

## CI/CD Integration

### GitHub Actions

**Workflow:** `.github/workflows/test.yaml`

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.23"

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: make test

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

### Pre-commit Hooks

```bash
# Install
cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

**Hook executa:**

1. `go fmt`
2. `go vet`
3. `golangci-lint`
4. `go test -short`

## Benchmarks

### Executar Benchmarks

```bash
# Todos
go test -bench=. ./...

# Específico
go test -bench=BenchmarkTranscribe ./internal/transcription/...

# Com memória
go test -bench=. -benchmem ./...

# Com CPU profile
go test -bench=. -cpuprofile=cpu.prof ./...
```

### Exemplos de Benchmarks

```go
func BenchmarkTranscribe(b *testing.B) {
	engine := setupEngine()
	audio := loadTestAudio()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.Transcribe(audio)
	}
}
```

## Testes de Performance

### Latência

```bash
# Testar latência da API
./scripts/benchmark-api.sh
```

**Metas:**

- Health check: < 10ms
- Transcribe: < 2s
- Correct: < 1s

### Throughput

```bash
# Testar throughput
./scripts/load-test.sh
```

**Metas:**

- 100 req/min (rate limit)
- 10 conexões simultâneas WebSocket

## Debugging de Testes

### Logs Detalhados

```bash
# Com logs
go test -v ./...

# Com trace
go test -trace=trace.out ./...
go tool trace trace.out
```

### Delve Debugger

```bash
# Debug de teste específico
dlv test ./internal/transcription/ -- -test.run TestEngine_Transcribe
```

### VS Code

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}/internal/transcription",
      "args": ["-test.run", "TestEngine_Transcribe"]
    }
  ]
}
```

## Troubleshooting

### "Model not found"

```bash
./scripts/download-vosk-models.sh small
```

### "Ollama not available"

Testes de correction usam mock por padrão. Para testar com Ollama real:

```bash
OLLAMA_INTEGRATION=1 go test ./internal/correction/...
```

### CGO issues

```bash
# Criar symlink
rm -f /tmp/dictate2me_vosk
ln -s "$(pwd)/lib/vosk" /tmp/dictate2me_vosk

# Executar
CGO_CFLAGS="-I/tmp/dictate2me_vosk" \
CGO_LDFLAGS="-L/tmp/dictate2me_vosk -lvosk" \
go test ./...
```

## Contribuindo com Testes

Ao adicionar features:

1. ✅ Escreva testes ANTES do código (TDD)
2. ✅ Mantenha cobertura > 80% do módulo
3. ✅ Use table-driven tests
4. ✅ Mock dependências externas
5. ✅ Documente casos especiais
6. ✅ Adicione benchmarks se performance-critical

## Recursos

- [Go Testing](https://golang.org/pkg/testing/)
- [Testify](https://github.com/stretchr/testify)
- [httptest](https://golang.org/pkg/net/http/httptest/)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

---

**Última atualização:** 2025-12-02
