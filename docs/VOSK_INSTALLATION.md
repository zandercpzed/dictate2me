# Instalação da Biblioteca Vosk

## Visão Geral

O dictate2me usa Vosk para transcrição de voz. Vosk requer uma biblioteca C nativa (`libvosk`) instalada no sistema.

## macOS

### Opção 1: Download de Binários Pré-compilados (Recomendado)

```bash
# Baixar a biblioteca pré-compilada
cd /tmp
curl -LO https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-osx-0.3.45.zip
unzip vosk-osx-0.3.45.zip

# Instalar biblioteca e headers
sudo cp vosk-osx-0.3.45/libvosk.dylib /usr/local/lib/
sudo cp vosk-osx-0.3.45/vosk_api.h /usr/local/include/

# Atualizar cache de bibliotecas
sudo update_dyld_shared_cache
```

### Opção 2: Compilar do Código-Fonte

```bash
# Instalar dependências
brew install cmake

# Clonar repositório
git clone https://github.com/alphacep/vosk-api.git
cd vosk-api/src

# Compilar
make

# Instalar
sudo cp libvosk.dylib /usr/local/lib/
sudo cp vosk_api.h /usr/local/include/
```

## Linux (Ubuntu/Debian)

```bash
# Baixar biblioteca pré-compilada
cd /tmp
curl -LO https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-x86_64-0.3.45.zip
unzip vosk-linux-x86_64-0.3.45.zip

# Instalar
sudo cp vosk-linux-x86_64-0.3.45/libvosk.so /usr/local/lib/
sudo cp vosk-linux-x86_64-0.3.45/vosk_api.h /usr/local/include/
sudo ldconfig
```

## Windows

```powershell
# Baixar biblioteca pré-compilada
curl -LO https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-win64-0.3.45.zip
Expand-Archive vosk-win64-0.3.45.zip -DestinationPath C:\vosk

# Adicionar ao PATH
setx PATH "%PATH%;C:\vosk\vosk-win64-0.3.45"
```

## Verificação da Instalação

Após instalar, verifique se a biblioteca está acessível:

```bash
# macOS/Linux
ls -la /usr/local/lib/libvosk.*
ls -la /usr/local/include/vosk_api.h

# Testar compilação
go build ./internal/transcription
```

## Variáveis de Ambiente (se necessário)

Se a biblioteca não for encontrada automaticamente:

```bash
# macOS/Linux
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lvosk"
export DYLD_LIBRARY_PATH="/usr/local/lib:$DYLD_LIBRARY_PATH"  # macOS
export LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"      # Linux

# Adicionar ao ~/.zshrc ou ~/.bashrc para tornar permanente
echo 'export DYLD_LIBRARY_PATH="/usr/local/lib:$DYLD_LIBRARY_PATH"' >> ~/.zshrc
```

## Troubleshooting

### Erro: "vosk_api.h file not found"

A biblioteca Vosk não está instalada ou não está no path de include.

**Solução:**

1. Verifique se `/usr/local/include/vosk_api.h` existe
2. Se não, siga as instruções de instalação acima
3. Configure `CGO_CFLAGS` se necessário

### Erro: "library not found for -lvosk"

A biblioteca compartilhada não está no path de bibliotecas.

**Solução:**

1. Verifique se `/usr/local/lib/libvosk.dylib` (macOS) ou `/usr/local/lib/libvosk.so` (Linux) existe
2. Configure `CGO_LDFLAGS` e `DYLD_LIBRARY_PATH`/`LD_LIBRARY_PATH`

### Erro em Runtime: "dyld: Library not loaded"

A biblioteca não está sendo encontrada em runtime.

**Solução:**

```bash
# macOS
export DYLD_LIBRARY_PATH="/usr/local/lib:$DYLD_LIBRARY_PATH"

# Linux
export LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"
```

## Links Úteis

- [Vosk API Releases](https://github.com/alphacep/vosk-api/releases)
- [Vosk Documentation](https://alphacephei.com/vosk/install)
- [Vosk GitHub](https://github.com/alphacep/vosk-api)
