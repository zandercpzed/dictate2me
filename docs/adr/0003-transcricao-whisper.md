# ADR-0003: Uso de Whisper.cpp para Transcrição Offline

## Status

✅ Aceito

## Contexto

O dictate2me tem como requisito fundamental a transcrição de voz **100% offline** e com **alta eficiência** de recursos (CPU/RAM). Precisamos de um motor de transcrição que:

1. Funcione sem internet.
2. Tenha suporte excelente a Português (PT-BR).
3. Rode em CPUs de consumo (Apple Silicon, Intel/AMD) sem exigir GPUs dedicadas de alta performance.
4. Tenha latência aceitável para ditado (quase tempo real).
5. Seja integrável com Go.

## Decisão

**Decidimos usar o [whisper.cpp](https://github.com/ggerganov/whisper.cpp) como motor de transcrição.**

Whisper.cpp é uma portabilidade em C/C++ do modelo Whisper da OpenAI, otimizada para inferência em CPU e Apple Silicon.

### Bindings Go

Usaremos os bindings oficiais ou comunitários que permitam integração via CGO, carregando o modelo na memória e processando chunks de áudio PCM.

## Alternativas Consideradas

### Alternativa 1: OpenAI Whisper (Python original)

**Descrição**: Implementação oficial em PyTorch.

**Prós**:

- Implementação de referência.
- Acesso imediato a novos modelos.

**Contras**:

- ❌ **Dependência de Python**: Viola nosso requisito de binário único nativo.
- ❌ **Pesado**: PyTorch é enorme e consome muita RAM.
- ❌ **Lento em CPU**: Otimizado para GPU NVIDIA (CUDA).

### Alternativa 2: Vosk

**Descrição**: Biblioteca de reconhecimento de fala offline baseada em Kaldi.

**Prós**:

- Leve e rápido.
- Maduro e estável.

**Contras**:

- ❌ **Precisão inferior**: Modelos menos precisos que Whisper, especialmente em pontuação e contexto.
- ❌ **Modelos antigos**: Baseado em tecnologia pré-Transformer.

### Alternativa 3: Coqui STT (DeepSpeech)

**Descrição**: Engine open-source baseada no DeepSpeech da Mozilla.

**Prós**:

- Boa performance.

**Contras**:

- ❌ **Descontinuado**: Projeto perdeu manutenção ativa.
- ❌ **Qualidade variável**: Inferior ao Whisper em testes gerais.

## Por que Whisper.cpp?

1. **Otimização Extrema**: Usa instruções SIMD (AVX, AVX2, NEON) e quantização (4-bit, 5-bit) para reduzir uso de memória e CPU drasticamente.
2. **Apple Silicon Native**: Usa framework Accelerate/Metal no macOS, garantindo performance incrível nos Macs (nosso público inicial).
3. **Sem Dependências Pesadas**: Não precisa de PyTorch, TensorFlow ou Python. É apenas C++.
4. **Qualidade SOTA**: Mantém a precisão do modelo Whisper original.

## Estratégia de Modelos

Para equilibrar performance e qualidade, suportaremos:

1. **Default**: `whisper-small` (Q5_K_M) - Bom equilíbrio, ~500MB RAM.
2. **Performance**: `whisper-base` ou `tiny` - Para máquinas mais fracas.
3. **Qualidade**: `whisper-medium` - Para quem tem >16GB RAM e quer precisão máxima.

## Consequências

### Positivas

- ✅ **Zero Runtime Externo**: Compila junto com o binário Go.
- ✅ **Baixo Consumo**: Roda confortavelmente com < 2GB RAM (modelo small).
- ✅ **Portabilidade**: Funciona em Mac, Linux, Windows, iOS, Android.

### Negativas

- ⚠️ **Build Complexo**: Requer compilador C++ compatível e flags específicas de otimização.
- ⚠️ **Download de Modelos**: Usuário precisa baixar arquivos de pesos (.bin/.gguf) separadamente (~500MB+).

## Referências

- [whisper.cpp Repository](https://github.com/ggerganov/whisper.cpp)
- [Whisper Paper (OpenAI)](https://cdn.openai.com/papers/whisper.pdf)
- [Benchmarks Whisper.cpp](https://github.com/ggerganov/whisper.cpp#benchmarks)

---

**Data da Decisão**: 2025-01-30
**Decisores**: @zandercpzed
