#!/bin/bash

# dictate2me Development Environment Setup Script
# This script sets up the development environment for dictate2me

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}dictate2me Development Setup${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo -e "${YELLOW}Warning: This script is optimized for macOS.${NC}"
    echo -e "${YELLOW}Some features may not work on other platforms.${NC}"
    echo ""
fi

# Check Go installation
echo -e "${CYAN}[1/7] Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    echo -e "Please install Go 1.23+ from https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓ Go ${GO_VERSION} is installed${NC}"

# Check if Go version is 1.23+
REQUIRED_VERSION="1.23"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo -e "${RED}✗ Go version must be 1.23 or higher${NC}"
    exit 1
fi

# Install golangci-lint
echo -e "${CYAN}[2/7] Installing golangci-lint...${NC}"
if ! command -v golangci-lint &> /dev/null; then
    if command -v brew &> /dev/null; then
        brew install golangci-lint
    else
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    fi
    echo -e "${GREEN}✓ golangci-lint installed${NC}"
else
    echo -e "${GREEN}✓ golangci-lint already installed${NC}"
fi

# Install goimports
echo -e "${CYAN}[3/7] Installing goimports...${NC}"
if ! command -v goimports &> /dev/null; then
    go install golang.org/x/tools/cmd/goimports@latest
    echo -e "${GREEN}✓ goimports installed${NC}"
else
    echo -e "${GREEN}✓ goimports already installed${NC}"
fi

# Install air (hot reload)
echo -e "${CYAN}[4/7] Installing air (hot reload)...${NC}"
if ! command -v air &> /dev/null; then
    go install github.com/air-verse/air@latest
    echo -e "${GREEN}✓ air installed${NC}"
else
    echo -e "${GREEN}✓ air already installed${NC}"
fi

# Install govulncheck
echo -e "${CYAN}[5/7] Installing govulncheck...${NC}"
if ! command -v govulncheck &> /dev/null; then
    go install golang.org/x/vuln/cmd/govulncheck@latest
    echo -e "${GREEN}✓ govulncheck installed${NC}"
else
    echo -e "${GREEN}✓ govulncheck already installed${NC}"
fi

# Install PortAudio (System Dependency)
echo -e "${CYAN}[6/8] Installing PortAudio...${NC}"
if [[ "$OSTYPE" == "darwin"* ]]; then
    if ! brew list portaudio &> /dev/null; then
        brew install portaudio
        echo -e "${GREEN}✓ PortAudio installed via Homebrew${NC}"
    else
        echo -e "${GREEN}✓ PortAudio already installed${NC}"
    fi
    
    # Check for pkg-config (needed for CGO)
    if ! command -v pkg-config &> /dev/null; then
        brew install pkg-config
        echo -e "${GREEN}✓ pkg-config installed${NC}"
    fi
else
    # Linux assumption (Debian/Ubuntu)
    if command -v apt-get &> /dev/null; then
        sudo apt-get update
        sudo apt-get install -y portaudio19-dev pkg-config
        echo -e "${GREEN}✓ PortAudio installed via apt${NC}"
    else
        echo -e "${YELLOW}⚠ Please install PortAudio manually for your OS${NC}"
    fi
fi

# Download Go dependencies
echo -e "${CYAN}[7/8] Downloading Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies downloaded${NC}"

# Setup pre-commit hooks (if pre-commit is installed)
echo -e "${CYAN}[8/8] Setting up pre-commit hooks...${NC}"
if command -v pre-commit &> /dev/null; then
    pre-commit install
    echo -e "${GREEN}✓ Pre-commit hooks installed${NC}"
else
    echo -e "${YELLOW}⚠ pre-commit not found, skipping hooks setup${NC}"
    echo -e "${YELLOW}  Install with: pip install pre-commit${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Development environment setup complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${CYAN}Next steps:${NC}"
echo -e "  1. Download AI models: ${YELLOW}make models${NC}"
echo -e "  2. Run tests: ${YELLOW}make test${NC}"
echo -e "  3. Build the app: ${YELLOW}make build${NC}"
echo -e "  4. Start coding! 🚀"
echo ""
