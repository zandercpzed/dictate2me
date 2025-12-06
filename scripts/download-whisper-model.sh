#!/usr/bin/env bash

# Script to download Whisper model for dictate2me
# Models available: tiny, base, small, medium, large
# We use 'small' for good balance of speed vs accuracy

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MODELS_DIR="$PROJECT_ROOT/models"

MODEL_NAME="${1:-small}"
MODEL_FILE="ggml-$MODEL_NAME.bin"
MODEL_URL="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/$MODEL_FILE"

echo "📥 Downloading Whisper model: $MODEL_NAME"
echo "   URL: $MODEL_URL"
echo ""

# Create models directory
mkdir -p "$MODELS_DIR"

# Download with progress
cd "$MODELS_DIR"
if [ -f "whisper-$MODEL_NAME.bin" ]; then
    echo "✓ Model already exists: whisper-$MODEL_NAME.bin"
    exit 0
fi

echo "Downloading... (this may take a few minutes)"
curl -L -o "whisper-$MODEL_NAME.bin" "$MODEL_URL"

echo ""
echo "✓ Model downloaded successfully!"
echo "  Path: $MODELS_DIR/whisper-$MODEL_NAME.bin"
echo ""
echo "Available models (in order of size):"
echo "  tiny   - ~75 MB  (fastest, lowest accuracy)"
echo "  base   - ~142 MB"
echo "  small  - ~466 MB (recommended, good balance)"
echo "  medium - ~1.5 GB"
echo "  large  - ~2.9 GB (slowest, best accuracy)"
echo ""
echo "To download a different model: $0 [tiny|base|small|medium|large]"
