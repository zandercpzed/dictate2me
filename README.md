<div align="center">
  
  # dictate2me
  
  **Transcrição de voz e correção textual 100% offline**
  
  [![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
  [![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
  [![CI](https://github.com/zandercpzed/dictate2me/actions/workflows/ci.yaml/badge.svg)](https://github.com/zandercpzed/dictate2me/actions)
  
  [Instalação](#-instalação) •
  [Uso Rápido](#-uso-rápido) •
  [Documentação](#-documentação) •
  [Contribuindo](#-contribuindo)
</div>

---

## ✨ Funcionalidades

- 🎤 **Captura de Áudio** - Gravação em tempo real do microfone
- 📝 **Transcrição Offline** - Powered by Vosk, sem enviar dados para nuvem
- ✏️ **Correção Inteligente** - LLM local (Ollama) para gramática, sintaxe e pontuação
- 🔌 **Integração com Obsidian** - Plugin nativo para inserção direta (em breve)
- 🖥️ **Cross-Platform** - macOS, Windows e Linux (em breve)
- 🔒 **Privacidade Total** - Seus dados nunca saem do seu computador
- 🌐 **API REST Local** - Integração com editores via HTTP/WebSocket 🆕

## 🚀 Instalação

### Pré-requisitos

- macOS 14+ (Sonoma) ou macOS 15 (Sequoia/Tahoe)
- 8GB RAM mínimo (16GB recomendado)
- 5GB de espaço em disco (para modelos de IA)

### Via Homebrew (Em breve)

```bash
brew tap zandercpzed/dictate2me
brew install dictate2me
```

### Download Direto

Baixe o binário mais recente em [Releases](https://github.com/zandercpzed/dictate2me/releases).

### Compilar do Código-Fonte

```bash
git clone https://github.com/zandercpzed/dictate2me.git
cd dictate2me
./scripts/setup-dev.sh
make build
```

## 📖 Uso Rápido

### 1. Baixar Modelos de IA

```bash
# Baixar modelo de transcrição Vosk
./scripts/download-vosk-models.sh small

# Instalar e configurar Ollama para correção de texto
./scripts/setup-ollama.sh
```

### 2. Iniciar Gravação

```bash
# Com correção de texto (requer Ollama)
dictate2me start

# Sem correção de texto
dictate2me start --no-correction
```

### 3. Transcrever Arquivo (em breve)

```bash
dictate2me transcribe audio.wav --output texto.txt
```

### 4. Usar API REST (para integrações)

```bash
# Iniciar daemon em background
dictate2me-daemon &

# A API estará disponível em http://localhost:8765
# Token salvo em ~/.dictate2me/api-token

# Testar health check
curl http://localhost:8765/api/v1/health

# Corrigir texto via API
curl -X POST http://localhost:8765/api/v1/correct \
  -H "Authorization: Bearer $(cat ~/.dictate2me/api-token)" \
  -H "Content-Type: application/json" \
  -d '{"text": "olá mundo"}'
```

Veja a [documentação completa da API](docs/API.md) para mais detalhes.

## ⚠️ Modo Degradado (sem transcrição)

Se o daemon for iniciado em um sistema sem a dependência nativa `libvosk` ou sem o modelo Vosk instalado, ele entrará em modo "degradado" — o servidor HTTP/WS continuará funcionando, mas a transcrição ficará desativada (retornos vazios). Isso é intencional para evitar que a aplicação falhe completamente em ambientes onde os binários nativos não estão disponíveis.

Como resolver (macOS):

- Instalar bibliotecas nativas (exemplo com Homebrew):

```bash
# Instalar dependências necessárias
brew install portaudio
# Baixar ou compilar libvosk (conforme instruções em docs/VOSK_INSTALLATION.md)
```

- Baixar modelo Vosk (recomendado: `vosk-model-small-pt-0.3`) e colocá-lo na pasta `models/` do projeto, ou passar `--model /caminho/para/modelo` ao iniciar o daemon:

```bash
./scripts/download-vosk-models.sh small
dictate2me-daemon --model models/vosk-model-small-pt-0.3
```

O arquivo `vosk_api.h` e a biblioteca compartilhada (`libvosk.dylib` ou `libvosk.so`) precisam estar disponíveis para que a integração CGO funcione. Se houver problemas de `#include` ou `dylib not found`, verifique os paths de `CGO_CFLAGS` / `CGO_LDFLAGS` e `DYLD_LIBRARY_PATH`.

Como instalar o Ollama (opcional, para correção de texto):

```bash
# Instale Ollama conforme https://ollama.com/docs
# Exemplo (macOS):
brew install ollama
ollama pull gemma2:2b
```

Se o Ollama não estiver disponível, o daemon ainda iniciará em modo degradado para correção (mas continuará realizando transcrição se o Vosk estiver ok). Para iniciar explicitamente sem correção use `--no-correction`.


## 📚 Documentação

| Documento                               | Descrição                  |
| --------------------------------------- | -------------------------- |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Visão geral da arquitetura |
| [API.md](docs/API.md)                   | Documentação da API REST   |
| [DEVELOPMENT.md](docs/DEVELOPMENT.md)   | Guia para desenvolvedores  |
| [ADRs](docs/adr/)                       | Decisões arquiteturais     |

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, leia nosso [Guia de Contribuição](CONTRIBUTING.md) antes de submeter PRs.

1. Fork o repositório
2. Crie sua branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'feat: add amazing feature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📊 Status do Projeto

🚀 **Em Desenvolvimento Ativo** - Fase 4: API REST Completa ✅

Próxima fase: Plugin Obsidian

Veja o [STATUS.md](STATUS.md) para detalhes do progresso.

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- [Vosk](https://alphacephei.com/vosk/) - Motor de transcrição offline
- [Ollama](https://ollama.com/) - Gerenciador de LLMs locais
- [Obsidian](https://obsidian.md/) - Editor de notas
- [PortAudio](http://www.portaudio.com/) - Cross-platform audio I/O

---

<div align="center">
  Feito com ❤️ pela comunidade open-source
</div>
