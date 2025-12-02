#!/bin/bash

# Script de teste da API REST do dictate2me
# Testa todos os endpoints principais

set -e

echo "🧪 Testando API dictate2me"
echo "=========================="
echo ""

# Cores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Config
API_URL="http://localhost:8765/api/v1"
TOKEN_FILE="$HOME/.dictate2me/api-token"

# Função para printar sucesso
success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Função para printar erro
error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

# Função para printar info
info() {
    echo -e "${YELLOW}ℹ${NC} $1"
}

# 1. Verificar se daemon está rodando
echo "1. Verificando se daemon está rodando..."
if ! curl -s "$API_URL/health" > /dev/null 2>&1; then
    error "Daemon não está rodando. Execute: dictate2me-daemon"
fi
success "Daemon está rodando"
echo ""

# 2. Testar Health Check
echo "2. Testando health check..."
HEALTH_RESPONSE=$(curl -s "$API_URL/health")
echo "$HEALTH_RESPONSE" | jq . > /dev/null 2>&1 || error "Resposta inválida"
STATUS=$(echo "$HEALTH_RESPONSE" | jq -r '.status')
if [ "$STATUS" != "healthy" ]; then
    error "Status não é healthy: $STATUS"
fi
success "Health check OK"
echo "   Services:"
echo "$HEALTH_RESPONSE" | jq -r '.services | to_entries[] | "     - \(.key): \(.value)"'
echo ""

# 3. Verificar token
echo "3. Verificando token..."
if [ ! -f "$TOKEN_FILE" ]; then
    error "Token não encontrado em $TOKEN_FILE"
fi
TOKEN=$(cat "$TOKEN_FILE")
if [ -z "$TOKEN" ]; then
    error "Token vazio"
fi
success "Token encontrado: ${TOKEN:0:8}..."
echo ""

# 4. Testar autenticação
echo "4. Testando autenticação..."

# 4a. Sem token (deve falhar)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$API_URL/correct" \
    -H "Content-Type: application/json" \
    -d '{"text": "test"}')
if [ "$HTTP_CODE" != "401" ]; then
    error "Autenticação sem token deveria retornar 401, retornou $HTTP_CODE"
fi
success "Autenticação sem token retorna 401 (correto)"

# 4b. Token inválido (deve falhar)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
    -X POST "$API_URL/correct" \
    -H "Authorization: Bearer invalid-token" \
    -H "Content-Type: application/json" \
    -d '{"text": "test"}')
if [ "$HTTP_CODE" != "401" ]; then
    error "Token inválido deveria retornar 401, retornou $HTTP_CODE"
fi
success "Token inválido retorna 401 (correto)"
echo ""

# 5. Testar endpoint de correção
echo "5. Testando endpoint /correct..."
CORRECTION_RESPONSE=$(curl -s -X POST "$API_URL/correct" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"text": "olá mundo como vai você"}')

# Verificar se resposta é JSON válido
echo "$CORRECTION_RESPONSE" | jq . > /dev/null 2>&1 || error "Resposta inválida"

# Verificar campos
ORIGINAL=$(echo "$CORRECTION_RESPONSE" | jq -r '.original')
CORRECTED=$(echo "$CORRECTION_RESPONSE" | jq -r '.corrected')
MODEL=$(echo "$CORRECTION_RESPONSE" | jq -r '.model')

if [ "$ORIGINAL" == "null" ] || [ "$CORRECTED" == "null" ]; then
    # Pode ser que correção não esteja disponível
    ERROR_MSG=$(echo "$CORRECTION_RESPONSE" | jq -r '.error')
    if [[ "$ERROR_MSG" == *"not available"* ]]; then
        info "Correção não disponível (Ollama não configurado)"
    else
        error "Resposta de correção inválida: $CORRECTION_RESPONSE"
    fi
else
    success "Endpoint /correct funcionando"
    echo "   Original:  $ORIGINAL"
    echo "   Corrigido: $CORRECTED"
    echo "   Modelo:    $MODEL"
fi
echo ""

# 6. Testar rate limiting (opcional, comentado para não poluir)
# echo "6. Testando rate limiting..."
# info "Pulando teste de rate limiting (evitar poluir logs)"
# echo ""

echo "=========================="
echo -e "${GREEN}✓ Todos os testes passaram!${NC}"
echo ""
echo "API está funcionando corretamente! 🎉"
echo ""
echo "Próximos comandos úteis:"
echo "  - Ver logs do daemon:"
echo "    tail -f /var/log/dictate2me-daemon.log"
echo ""
echo "  - Testar WebSocket:"
echo "    websocat -H \"Authorization: Bearer $TOKEN\" ws://localhost:8765/api/v1/stream"
echo ""
echo "  - Ver documentação completa:"
echo "    cat docs/API.md"
