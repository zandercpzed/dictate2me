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
