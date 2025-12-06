#!/bin/bash
# Launcher script for dictate2me daemon

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🚀 Starting dictate2me daemon..."
echo "📍 Project: $PROJECT_ROOT"

# Set library path for Vosk
export DYLD_LIBRARY_PATH="$PROJECT_ROOT/lib/vosk:$DYLD_LIBRARY_PATH"

# Check if daemon is already running
if curl -s http://localhost:8765/api/v1/health > /dev/null 2>&1; then
    echo "⚠️  Daemon is already running!"
    echo "✅ Health check: OK"
    exit 0
fi

# --- OLLAMA SETUP ---
echo "🦙 Checking Ollama setup..."

if ! command -v ollama &> /dev/null; then
    echo "❌ Ollama not found. Please install it from https://ollama.com"
    # We continue anyway, but correction won't work
else
    # Check if Ollama server is running
    if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "🔄 Starting Ollama server..."
        ollama serve > /dev/null 2>&1 &
        
        # Wait for Ollama to start
        echo "⏳ Waiting for Ollama..."
        count=0
        while ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; do
            sleep 1
            count=$((count+1))
            if [ $count -ge 10 ]; then
                echo "⚠️  Ollama took too long to start. Correction might be disabled."
                break
            fi
        done
        echo "✅ Ollama started!"
    else
        echo "✅ Ollama is already running."
    fi

    # Check for model (using llama3 as default for correction)
    MODEL="llama3"
    if curl -s http://localhost:11434/api/tags | grep -q "$MODEL"; then
        echo "✅ Model $MODEL found."
    else
        echo "⬇️  Downloading model $MODEL (this may take a while)..."
        ollama pull $MODEL
        echo "✅ Model $MODEL ready."
    fi
fi
# --------------------

# Load Groq API Key from .env if exists
ENV_FILE="$HOME/.dictate2me/.env"
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
    echo "✓ Loaded API key from ~/.dictate2me/.env"
fi

# Check for Groq API Key
if [ -z "$GROQ_API_KEY" ]; then
    echo ""
    echo "❌ GROQ_API_KEY not configured!"
    echo ""
    echo "To fix (choose one):"
    echo "  Option 1 (Obsidian Plugin):"
    echo "    Go to Settings > Dictate2Me > API Configuration"
    echo "    Enter your Groq API key there"
    echo ""
    echo "  Option 2 (Manual):"
    echo "    1. Get API key at: https://console.groq.com/keys"
    echo "    2. Add to your ~/.zshrc:"
    echo "       export GROQ_API_KEY='your-key-here'"
    echo "    3. Run: source ~/.zshrc"
    echo ""
    exit 1
fi

# Start daemon
cd "$PROJECT_ROOT"
echo "🎤 Starting daemon..."
exec ./bin/dictate2me-daemon > /tmp/dictate2me-daemon.log 2>&1

# Note: This script will block until daemon is stopped (Ctrl+C)
