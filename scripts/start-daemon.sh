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

# Start daemon
cd "$PROJECT_ROOT"
echo "🎤 Starting daemon..."
./bin/dictate2me-daemon

# Note: This script will block until daemon is stopped (Ctrl+C)
