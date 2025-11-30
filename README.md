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
- 📝 **Transcrição Offline** - Powered by Whisper, sem enviar dados para nuvem
- ✏️ **Correção Inteligente** - LLM local para gramática, sintaxe e pontuação
- 🔌 **Integração com Obsidian** - Plugin nativo para inserção direta
- 🖥️ **Cross-Platform** - macOS, Windows e Linux (em breve)
- 🔒 **Privacidade Total** - Seus dados nunca saem do seu computador

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
dictate2me models download
```

### 2. Iniciar Gravação

```bash
dictate2me start
```

### 3. Transcrever Arquivo

```bash
dictate2me transcribe audio.wav --output texto.txt
```

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

🚧 **Em Desenvolvimento Ativo** - Fase 0: Bootstrap (Semana 1)

Veja o [CHANGELOG.md](CHANGELOG.md) para histórico de versões.

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- [Whisper.cpp](https://github.com/ggerganov/whisper.cpp) - Motor de transcrição
- [llama.cpp](https://github.com/ggerganov/llama.cpp) - Inferência de LLM
- [Obsidian](https://obsidian.md/) - Editor de notas
- [PortAudio](http://www.portaudio.com/) - Cross-platform audio I/O

---

<div align="center">
  Feito com ❤️ pela comunidade open-source
</div>
