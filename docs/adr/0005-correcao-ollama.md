# ADR-0005: Correção de Texto com Ollama

## Status

**ACEITO** - 2025-12-01

## Contexto

Após implementar a transcrição de voz com Vosk (ADR-0004), precisamos de um sistema de correção automática de texto. O texto transcrito geralmente contém:

1. **Erros de pontuação**: Falta de vírgulas, pontos, etc
2. **Erros de capitalização**: Início de frases, nomes próprios
3. **Erros gramaticais**: Concordância, conjugação
4. **Falta de formatação**: Parágrafos, quebras de linha

### Requisitos

- ✅ **100% offline**: Nenhum dado pode ser enviado para servidores externos
- ✅ **Baixo consumo de recursos**: Máximo 4GB RAM adicional
- ✅ **Latência aceitável**: Correção em até 2 segundos
- ✅ **Português BR**: Otimizado para português brasileiro
- ✅ **Fácil instalação**: Setup simples para usuários

### Opções Avaliadas

#### 1. llama.cpp com bindings Go

**Prós:**

- Controle total sobre o modelo
- Performance excelente
- Suporte a quantização (Q4, Q5, Q8)

**Contras:**

- Requer CGO complexo (similar ao Whisper.cpp que removemos)
- Build system complicado
- Necessário gerenciar lifecycle do modelo manualmente
- Difícil atualizar modelos

#### 2. Ollama

**Prós:**

- API REST simples e bem documentada
- Gerenciamento automático de modelos
- Instalação trivial (`brew install ollama`)
- Suporte nativo a streaming
- Comunidade ativa (similar ao Docker para LLMs)
- Cliente Go oficial (`github.com/ollama/ollama/api`)
- Modelos facilmente atualizáveis (`ollama pull`)

**Contras:**

- Dependência externa (daemon deve estar rodando)
- Menos controle sobre otimizações

#### 3. LocalAI

**Prós:**

- API compatível com OpenAI
- Suporta múltiplos backends

**Contras:**

- Mais pesado que Ollama
- Menos focado em desktop
- Setup mais complexo

## Decisão

**Escolhemos Ollama como engine de correção.**

### Justificativa

1. **Simplicidade**: Ollama abstrai toda complexidade de gerenciamento de modelos
2. **Developer Experience**: API REST + cliente Go = integração trivial
3. **Manutenção**: Sem CGO = builds mais rápidos e estáveis
4. **Usuários**: `brew install ollama` é muito mais simples que baixar modelos GGUF manualmente
5. **Comunidade**: Ollama tem momentum forte, modelos constantemente atualizados

### Modelo Selecionado

**Primário**: `gemma2:2b` (Google Gemma 2B)

- Tamanho: ~1.7GB (Q4_K_M)
- RAM: ~2.5GB
- Latência: ~500ms (M1/M2)
- Qualidade: Excelente para português após fine-tuning do prompt

**Alternativo**: `qwen2.5:3b` (Alibaba Qwen 2.5 3B)

- Tamanho: ~2.1GB
- RAM: ~3GB
- Latência: ~700ms
- Multilíngue superior

## Implementação

### 1. Arquitetura

```
dictate2me (Go)
    ↓ HTTP
Ollama Daemon (localhost:11434)
    ↓
Modelo LLM (gemma2:2b)
```

### 2. Fluxo

```go
texto_transcrito := "olá mundo como vai você"

corrigido, err := corrector.Correct(texto_transcrito)
// corrigido = "Olá, mundo! Como vai você?"
```

### 3. Prompt System

```
Você é um corretor de texto em português brasileiro.
Sua tarefa é corrigir o texto fornecido adicionando:
- Pontuação adequada
- Capitalização correta
- Correções gramaticais mínimas

IMPORTANTE: Mantenha o conteúdo original. Não reescreva, apenas corrija.

Texto: {input}
Correto:
```

### 4. Dependências

```go
require (
    github.com/ollama/ollama v0.1.17
)
```

## Consequências

### Positivas

✅ **Zero CGO**: Builds limpos e rápidos
✅ **Manutenção Simples**: Ollama atualiza modelos automaticamente
✅ **Flexibilidade**: Fácil trocar de modelo sem rebuild
✅ **Developer-Friendly**: API REST familiar
✅ **Testabilidade**: Fácil mockar HTTP requests

### Negativas

⚠️ **Dependência Externa**: Usuários precisam instalar Ollama
⚠️ **Daemon**: Ollama deve estar rodando (mas pode auto-iniciar)
⚠️ **Controle Limitado**: Menos otimizações low-level possíveis

### Mitigações

- **Auto-start**: dictate2me pode verificar/iniciar Ollama automaticamente
- **Health Check**: Validar que Ollama está rodando antes de usar
- **Fallback**: Modo "sem correção" se Ollama não disponível
- **Documentação**: Instruções claras de instalação do Ollama

## Instalação do Ollama

### macOS

```bash
brew install ollama
ollama pull gemma2:2b
ollama serve  # Inicia o daemon
```

### Linux

```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama pull gemma2:2b
ollama serve
```

## Referências

- [Ollama GitHub](https://github.com/ollama/ollama)
- [Ollama Models](https://ollama.com/library)
- [Gemma 2B](https://ollama.com/library/gemma2)
- [Go Client](https://github.com/ollama/ollama/tree/main/api)
