# Migração para Vosk - Resumo e Próximos Passos

## ✅ Concluído

### 1. Documentação Atualizada

- ✅ ADR-0004 criado documentando a decisão de migrar para Vosk
- ✅ `prompt.md` atualizado com stack Vosk
- ✅ `VOSK_INSTALLATION.md` criado com instruções de instalação

### 2. Código Implementado

- ✅ `internal/transcription/vosk.go` - Engine completo com:
  - Streaming API
  - Partial results para feedback em tempo real
  - Word-level timestamps
  - Confidence scores
- ✅ `internal/transcription/vosk_test.go` - Testes completos
- ✅ `go.mod` atualizado para usar `github.com/alphacep/vosk-api/go`

### 3. Scripts e Ferramentas

- ✅ `scripts/download-vosk-models.sh` - Download automático de modelos
- ✅ Modelo pequeno PT baixado (50MB) em `models/vosk-model-small-pt-0.3/`

### 4. Limpeza

- ✅ Removido código Whisper.cpp (`whisper.go`, `binding/`, `whisper.cpp/`)
- ✅ Removidas dependências Whisper do `go.mod`

## ⚠️ Pendente

### 1. Instalação da Biblioteca Vosk C

**Problema**: Vosk Go bindings requerem `libvosk` (biblioteca C) instalada no sistema.

**Opções**:

#### Opção A: Usar binários pré-compilados (Mais Simples)

```bash
# Download e instalação manual
curl -L -o /tmp/vosk.zip https://github.com/alphacep/vosk-api/releases/download/v0.3.50/vosk-osx-0.3.50.zip
unzip /tmp/vosk.zip -d /tmp/vosk
sudo cp /tmp/vosk/libvosk.dylib /usr/local/lib/
sudo cp /tmp/vosk/vosk_api.h /usr/local/include/
```

**Nota**: Verificar se v0.3.50 tem binários macOS disponíveis.

#### Opção B: Compilar do código-fonte (Mais Complexo)

Requer:

- Kaldi (framework de ASR)
- OpenBLAS
- OpenFST
- Várias horas de compilação

**Não recomendado** para desenvolvimento rápido.

#### Opção C: Usar Docker (Alternativa)

```dockerfile
FROM golang:1.23-alpine
RUN apk add --no-cache vosk-api-dev
# ... resto do Dockerfile
```

### 2. Testes

Após instalar `libvosk`:

```bash
# Testar compilação
go build ./internal/transcription

# Executar testes
go test -v ./internal/transcription

# Testar com modelo real
go test -v -run TestTranscribeStream ./internal/transcription
```

### 3. Integração com o Resto do Sistema

Atualizar:

- `cmd/dictate2me/main.go` - Usar novo engine Vosk
- `internal/api/handlers.go` - Endpoints de transcrição
- Documentação de uso

## 📊 Comparação: Antes vs Depois

| Aspecto                | Whisper.cpp              | Vosk                    |
| ---------------------- | ------------------------ | ----------------------- |
| **Tamanho do modelo**  | 500MB+                   | 50MB ✅                 |
| **Latência**           | 1-2s (batch)             | <100ms (streaming) ✅   |
| **Complexidade build** | Alta (CGO + C++ + Metal) | Média (CGO + C)         |
| **Instalação**         | Compilação manual        | Binários disponíveis ✅ |
| **API**                | Batch processing         | Streaming nativo ✅     |
| **Uso de RAM**         | ~2GB                     | ~500MB ✅               |
| **Acurácia**           | Excelente                | Muito boa               |

## 🎯 Próximos Passos Recomendados

1. **Instalar libvosk** usando Opção A (binários pré-compilados)
2. **Testar compilação** do módulo de transcrição
3. **Executar testes** para validar funcionalidade
4. **Atualizar CLI** para usar novo engine
5. **Documentar** processo de instalação no README.md

## 🔗 Recursos

- [Vosk API](https://alphacephei.com/vosk/)
- [Vosk Releases](https://github.com/alphacep/vosk-api/releases)
- [Vosk Models](https://alphacephei.com/vosk/models)
- [Go Bindings](https://github.com/alphacep/vosk-api/tree/master/go)

## 💡 Alternativa: Usar Vosk Server

Se a instalação da biblioteca C for muito complexa, considerar usar **Vosk Server** via HTTP/WebSocket:

```bash
# Rodar servidor Vosk em Docker
docker run -d -p 2700:2700 alphacep/kaldi-pt:latest

# Cliente Go faz requisições HTTP
# Não requer libvosk instalada localmente
```

Vantagens:

- ✅ Sem dependências C no cliente
- ✅ Fácil deployment
- ✅ Escalável

Desvantagens:

- ❌ Requer Docker/servidor separado
- ❌ Latência de rede adicional
- ❌ Não é 100% offline (depende de localhost)
