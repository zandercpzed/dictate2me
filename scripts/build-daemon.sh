#!/bin/bash
# Build script for dictate2me-daemon
# Handles CGO compilation with Vosk library

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "🔨 Building dictate2me-daemon..."
echo "📍 Project: $PROJECT_ROOT"

# Set CGO flags for Vosk
export CGO_CFLAGS="-I$PROJECT_ROOT/lib/vosk/include"
export CGO_LDFLAGS="-L$PROJECT_ROOT/lib/vosk -lvosk"
export DYLD_LIBRARY_PATH="$PROJECT_ROOT/lib/vosk:$DYLD_LIBRARY_PATH"

# Build
go build -o bin/dictate2me-daemon ./cmd/dictate2me-daemon

echo "✅ Build successful!"
echo "📦 Binary: $PROJECT_ROOT/bin/dictate2me-daemon"
