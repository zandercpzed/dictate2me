# dictate2me v0.1.0 - First Public Release 🎉

We're excited to announce the first public release of **dictate2me**!

## 🎯 What is dictate2me?

dictate2me is a **100% offline voice transcription tool** with AI-powered text correction. Perfect for privacy-conscious users who want local AI without sending data to the cloud.

## ✨ Features

### Core Functionality

- 🎤 **Real-time audio capture** with PortAudio
- 📝 **Offline transcription** using Vosk (Portuguese, English, Spanish, and more)
- ✏️ **AI text correction** using local LLM (Ollama)
- 🔒 **100% offline** - your data never leaves your machine

### API

- 🌐 **REST API** for editor integrations
- 📡 **WebSocket streaming** for real-time results
- 🔐 **Token-based authentication** for security
- 🚀 **<10ms latency** for fast responses

### Obsidian Plugin

- 🔌 **Native Obsidian integration**
- ⌨️ **Hotkey support** (Cmd/Ctrl+Shift+D)
- 💭 **Live partial results** as you speak
- 🎨 **Visual feedback** with pulsing microphone icon
- ⚙️ **Configurable settings** for customization

## 📊 Project Stats

- **~3,000 lines** of production code
- **~9,000 lines** of comprehensive documentation
- **85%+ test coverage** ensuring reliability
- **<10ms API latency** for responsive performance

## 🚀 Getting Started

### Prerequisites

- **Go 1.23+** for building the tools
- **PortAudio** for audio capture
- **Vosk model** for transcription (auto-downloadable)
- **Ollama** (optional) for text correction

### Quick Start

1. **Clone the repository:**

   ```bash
   git clone https://github.com/zandercpzed/dictate2me.git
   cd dictate2me
   ```

2. **Download dependencies:**

   ```bash
   make deps
   ```

3. **Download Vosk model:**

   ```bash
   ./scripts/download-vosk-models.sh small
   ```

4. **Build:**

   ```bash
   make build
   ```

5. **Run the daemon:**
   ```bash
   ./bin/dictate2me-daemon
   ```

Check out our [README](https://github.com/zandercpzed/dictate2me#readme) for detailed installation instructions.

## 📖 Documentation

This release includes extensive documentation:

- **[API Documentation](docs/API.md)** - Complete API reference
- **[Testing Guide](docs/TESTING.md)** - How to run and write tests
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute
- **[Architecture](docs/ARCHITECTURE.md)** - System design overview
- **[ADRs](docs/adr/)** - Architecture decision records
- **[Plugin Guide](plugins/obsidian-dictate2me/README.md)** - Obsidian plugin usage

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Areas open for contribution:

- 🐧 Linux support
- 🪟 Windows support
- 🌍 Internationalization (i18n)
- 🧪 Increase test coverage to 100%
- 📖 Documentation translations
- 🎨 UI improvements
- ⚡ Performance optimizations

## 🔒 Security

This project takes security seriously:

- All processing happens locally
- No data leaves your machine
- Token-based API authentication
- Localhost-only API access by default

For security issues, please see [SECURITY.md](SECURITY.md).

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

This project was built with a focus on:

- ✨ Technical excellence
- 📖 Abundant documentation
- 🤝 Open collaboration
- 🔒 User privacy

## 🛠️ Tech Stack

- **Backend:** Go 1.23+
- **Frontend:** TypeScript (Obsidian plugin)
- **Speech Recognition:** Vosk
- **LLM:** Ollama
- **Audio:** PortAudio
- **Streaming:** WebSocket

---

**Full Changelog**: https://github.com/zandercpzed/dictate2me/commits/v0.1.0

## 📞 Get Help

- **Documentation:** [docs/](https://github.com/zandercpzed/dictate2me/tree/main/docs)
- **Issues:** [GitHub Issues](https://github.com/zandercpzed/dictate2me/issues)
- **Discussions:** [GitHub Discussions](https://github.com/zandercpzed/dictate2me/discussions)

---

Thank you for your interest in dictate2me! We're excited to see what you build with it. 🚀
