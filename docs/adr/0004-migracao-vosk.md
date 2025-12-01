# ADR-0004: Migração de Whisper.cpp para Vosk

## Status

**ACEITO** - 2025-12-01

## Contexto

Inicialmente escolhemos Whisper.cpp como engine de transcrição (ADR-0003), mas durante a implementação encontramos diversos desafios:

1. **Complexidade de Build**: Whisper.cpp requer CGO com configurações complexas de LDFLAGS e CFLAGS
2. **Tamanho dos Modelos**: Modelos Whisper são grandes (500MB+ para small)
3. **Latência**: Whisper processa em batch, não é otimizado para streaming em tempo real
4. **Dependências**: Requer compilação de C++, Metal (macOS), CUDA (opcional), aumentando complexidade
5. **Manutenção**: Bindings Go não oficiais, com problemas de compatibilidade

## Decisão

**Migrar para Vosk como engine de transcrição principal.**

### Vosk - Características

- **Offline-first**: Totalmente offline, sem dependências de rede
- **Leve**: Modelos pequenos (50MB para português)
- **Streaming**: API nativa para streaming em tempo real com zero-latency
- **Simples**: Bindings Go oficiais, sem CGO complexo
- **Multilíngue**: Suporte nativo a 20+ idiomas incluindo português
- **Vocabulário Reconfigurável**: Permite ajustar vocabulário em runtime
- **Speaker Identification**: Suporte a identificação de falantes

### Comparação

| Critério               | Whisper.cpp           | Vosk               | Vencedor |
| ---------------------- | --------------------- | ------------------ | -------- |
| Tamanho do modelo (PT) | 500MB+                | 50MB               | ✅ Vosk  |
| Latência               | Batch (~1-2s)         | Streaming (<100ms) | ✅ Vosk  |
| Complexidade de build  | Alta (CGO + C++)      | Baixa (Go puro)    | ✅ Vosk  |
| Acurácia               | Excelente             | Muito boa          | Whisper  |
| Streaming API          | Não nativo            | Nativo             | ✅ Vosk  |
| Manutenção             | Bindings não oficiais | Bindings oficiais  | ✅ Vosk  |
| Uso de RAM             | ~2GB                  | ~500MB             | ✅ Vosk  |

## Implementação

### 1. Dependências

```go
// go.mod
require (
    github.com/alphacep/vosk-api/go v0.3.45
)
```

### 2. Estrutura

```
internal/transcription/
├── vosk.go           # Engine principal
├── vosk_test.go      # Testes
├── models.go         # Gerenciamento de modelos
└── doc.go            # Documentação
```

### 3. API

```go
type Engine struct {
    model      *vosk.VoskModel
    recognizer *vosk.VoskRecognizer
    config     Config
}

type Config struct {
    ModelPath  string
    SampleRate float64  // 16000 Hz
    Language   string   // "pt"
}

func New(cfg Config) (*Engine, error)
func (e *Engine) TranscribeStream(samples []int16) ([]Segment, error)
func (e *Engine) TranscribeFile(path string) ([]Segment, error)
func (e *Engine) Close() error
```

### 4. Modelos Disponíveis

**Português:**

- `vosk-model-small-pt-0.3` (50MB) - Uso geral, rápido
- `vosk-model-pt-fb-v0.1.1-20220516_2113` (1.6GB) - Alta precisão

**Download:**

```bash
# Modelo pequeno (recomendado)
wget https://alphacephei.com/vosk/models/vosk-model-small-pt-0.3.zip

# Modelo grande (opcional)
wget https://alphacephei.com/vosk/models/vosk-model-pt-fb-v0.1.1-20220516_2113.zip
```

## Consequências

### Positivas

✅ **Redução de Complexidade**: Eliminação de CGO complexo e dependências C++
✅ **Menor Footprint**: Modelos 10x menores, menos RAM
✅ **Melhor UX**: Streaming em tempo real com feedback instantâneo
✅ **Manutenibilidade**: Bindings oficiais com melhor suporte
✅ **Portabilidade**: Build mais simples em diferentes plataformas
✅ **Performance**: Menor latência para casos de uso em tempo real

### Negativas

⚠️ **Acurácia**: Whisper tem acurácia ligeiramente superior em alguns casos
⚠️ **Reconhecimento**: Whisper é mais conhecido na comunidade
⚠️ **Recursos**: Whisper tem mais recursos avançados (tradução, etc)

### Mitigações

- Para casos que exigem máxima acurácia, podemos oferecer Whisper como opção futura
- Vosk permite ajuste de vocabulário para melhorar acurácia em domínios específicos
- Comunidade Vosk é ativa e modelos estão em constante melhoria

## Notas

- Vosk é usado em produção por empresas como Mozilla, Telegram
- Modelos são treinados com Kaldi, framework robusto de ASR
- Suporte a GPU é opcional, funciona bem em CPU
- API permite partial results para feedback em tempo real

## Referências

- [Vosk Website](https://alphacephei.com/vosk/)
- [Vosk GitHub](https://github.com/alphacep/vosk-api)
- [Vosk Models](https://alphacephei.com/vosk/models)
- [Vosk Go Bindings](https://github.com/alphacep/vosk-api/tree/master/go)
