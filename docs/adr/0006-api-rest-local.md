# ADR-0006: API REST Local para Integração de Editores

## Status

**ACEITO** - 2025-12-01

## Contexto

Com os módulos de captura de áudio, transcrição (Vosk) e correção (Ollama) implementados, precisamos de uma forma de integrar o dictate2me com editores de texto como Obsidian, VS Code, etc.

### Requisitos

- ✅ **Comunicação local**: API deve rodar em localhost
- ✅ **Segurança**: Autenticação via token para evitar acesso não autorizado
- ✅ **Streaming**: Suporte a resultados parciais em tempo real
- ✅ **Simplicidade**: API RESTful simples e intuitiva
- ✅ **Documentação**: OpenAPI/Swagger para facilitar integração

### Casos de Uso

1. **Plugin Obsidian** inicia gravação via API
2. **VS Code Extension** envia áudio gravado para transcrição
3. **Aplicação web local** consome stream de transcrição em tempo real

## Decisão

**Implementar API REST local usando net/http padrão do Go + WebSocket para streaming.**

### Justificativa

1. **net/http stdlib**: Sem dependências externas, estável, performático
2. **gorilla/websocket**: Padrão da indústria para WebSocket em Go
3. **Localhost only**: Sem exposição a rede externa, mais seguro
4. **Token-based auth**: Simples e efetivo para uso local

## Arquitetura

### Endpoints

#### POST /api/v1/transcribe

Transcreve áudio enviado no body

**Request:**

```json
{
  "audio": "base64_encoded_wav",
  "sampleRate": 16000,
  "language": "pt"
}
```

**Response:**

```json
{
  "text": "Texto transcrito",
  "confidence": 0.95,
  "segments": [
    {
      "text": "Texto transcrito",
      "start": 0.0,
      "end": 2.5,
      "confidence": 0.95
    }
  ]
}
```

#### POST /api/v1/correct

Corrige texto usando LLM

**Request:**

```json
{
  "text": "olá mundo como vai você"
}
```

**Response:**

```json
{
  "original": "olá mundo como vai você",
  "corrected": "Olá, mundo! Como vai você?",
  "model": "gemma2:2b"
}
```

#### GET /api/v1/health

Health check da API e serviços

**Response:**

```json
{
  "status": "healthy",
  "services": {
    "transcription": "ready",
    "correction": "ready"
  },
  "uptime": 3600
}
```

#### WS /api/v1/stream

WebSocket para streaming em tempo real

**Messages (Client → Server):**

```json
{
  "type": "start",
  "config": {
    "language": "pt",
    "enableCorrection": true
  }
}
```

```json
{
  "type": "audio",
  "data": "base64_audio_chunk"
}
```

```json
{
  "type": "stop"
}
```

**Messages (Server → Client):**

```json
{
  "type": "partial",
  "text": "resultado parcial..."
}
```

```json
{
  "type": "final",
  "transcript": "Texto transcrito completo.",
  "corrected": "Texto corrigido completo.",
  "confidence": 0.95
}
```

```json
{
  "type": "error",
  "message": "Erro na transcrição"
}
```

### Autenticação

**Token simples gerado na inicialização:**

```
Authorization: Bearer <token>
```

Token é:

1. Gerado aleatoriamente na primeira execução
2. Salvo em `~/.dictate2me/api-token`
3. Lido pelo plugin/cliente
4. Enviado em todas as requisições

### Segurança

- ✅ Bind apenas em `127.0.0.1` (localhost)
- ✅ CORS restrito
- ✅ Rate limiting (100 req/min por padrão)
- ✅ Token obrigatório em produção
- ✅ Logs de acesso

## Implementação

### Estrutura

```
internal/api/
├── server.go          # HTTP server principal
├── server_test.go     # Testes do server
├── handlers.go        # Handlers dos endpoints
├── handlers_test.go   # Testes dos handlers
├── middleware.go      # Auth, CORS, logging
├── websocket.go       # Handler WebSocket
├── websocket_test.go  # Testes WebSocket
└── doc.go             # Documentação do pacote
```

### Daemon

```
cmd/dictate2me-daemon/
└── main.go            # Daemon que roda a API
```

**Uso:**

```bash
# Iniciar daemon em background
dictate2me-daemon start

# Parar daemon
dictate2me-daemon stop

# Status
dictate2me-daemon status
```

### Configuração

```yaml
# ~/.dictate2me/config.yaml
api:
  host: 127.0.0.1
  port: 8765
  token: auto-generated-or-custom
  cors:
    allowedOrigins:
      - http://localhost:*
  rateLimit:
    requestsPerMinute: 100
```

## Consequências

### Positivas

✅ **Integração Simples**: Qualquer linguagem pode consumir REST
✅ **Streaming Real-time**: WebSocket para feedback instantâneo
✅ **Segurança Local**: Token + localhost = seguro o suficiente
✅ **Escalável**: Fácil adicionar novos endpoints
✅ **Testável**: HTTP handlers são fáceis de testar

### Negativas

⚠️ **Daemon Necessário**: Editor precisa que daemon esteja rodando
⚠️ **Complexidade**: Mais código para manter
⚠️ **Token Management**: Usuário precisa configurar

### Mitigações

- **Auto-start**: Plugin tenta iniciar daemon se não estiver rodando
- **Status Check**: Endpoint /health para verificar disponibilidade
- **Documentação**: Guia claro de integração para desenvolvedores

## Alternativas Consideradas

### 1. gRPC

**Prós:** Performático, tipado, streaming nativo
**Contras:** Mais complexo, menos acessível para plugins JS/TS

### 2. Unix Domain Socket

**Prós:** Mais seguro, sem porta de rede
**Contras:** Não funciona no Windows, mais complexo para WebSocket

### 3. Named Pipes

**Prós:** Cross-platform
**Contras:** Complexo, sem suporte a HTTP/WebSocket

## Referências

- [Go net/http](https://pkg.go.dev/net/http)
- [gorilla/websocket](https://github.com/gorilla/websocket)
- [OpenAPI 3.0](https://swagger.io/specification/)
