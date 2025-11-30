#!/bin/bash

# Script to download AI models for dictate2me
# Models are stored in the models/ directory

set -e

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

MODELS_DIR="models"
mkdir -p "$MODELS_DIR"

# Whisper Models (ggml format for whisper.cpp)
# Source: https://huggingface.co/ggerganov/whisper.cpp
WHISPER_MODEL="ggml-small.bin"
WHISPER_URL="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/$WHISPER_MODEL"

# LLM Models (GGUF format for llama.cpp)
# Using Phi-3 Mini 4k Instruct (Q4_K_M quantization)
# Source: https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf
LLM_MODEL="Phi-3-mini-4k-instruct-q4.gguf"
# Note: This is a placeholder URL, actual URL might vary slightly on HF
LLM_URL="https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf/resolve/main/Phi-3-mini-4k-instruct-q4.gguf"

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}dictate2me Model Downloader${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Function to download file if not exists
download_if_missing() {
    local file="$1"
    local url="$2"
    local path="$MODELS_DIR/$file"

    if [ -f "$path" ]; then
        echo -e "${GREEN}✓ $file already exists${NC}"
    else
        echo -e "${YELLOW}⬇️  Downloading $file...${NC}"
        echo -e "   URL: $url"
        
        if command -v curl &> /dev/null; then
            curl -L -o "$path" "$url" --progress-bar
        elif command -v wget &> /dev/null; then
            wget -O "$path" "$url" -q --show-progress
        else
            echo -e "${RED}Error: curl or wget not found${NC}"
            exit 1
        fi
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ Download complete${NC}"
        else
            echo -e "${RED}✗ Download failed${NC}"
            rm -f "$path" # Remove partial file
            exit 1
        fi
    fi
}

# 1. Download Whisper Model
echo -e "${CYAN}[1/2] Checking Whisper model (Transcription)...${NC}"
download_if_missing "$WHISPER_MODEL" "$WHISPER_URL"

# 2. Download LLM Model (Optional for now)
echo -e "${CYAN}[2/2] Checking LLM model (Correction)...${NC}"
echo -e "${YELLOW}Skipping LLM download for Phase 2 (Transcription focus)${NC}"
# download_if_missing "$LLM_MODEL" "$LLM_URL"

echo ""
echo -e "${GREEN}All required models are ready!${NC}"
echo -e "Models location: $(pwd)/$MODELS_DIR"
