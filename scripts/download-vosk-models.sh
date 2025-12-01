#!/bin/bash
# Script para download dos modelos Vosk para português
# Uso: ./scripts/download-vosk-models.sh [small|large|both]

set -e

MODELS_DIR="models"
VOSK_BASE_URL="https://alphacephei.com/vosk/models"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Criar diretório de modelos se não existir
mkdir -p "$MODELS_DIR"

# Função para download e extração de modelo
download_model() {
    local model_name=$1
    local model_url=$2
    local model_size=$3
    
    print_info "Baixando modelo: $model_name ($model_size)"
    
    # Verificar se já existe
    if [ -d "$MODELS_DIR/$model_name" ]; then
        print_warn "Modelo $model_name já existe. Pulando..."
        return 0
    fi
    
    # Download
    local zip_file="$MODELS_DIR/${model_name}.zip"
    
    if [ ! -f "$zip_file" ]; then
        print_info "Fazendo download de $model_url..."
        curl -L -o "$zip_file" "$model_url" || {
            print_error "Falha ao baixar $model_name"
            return 1
        }
    else
        print_warn "Arquivo ZIP já existe, pulando download"
    fi
    
    # Extrair
    print_info "Extraindo $model_name..."
    unzip -q "$zip_file" -d "$MODELS_DIR" || {
        print_error "Falha ao extrair $model_name"
        return 1
    }
    
    # Remover ZIP
    rm "$zip_file"
    
    print_info "✓ Modelo $model_name instalado com sucesso!"
}

# Definir modelos disponíveis
SMALL_MODEL_NAME="vosk-model-small-pt-0.3"
SMALL_MODEL_URL="$VOSK_BASE_URL/${SMALL_MODEL_NAME}.zip"
SMALL_MODEL_SIZE="50MB"

LARGE_MODEL_NAME="vosk-model-pt-fb-v0.1.1-20220516_2113"
LARGE_MODEL_URL="$VOSK_BASE_URL/${LARGE_MODEL_NAME}.zip"
LARGE_MODEL_SIZE="1.6GB"

# Processar argumentos
MODE=${1:-small}

case $MODE in
    small)
        print_info "Baixando modelo pequeno (recomendado para uso geral)"
        download_model "$SMALL_MODEL_NAME" "$SMALL_MODEL_URL" "$SMALL_MODEL_SIZE"
        ;;
    large)
        print_info "Baixando modelo grande (melhor acurácia, mais lento)"
        download_model "$LARGE_MODEL_NAME" "$LARGE_MODEL_URL" "$LARGE_MODEL_SIZE"
        ;;
    both)
        print_info "Baixando ambos os modelos"
        download_model "$SMALL_MODEL_NAME" "$SMALL_MODEL_URL" "$SMALL_MODEL_SIZE"
        download_model "$LARGE_MODEL_NAME" "$LARGE_MODEL_URL" "$LARGE_MODEL_SIZE"
        ;;
    *)
        print_error "Modo inválido: $MODE"
        echo "Uso: $0 [small|large|both]"
        echo ""
        echo "Modelos disponíveis:"
        echo "  small - $SMALL_MODEL_NAME ($SMALL_MODEL_SIZE) - Recomendado"
        echo "  large - $LARGE_MODEL_NAME ($LARGE_MODEL_SIZE) - Melhor acurácia"
        echo "  both  - Baixa ambos os modelos"
        exit 1
        ;;
esac

print_info "Modelos instalados em: $MODELS_DIR/"
print_info ""
print_info "Para usar o modelo pequeno:"
print_info "  dictate2me --model=$MODELS_DIR/$SMALL_MODEL_NAME"
print_info ""
print_info "Para usar o modelo grande:"
print_info "  dictate2me --model=$MODELS_DIR/$LARGE_MODEL_NAME"
