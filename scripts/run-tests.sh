#!/usr/bin/env bash
set -euo pipefail

# This script sets up env vars for CGO and runtime linking and runs go test
# Usage: ./scripts/run-tests.sh [go test args]

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VOSK_DIR="/tmp/dictate2me_vosk"

# Ensure symlink exists
ln -sf "$ROOT_DIR/lib/vosk" "$VOSK_DIR"

export CGO_CFLAGS="-I$VOSK_DIR"
export CGO_LDFLAGS="-L$VOSK_DIR"
export DYLD_LIBRARY_PATH="$VOSK_DIR:${DYLD_LIBRARY_PATH:-}"

echo "Running tests with CGO_CFLAGS=$CGO_CFLAGS CGO_LDFLAGS=$CGO_LDFLAGS DYLD_LIBRARY_PATH=$DYLD_LIBRARY_PATH"

if [ "$#" -eq 0 ]; then
  go test ./... -coverprofile=coverage.out -covermode=atomic
else
  go test "$@"
fi
