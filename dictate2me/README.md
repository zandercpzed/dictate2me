# Dictate2Me — Plugin Obsidian

Transcrição de voz 100% local, privacidade total, integração nativa e correção inteligente.

## Instalação Manual

1. Clone ou baixe este repositório.
2. Execute:
   ```bash
   npm install
   npm run build
   ```
3. Copie os arquivos `manifest.json` e `dist/main.js` para a pasta de plugins do seu vault Obsidian:
   ```
   SeuVault/.obsidian/plugins/dictate2me/
   ```
4. Ative o plugin nas configurações do Obsidian.

## Uso

- Use o comando "Iniciar transcrição de voz" na paleta de comandos do Obsidian.
- O plugin simula captura, transcrição e correção, exibindo logs estruturados no console.

## Estrutura
- `src/main.ts` — entrypoint do plugin
- `src/audio.ts` — captura de áudio (simulado)
- `src/transcription.ts` — transcrição (simulada)
- `src/correction.ts` — correção (simulada)

## Roadmap
- MVP funcional
- Integração com interface do usuário
- Testes de integração
- Suporte a modelos reais (futuro)

# Dictate2Me Plugin

## Fluxo
- Captura áudio via WebAudio API
- Transcreve via API (`/api/v1/transcribe`) com fallback
- Corrige via API (`/api/v1/correct`) com fallback
- Insere texto no editor ativo
- StatusBar mostra progresso, sucesso, erro
- Comando para repetir última transcrição

## Logs
- Todos os módulos logam eventos e erros

## Testes
- Unitários em `src/*.test.ts`
- Manual: rodar plugin, verificar inserção e status

## Comandos
- Iniciar transcrição
- Repetir última transcrição
