#!/bin/bash

# Comprehensive test script for dictate2me
# Tests: build, daemon startup, API endpoints, and cleanup

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Config
API_URL="http://localhost:8765/api/v1"
DAEMON_PID=""

# Functions
success() {
    echo -e "${GREEN}✓${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
    cleanup
    exit 1
}

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

cleanup() {
    echo ""
    info "Cleaning up..."
    
    if [ -n "$DAEMON_PID" ] && ps -p "$DAEMON_PID" > /dev/null 2>&1; then
        info "Stopping daemon (PID: $DAEMON_PID)..."
        kill "$DAEMON_PID" 2>/dev/null || true
        wait "$DAEMON_PID" 2>/dev/null || true
        success "Daemon stopped"
    fi
}

# Trap cleanup on exit
trap cleanup EXIT INT TERM

echo "========================================"
echo "  Dictate2Me - Comprehensive Test Suite"
echo "========================================"
echo ""

# 1. Test Build
echo "1. Testing build..."
info "Running: make build"
if make build > /dev/null 2>&1; then
    success "Build successful"
else
    error "Build failed"
fi
echo ""

# 2. Check binaries
echo "2. Checking binaries..."
if [ -f "bin/dictate2me" ]; then
    success "dictate2me binary exists"
else
    error "dictate2me binary not found"
fi

if [ -f "bin/dictate2me-daemon" ]; then
    success "dictate2me-daemon binary exists"
else
    error "dictate2me-daemon binary not found"
fi
echo ""

# 3. Check Vosk model
echo "3. Checking Vosk model..."
if [ -d "models/vosk-model-small-pt-0.3" ]; then
    success "Vosk model found"
else
    warning "Vosk model not found (some tests may be skipped)"
    info "Download with: ./scripts/download-vosk-models.sh small"
fi
echo ""

# 4. Start daemon
echo "4. Starting daemon..."
info "Starting dictate2me-daemon in background..."

# Check if daemon is already running
if curl -s "$API_URL/health" > /dev/null 2>&1; then
    warning "Daemon already running, using existing instance"
else
    # Start daemon in background
    DYLD_LIBRARY_PATH=/tmp/dictate2me_vosk ./bin/dictate2me-daemon \
        --no-correction \
        > /tmp/dictate2me-test.log 2>&1 &
    DAEMON_PID=$!
    
    info "Daemon PID: $DAEMON_PID"
    info "Waiting for daemon to start..."
    
    # Wait for daemon to be ready (max 10 seconds)
    for i in {1..20}; do
        if curl -s "$API_URL/health" > /dev/null 2>&1; then
            success "Daemon started successfully"
            break
        fi
        sleep 0.5
        
        # Check if daemon crashed
        if ! ps -p "$DAEMON_PID" > /dev/null 2>&1; then
            error "Daemon crashed on startup. Check logs: /tmp/dictate2me-test.log"
        fi
        
        if [ $i -eq 20 ]; then
            error "Daemon failed to start. Check logs: /tmp/dictate2me-test.log"
        fi
    done
fi
echo ""

# 5. Test API health endpoint
echo "5. Testing API endpoints..."
info "GET /api/v1/health"

HEALTH_RESPONSE=$(curl -s "$API_URL/health")
if echo "$HEALTH_RESPONSE" | jq . > /dev/null 2>&1; then
    success "Health endpoint returns valid JSON"
    
    STATUS=$(echo "$HEALTH_RESPONSE" | jq -r '.status')
    if [ "$STATUS" = "healthy" ]; then
        success "Status is healthy"
    else
        error "Status is not healthy: $STATUS"
    fi
    
    # Show services status
    info "Services status:"
    echo "$HEALTH_RESPONSE" | jq -r '.services | to_entries[] | "  - \(.key): \(.value)"'
else
    error "Health endpoint returned invalid JSON"
fi
echo ""

# 6. Test authentication
echo "6. Testing authentication..."

