#!/usr/bin/env bash
set -euo pipefail

ROOTDIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOTDIR"

echo "Testing services..."

echo
echo "1) HTTP daemon health check"
if curl -sS http://127.0.0.1:8765/api/v1/health >/dev/null 2>&1; then
  echo "  ✓ daemon HTTP health: OK"
else
  echo "  ✗ daemon HTTP health: failed (is the daemon running?)"
fi

echo
echo "2) Vosk model load (attempt)"
python3 - <<'PY'
import os,sys
from pathlib import Path
model = Path('models/vosk-model-small-pt-0.3')
if model.exists():
    print('  ✓ model directory found at', model)
else:
    print('  ✗ model directory not found at', model)
    print('    Try: ./scripts/download-vosk-models.sh small')
PY

echo
echo "3) Ollama health check (if installed)"
if command -v ollama >/dev/null 2>&1; then
  if curl -sS http://127.0.0.1:11434/api/tags >/dev/null 2>&1; then
    echo "  ✓ Ollama reachable"
  else
    echo "  ✗ Ollama not reachable at http://127.0.0.1:11434 (is ollama running?)"
  fi
else
  echo "  ⚠ Ollama CLI not installed (skip). Install: brew install ollama" 
fi

echo
echo "4) PortAudio init test (build and run)"
go run ./cmd/check-portaudio || echo "  ✗ PortAudio initialization failed (install PortAudio)."

echo
echo "Service tests complete."
