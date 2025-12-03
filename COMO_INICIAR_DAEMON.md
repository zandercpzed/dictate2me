# 🚀 Como Iniciar o Daemon dictate2me

## Método 1: Via Terminal (Recomendado)

### Passo a Passo:

1. **Abra o Terminal:**

   - Pressione `⌘ + Espaço` (Spotlight)
   - Digite "Terminal" e pressione Enter

2. **Cole e execute este comando:**

   ```bash
   cd '/Users/zander/Library/CloudStorage/GoogleDrive-zander.cattapreta@zedicoes.com/My Drive/_ programação/_ dictate2me/dictate2me' && ./scripts/start-daemon.sh
   ```

3. **Aguarde a mensagem:**

   ```
   🚀 Starting dictate2me daemon...
   🎤 Daemon started successfully!
   ```

4. **Mantenha o Terminal aberto** enquanto usar o plugin

---

## Método 2: Via Plugin do Obsidian

1. **No Obsidian:**
   - Settings → Dictate2Me
2. **Na seção "Daemon Control":**

   - Clique no botão **"Copy Command"**
   - Um modal aparecerá com instruções

3. **Abra o Terminal e cole** (`⌘ + V`)

4. **Pressione Enter**

---

## Verificar se está Rodando

No Terminal:

```bash
curl http://localhost:8765/api/v1/health
```

Deve retornar:

```json
{
  "status": "healthy",
  "services": { "transcription": "ready", "correction": "disabled" }
}
```

---

## Parar o Daemon

No Terminal onde está rodando:

```
Ctrl + C
```

---

## Troubleshooting

### "Connection failed"

- ✅ Verifique se o daemon está rodando
- ✅ Use o comando curl acima para testar

### "Library not loaded: libvosk.dylib"

- O script `start-daemon.sh` já configura isso automaticamente
- Se persistir, execute manualmente:
  ```bash
  cd '/Users/zander/Library/CloudStorage/GoogleDrive-zander.cattapreta@zedicoes.com/My Drive/_ programação/_ dictate2me/dictate2me'
  export DYLD_LIBRARY_PATH="$(pwd)/lib/vosk:$DYLD_LIBRARY_PATH"
  ./bin/dictate2me-daemon
  ```

### Daemon já está rodando

- Você verá: "⚠️ Daemon is already running!"
- Tudo certo, pode usar normalmente!

---

## Token da API

Cole este token nas configurações do plugin:

```
60c441335af1d363447ee95d6817834f4aeec3600c0c054aec385808f2c6ca11
```

---

**Pronto!** Agora você pode usar o ditado no Obsidian! 🎤