# Get token
TOKEN_FILE="$HOME/.dictate2me/api-token"
if [ -f "$TOKEN_FILE" ]; then
    TOKEN=$(cat "$TOKEN_FILE")
    success "Token found"
else
    error "Token file not found: $TOKEN_FILE"
fi

# Test without token
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$API_URL/correct" \
    -H "Content-Type: application/json" \
    -d '{"text": "test"}')

if [ "$HTTP_CODE" = "401" ]; then
    success "Request without token correctly returns 401"
else
    error "Expected 401, got $HTTP_CODE"
fi

# Test with invalid token
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$API_URL/correct" \
    -H "Authorization: Bearer invalid-token" \
    -H "Content-Type: application/json" \
    -d '{"text": "test"}')

if [ "$HTTP_CODE" = "401" ]; then
    success "Request with invalid token correctly returns 401"
else
    error "Expected 401, got $HTTP_CODE"
fi
echo ""

# 7. Test correction endpoint (might fail if Ollama not available)
echo "7. Testing /correct endpoint..."
info "POST /api/v1/correct"

CORRECT_RESPONSE=$(curl -s -X POST "$API_URL/correct" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"text": "olá mundo"}')

if echo "$CORRECT_RESPONSE" | jq . > /dev/null 2>&1; then
    # Check if it's an error response
    ERROR=$(echo "$CORRECT_RESPONSE" | jq -r '.error // empty')
    
    if [ -n "$ERROR" ]; then
        if [[ "$ERROR" == *"not available"* ]]; then
            warning "Correction service not available (expected with --no-correction)"
        else
            error "Unexpected error: $ERROR"
        fi
    else
        success "Correct endpoint works"
        info "Response:"
        echo "$CORRECT_RESPONSE" | jq .
    fi
else
    error "Correct endpoint returned invalid JSON"
fi
echo ""

# 8. Test transcription endpoint (needs audio data)
echo "8. Testing /transcribe endpoint..."
info "POST /api/v1/transcribe"

# Create minimal valid audio data (silence)
# 16-bit PCM, 1 second at 16000 Hz
SILENCE_BYTES=$(python3 -c "import base64; print(base64.b64encode(b'\\x00' * 32000).decode())")

TRANSCRIBE_RESPONSE=$(curl -s -X POST "$API_URL/transcribe" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"audio\": \"$SILENCE_BYTES\"}")

if echo "$TRANSCRIBE_RESPONSE" | jq . > /dev/null 2>&1; then
    success "Transcribe endpoint returns valid JSON"
    
    TEXT=$(echo "$TRANSCRIBE_RESPONSE" | jq -r '.text // empty')
    info "Transcribed text: '$TEXT' (empty is ok for silence)"
else
    error "Transcribe endpoint returned invalid JSON"
fi
echo ""

# 9. Performance check
echo "9. Performance check..."
info "Measuring API latency..."

START_TIME=$(date +%s%N)
curl -s "$API_URL/health" > /dev/null
END_TIME=$(date +%s%N)
LATENCY=$(( (END_TIME - START_TIME) / 1000000 ))

if [ $LATENCY -lt 100 ]; then
    success "API latency: ${LATENCY}ms (excellent)"
elif [ $LATENCY -lt 500 ]; then
    success "API latency: ${LATENCY}ms (good)"
else
    warning "API latency: ${LATENCY}ms (slow)"
fi
echo ""

# Summary
echo "========================================"
echo -e "${GREEN}✓ All tests passed!${NC}"
echo "========================================"
echo ""
echo "Summary:"
echo "  - Build: ✓"
echo "  - Binaries: ✓"
echo "  - Daemon: ✓"
echo "  - Health: ✓"
echo "  - Auth: ✓"
echo "  - Endpoints: ✓"
echo "  - Performance: ${LATENCY}ms"
echo ""
echo "Daemon logs: /tmp/dictate2me-test.log"
echo ""

exit 0
