#!/bin/bash
# Script to download and setup Vosk C library for macOS

set -e

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

VOSK_VERSION="0.3.42"
LIB_DIR="lib"
VOSK_DIR="$LIB_DIR/vosk"

echo -e "${CYAN}Setting up Vosk C Library v${VOSK_VERSION}...${NC}"

# Create lib directory
mkdir -p "$LIB_DIR"

# Check if already installed
if [ -f "$VOSK_DIR/libvosk.dylib" ] && [ -f "$VOSK_DIR/vosk_api.h" ]; then
    echo -e "${GREEN}✓ Vosk library already installed in $VOSK_DIR${NC}"
    exit 0
fi

# Determine OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    ZIP_NAME="vosk-osx-${VOSK_VERSION}.zip"
    URL="https://github.com/alphacep/vosk-api/releases/download/v${VOSK_VERSION}/${ZIP_NAME}"
else
    echo -e "${RED}This script currently supports macOS only.${NC}"
    exit 1
fi

# Download
echo -e "Downloading $URL..."
curl -L -o "$LIB_DIR/$ZIP_NAME" "$URL"

# Extract
echo -e "Extracting..."
unzip -q "$LIB_DIR/$ZIP_NAME" -d "$LIB_DIR"

# Rename directory to standard name if needed (the zip usually extracts to vosk-osx-0.3.45)
EXTRACTED_DIR="$LIB_DIR/vosk-osx-${VOSK_VERSION}"
if [ -d "$EXTRACTED_DIR" ]; then
    rm -rf "$VOSK_DIR"
    mv "$EXTRACTED_DIR" "$VOSK_DIR"
fi

# Cleanup
rm "$LIB_DIR/$ZIP_NAME"

echo -e "${GREEN}✓ Vosk library installed to $VOSK_DIR${NC}"
echo -e ""
echo -e "To run tests/build, you need to set these environment variables:"
echo -e "export CGO_CFLAGS=\"-I$(pwd)/$VOSK_DIR\""
echo -e "export CGO_LDFLAGS=\"-L$(pwd)/$VOSK_DIR -lvosk\""
echo -e "export DYLD_LIBRARY_PATH=\"$(pwd)/$VOSK_DIR:\$DYLD_LIBRARY_PATH\""
