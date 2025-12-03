#!/bin/bash
# Script para instalar o plugin dictate2me em uma vault do Obsidian

set -e

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔌 Instalador do Plugin dictate2me${NC}"
echo ""

# Verificar argumentos
if [ -z "$1" ]; then
    echo -e "${RED}❌ Erro: Caminho da vault não fornecido${NC}"
    echo ""
    echo "Uso: $0 /caminho/para/sua/vault"
    echo ""
    echo "Exemplo:"
    echo "  $0 ~/Documents/MinhaNotes"
    exit 1
fi

VAULT_PATH="$1"
PLUGIN_SOURCE="$(dirname "$0")/../plugins/obsidian-dictate2me"
PLUGIN_DEST="$VAULT_PATH/.obsidian/plugins/dictate2me"

# Verificar se a vault existe
if [ ! -d "$VAULT_PATH" ]; then
    echo -e "${RED}❌ Erro: Vault não encontrada: $VAULT_PATH${NC}"
    exit 1
fi

# Criar diretório .obsidian se não existir
if [ ! -d "$VAULT_PATH/.obsidian" ]; then
    echo -e "${BLUE}📂 Criando diretório .obsidian...${NC}"
    mkdir -p "$VAULT_PATH/.obsidian"
fi

# Criar diretório de plugins se não existir
if [ ! -d "$VAULT_PATH/.obsidian/plugins" ]; then
    echo -e "${BLUE}📂 Criando diretório de plugins...${NC}"
    mkdir -p "$VAULT_PATH/.obsidian/plugins"
fi

# Criar diretório do plugin
echo -e "${BLUE}📂 Criando diretório do plugin dictate2me...${NC}"
mkdir -p "$PLUGIN_DEST"

# Copiar arquivos do plugin
echo -e "${BLUE}📋 Copiando arquivos do plugin...${NC}"

if [ ! -f "$PLUGIN_SOURCE/main.js" ]; then
    echo -e "${RED}❌ Erro: main.js não encontrado. Execute 'npm run build' primeiro.${NC}"
    exit 1
fi

cp "$PLUGIN_SOURCE/main.js" "$PLUGIN_DEST/"
cp "$PLUGIN_SOURCE/manifest.json" "$PLUGIN_DEST/"
cp "$PLUGIN_SOURCE/styles.css" "$PLUGIN_DEST/"

echo -e "${GREEN}✅ Arquivos copiados:${NC}"
ls -lh "$PLUGIN_DEST"

# Adicionar plugin à lista de plugins habilitados
COMMUNITY_PLUGINS="$VAULT_PATH/.obsidian/community-plugins.json"

if [ ! -f "$COMMUNITY_PLUGINS" ]; then
    echo -e "${BLUE}📝 Criando community-plugins.json...${NC}"
    echo '["dictate2me"]' > "$COMMUNITY_PLUGINS"
else
    echo -e "${BLUE}📝 Atualizando community-plugins.json...${NC}"
    # Adicionar dictate2me se não existir
    if ! grep -q "dictate2me" "$COMMUNITY_PLUGINS"; then
        # Usar jq se disponível, senão fazer manualmente
        if command -v jq &> /dev/null; then
            jq '. + ["dictate2me"]' "$COMMUNITY_PLUGINS" > "${COMMUNITY_PLUGINS}.tmp"
            mv "${COMMUNITY_PLUGINS}.tmp" "$COMMUNITY_PLUGINS"
        else
            # Fallback simples
            sed -i.bak 's/\]$/, "dictate2me"]/' "$COMMUNITY_PLUGINS"
            rm "${COMMUNITY_PLUGINS}.bak"
        fi
    fi
fi

echo ""
echo -e "${GREEN}✅ Plugin instalado com sucesso!${NC}"
echo ""
echo -e "${BLUE}📍 Localização:${NC} $PLUGIN_DEST"
echo ""
echo -e "${BLUE}🔑 Token da API:${NC}"
if [ -f ~/.dictate2me/api-token ]; then
    cat ~/.dictate2me/api-token
    echo ""
else
    echo -e "${RED}Token não encontrado. Inicie o daemon primeiro.${NC}"
fi
echo ""
echo -e "${BLUE}📖 Próximos passos:${NC}"
echo "1. Abra ou recarregue a vault no Obsidian (Cmd+R / Ctrl+R)"
echo "2. Vá em Settings → Community plugins"
echo "3. Habilite 'dictate2me' se necessário"
echo "4. Configure o token nas settings do plugin"
echo "5. Certifique-se de que o daemon está rodando:"
echo "   ./bin/dictate2me-daemon"
echo ""
echo -e "${GREEN}🎉 Pronto para usar!${NC}"
