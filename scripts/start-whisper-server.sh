#!/usr/bin/env bash

# Start Whisper HTTP Server for dictate2me
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

WHISPER_BIN="/tmp/whisper.cpp/bin/whisper-server"
MODEL_PATH="$PROJECT_ROOT/models/whisper-small.bin"
HOST="127.0.0.1"
PORT="8766"

if [ ! -f "$WHISPER_BIN" ]; then
    echo "❌ Whisper server not found. Please run: make install-whisper"
    exit 1
fi

if [ ! -f "$MODEL_PATH" ]; then
    echo "❌ Whisper model not found at: $MODEL_PATH"
    echo "   Run: make download-whisper-model"
    exit 1
fi

echo "🎙️  Starting Whisper Server..."
echo "   Model: $MODEL_PATH"
echo "   Endpoint: http://$HOST:$PORT"
echo ""

exec "$WHISPER_BIN" \
    --model "$MODEL_PATH" \
    --host "$HOST" \
    --port "$PORT" \
    --language pt \
    --threads 4
