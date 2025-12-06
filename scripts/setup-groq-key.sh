#!/usr/bin/env bash

# Setup Groq API Key for dictate2me

echo "🔑 Groq API Key Setup"
echo ""
echo "To use dictate2me with Groq's Whisper API, you need an API key."
echo ""
echo "1. Visit: https://console.groq.com/keys"
echo "2. Sign up (free) and create an API key"
echo "3. Copy the key and paste it below"
echo ""
read -p "Paste your Groq API key: " API_KEY

if [ -z "$API_KEY" ]; then
    echo "❌ No API key provided"
    exit 1
fi

# Add to shell profile
SHELL_PROFILE="$HOME/.zshrc"
if [ -f "$HOME/.bashrc" ]; then
    SHELL_PROFILE="$HOME/.bashrc"
fi

echo "" >> "$SHELL_PROFILE"
echo "# Groq API Key for dictate2me" >> "$SHELL_PROFILE"
echo "export GROQ_API_KEY='$API_KEY'" >> "$SHELL_PROFILE"

echo ""
echo "✓ API key saved to $SHELL_PROFILE"
echo ""
echo "To activate in current session:"
echo "  export GROQ_API_KEY='$API_KEY'"
echo ""
echo "Or restart your terminal."
