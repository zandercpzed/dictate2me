#!/usr/bin/env bash
set -euo pipefail

# Create symlink to lib/vosk to avoid issues with spaces in path
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ln -sf "$ROOT_DIR/lib/vosk" /tmp/dictate2me_vosk

# Export CGO and runtime variables
export CGO_CFLAGS="-I/tmp/dictate2me_vosk"
export CGO_LDFLAGS="-L/tmp/dictate2me_vosk"
export DYLD_LIBRARY_PATH="/tmp/dictate2me_vosk:${DYLD_LIBRARY_PATH:-}"

echo "Test environment configured: /tmp/dictate2me_vosk -> $ROOT_DIR/lib/vosk"

echo "Run tests with: go test ./... -coverprofile=coverage.out -covermode=atomic"