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

# Start daemon
cd "$PROJECT_ROOT"
echo "🎤 Starting daemon..."
./bin/dictate2me-daemon

# Note: This script will block until daemon is stopped (Ctrl+C)
