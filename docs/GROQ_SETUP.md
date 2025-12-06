# Configuração da API Groq

O `dictate2me` usa a API do Groq para transcrição de voz (Whisper-large-v3).

## Como Configurar

### 1. Obter API Key (Grátis)

1. Acesse: https://console.groq.com/keys
2. Faça login ou cadastro (é gratuito)
3. Clique em "Create API Key"
4. Copie a chave gerada

### 2. Configurar no Sistema

Adicione a chave ao seu arquivo `~/.zshrc` (ou `~/.bashrc` se usar Bash):

```bash
echo "export GROQ_API_KEY='sua-chave-aqui'" >> ~/.zshrc
source ~/.zshrc
```

**Importante:** Substitua `'sua-chave-aqui'` pela chave real que você copiou.

### 3. Verificar

```bash
echo $GROQ_API_KEY
```

Deve exibir sua chave.

### 4. Iniciar o Daemon

```bash
./scripts/start-daemon.sh
```

## Limites do Plano Gratuito

- **Whisper**: 10 requisições/minuto
- **Suficiente** para uso pessoal de ditado

## Troubleshooting

**Erro: "GROQ_API_KEY environment variable not set"**

- Execute: `source ~/.zshrc` no terminal atual
- Ou abra um novo terminal

**Erro: "invalid API key"**

- Verifique se copiou a chave completa
- Gere uma nova chave em: https://console.groq.com/keys
