#!/bin/bash
# Script para configurar Ollama e modelos para correção de texto

set -e

# Cores
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}dictate2me - Ollama Setup${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Check OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
else
    echo -e "${RED}Unsupported OS: $OSTYPE${NC}"
    exit 1
fi

# Check if Ollama is installed
echo -e "${CYAN}[1/3] Checking Ollama installation...${NC}"
if command -v ollama &> /dev/null; then
    echo -e "${GREEN}✓ Ollama is already installed${NC}"
    OLLAMA_VERSION=$(ollama --version | head -n1)
    echo -e "  Version: $OLLAMA_VERSION"
else
    echo -e "${YELLOW}⚠ Ollama is not installed${NC}"
    echo -e "Installing Ollama..."
    
    if [[ "$OS" == "macos" ]]; then
        if command -v brew &> /dev/null; then
            brew install ollama
        else
            echo -e "${RED}Homebrew not found. Please install Homebrew first:${NC}"
            echo -e "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
            exit 1
        fi
    elif [[ "$OS" == "linux" ]]; then
        curl -fsSL https://ollama.com/install.sh | sh
    fi
    
    echo -e "${GREEN}✓ Ollama installed${NC}"
fi

# Check if Ollama is running
echo -e "${CYAN}[2/3] Checking Ollama daemon...${NC}"
if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Ollama daemon is running${NC}"
else
    echo -e "${YELLOW}⚠ Ollama daemon is not running${NC}"
    echo -e "Starting Ollama daemon..."
    
    if [[ "$OS" == "macos" ]]; then
        # On macOS, ollama serve runs in foreground, so we run it in background
        nohup ollama serve > /dev/null 2>&1 &
        sleep 2
    else
        # On Linux with systemd
        if command -v systemctl &> /dev/null; then
            sudo systemctl start ollama
        else
            nohup ollama serve > /dev/null 2>&1 &
            sleep 2
        fi
    fi
    
    # Verify it started
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Ollama daemon started${NC}"
    else
        echo -e "${RED}✗ Failed to start Ollama daemon${NC}"
        echo -e "  Try manually: ollama serve"
        exit 1
    fi
fi

# Pull the model
MODEL="gemma2:2b"
echo -e "${CYAN}[3/3] Downloading model: $MODEL${NC}"
echo -e "${YELLOW}This may take a few minutes (~1.7GB)...${NC}"

if ollama list | grep -q "$MODEL"; then
    echo -e "${GREEN}✓ Model $MODEL is already available${NC}"
else
    ollama pull "$MODEL"
    echo -e "${GREEN}✓ Model $MODEL downloaded${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Ollama setup complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${CYAN}Test correction:${NC}"
echo -e "  ollama run $MODEL 'Corrija: olá mundo como vai você'"
echo ""
echo -e "${CYAN}Start dictate2me:${NC}"
echo -e "  make run ARGS=\"start\""
echo ""
