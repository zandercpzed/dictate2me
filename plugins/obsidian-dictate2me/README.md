# Dictate2Me - Obsidian Plugin

🎤 **Real-time voice transcription plugin for Obsidian with offline AI.**

Dictate text directly into your notes using local AI models, with automatic grammar correction. 100% offline, privacy-first.

## Features

- 🎤 **Voice Transcription** - Real-time speech-to-text
- ✏️ **Auto-Correction** - Grammar, punctuation, and capitalization via local LLM
- 🔒 **100% Offline** - All processing happens locally on your machine
- 🚀 **Low Latency** - Streaming results as you speak
- 🌍 **Multi-language** - Support for Portuguese, English, Spanish, and more
- ⌨️ **Hotkey Support** - Quick access via keyboard shortcut (default: `Cmd/Ctrl+Shift+D`)

## Requirements

1. **dictate2me daemon** must be running

   - Download from: [https://github.com/zandercpzed/dictate2me](https://github.com/zandercpzed/dictate2me)
   - Or install via: `brew install dictate2me` (coming soon)

2. **System Requirements**:
   - macOS 14+ (Sonoma) recommended
   - 8GB RAM minimum (16GB recommended)
   - Microphone access

## Installation

### From Obsidian Community Plugins (Coming Soon)

1. Open Obsidian Settings
2. Go to Community Plugins → Browse
3. Search for "Dictate2Me"
4. Click Install
5. Enable the plugin

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/zandercpzed/dictate2me/releases)
2. Extract `main.js`, `manifest.json`, and `styles.css` to your vault:
   ```
   YourVault/.obsidian/plugins/dictate2me/
   ```
3. Reload Obsidian
4. Enable the plugin in Settings → Community Plugins

## Setup

### 1. Start the Daemon

Before using the plugin, start the dictate2me daemon:

```bash
# Install daemon (if not already installed)
brew install dictate2me

# Start daemon
dictate2me-daemon
```

The daemon will:

- Start API server on `http://localhost:8765`
- Generate API token in `~/.dictate2me/api-token`
- Load AI models (Vosk for transcription, Ollama for correction)

### 2. Configure Plugin

1. Open Obsidian Settings → Dictate2Me
2. **API Token**: Copy token from `~/.dictate2me/api-token`
   ```bash
   cat ~/.dictate2me/api-token
   ```
3. **API URL**: Default is `http://localhost:8765/api/v1` (should work out of the box)
4. **Language**: Set to your language (e.g., `pt`, `en`, `es`)
5. **Enable Correction**: Toggle on/off for grammar correction
6. Click **Test Connection** to verify setup

### 3. First Use

1. Open a note in Obsidian
2. Click the microphone icon in the left ribbon, or press `Cmd/Ctrl+Shift+D`
3. Start speaking
4. Click again (or press hotkey) to stop
5. Text will be inserted at cursor position

## Usage

### Starting Dictation

**Method 1: Ribbon Icon**

- Click the microphone icon in the left sidebar
- Icon will pulse red while recording

**Method 2: Command Palette**

- Press `Cmd/Ctrl+P`
- Type "dictation"
- Select "Start/Stop Dictation"

**Method 3: Hotkey**

- Press `Cmd/Ctrl+Shift+D` (configurable in Obsidian hotkeys)

### Stopping Dictation

Use any of the same methods above. The recording will stop and final text will be inserted.

### Status Bar

The status bar (bottom of Obsidian) shows:

- **Ready**: Plugin is ready to record
- **🎤 Recording...**: Currently recording
- **💭 [partial text]**: Live transcription (if enabled in settings)
- **✓ Inserted (95% confidence)**: Text inserted successfully

## Settings

### API Configuration

| Setting   | Description            | Default                        |
| --------- | ---------------------- | ------------------------------ |
| API URL   | Daemon API endpoint    | `http://localhost:8765/api/v1` |
| API Token | Auth token from daemon | (empty)                        |
| Language  | Transcription language | `pt`                           |

### Features

| Setting                | Description                     | Default |
| ---------------------- | ------------------------------- | ------- |
| Enable text correction | Use LLM for grammar/punctuation | ✅ On   |
| Show partial results   | Display live transcription      | ✅ On   |
| Show confidence score  | Display confidence %            | ✅ On   |
| Auto-check daemon      | Verify daemon before recording  | ✅ On   |

### Hotkey

You can customize the hotkey in Obsidian Settings → Hotkeys → search for "Dictate2Me"

## Troubleshooting

### "Daemon is not running" error

**Solution:**

```bash
# Start daemon in terminal
dictate2me-daemon

# Or check if already running
curl http://localhost:8765/api/v1/health
```

### "Connection failed" error

1. Verify daemon is running:

   ```bash
   ps aux | grep dictate2me-daemon
   ```

2. Check API token:

   ```bash
   cat ~/.dictate2me/api-token
   ```

3. Test connection manually:

   ```bash
   curl http://localhost:8765/api/v1/health
   ```

4. Check plugin settings have correct URL and token

### No text inserted

1. Check microphone permissions in System Settings
2. Verify you're speaking clearly
3. Check status bar for errors
4. Look at Obsidian developer console (View → Toggle Developer Tools)

### Low quality transcription

1. Speak closer to microphone
2. Reduce background noise
3. Download larger Vosk model:
   ```bash
   ./scripts/download-vosk-models.sh large
   ```
4. Adjust language setting to match your speech

### Correction not working

1. Verify Ollama is installed and running:

   ```bash
   ollama list
   ```

2. Pull correction model:

   ```bash
   ollama pull gemma2:2b
   ```

3. Check daemon logs for errors

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/zandercpzed/dictate2me
cd dictate2me/plugins/obsidian-dictate2me

# Install dependencies
npm install

# Build
npm run build

# Or watch for changes
npm run dev
```

### Project Structure

```
obsidian-dictate2me/
├── src/
│   ├── main.ts       # Main plugin class
│   ├── client.ts     # WebSocket client
│   ├── settings.ts   # Settings interface
│   └── styles.css    # Icon animations
├── manifest.json     # Plugin metadata
├── package.json      # Dependencies
├── tsconfig.json     # TypeScript config
└── esbuild.config.mjs # Build config
```

## Privacy & Security

- ✅ **100% Offline**: All processing happens locally
- ✅ **No Cloud**: Audio never leaves your device
- ✅ **Local API**: Daemon runs on localhost only
- ✅ **Secure**: Token-based authentication
- ✅ **Open Source**: Full transparency

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/zandercpzed/dictate2me/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/zandercpzed/dictate2me/discussions)
- 📖 **Documentation**: [docs/](../../docs/)

## License

MIT License - see [LICENSE](../../LICENSE)

## Credits

- Built with ❤️ by the dictate2me team
- Powered by [Vosk](https://alphacephei.com/vosk/) (transcription)
- Powered by [Ollama](https://ollama.com/) (correction)
- Obsidian Plugin API

---

**Made with 🎤 for the Obsidian community**
