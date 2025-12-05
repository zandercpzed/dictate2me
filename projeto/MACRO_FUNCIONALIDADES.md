# Macro-Funcionalidades - dictate2me

Este documento lista as principais capacidades funcionais do sistema `dictate2me`.

## 1. Captura e Processamento de Áudio

**Status: ✅ Concluído**

- **Captura de Microfone:** Capacidade de capturar áudio bruto diretamente do microfone padrão do sistema usando PortAudio.
- **Detecção de Atividade de Voz (VAD):** Algoritmo inteligente para detectar quando o usuário está falando e quando há silêncio, otimizando o processamento.
- **Buffer Circular:** Gerenciamento eficiente de memória para fluxo contínuo de áudio sem perdas.
- **Monitoramento de Nível:** Feedback visual ou métrico do nível de entrada de áudio (input level).

## 2. Transcrição Offline (Speech-to-Text)

**Status: ✅ Concluído**

- **Motor Vosk:** Utilização da biblioteca Vosk para reconhecimento de fala offline, garantindo privacidade e independência de internet.
- **Gestão de Modelos:** Scripts e lógica para baixar e carregar modelos de idioma (padrão: Português-BR Small).
- **Streaming em Tempo Real:** Processamento do áudio à medida que é capturado, fornecendo resultados parciais instantâneos.
- **Suporte Multi-idioma:** Arquitetura pronta para suportar outros idiomas bastando trocar o modelo do Vosk.

## 3. Correção e Pós-Processamento (AI)

**Status: ✅ Concluído**

- **Motor Ollama (Local LLM):** Integração com modelos LLM locais (Gemma 2, Llama 3, etc.) rodando via Ollama.
- **Correção Contextual:** O texto transcrito é enviado ao LLM para correção de pontuação, gramática e coerência, transformando a "fala" em "texto escrito".
- **Prompt Engineering:** Prompts otimizados para garantir que o estilo original seja mantido enquanto a forma é corrigida.
- **Health Checks:** Verificação automática da disponibilidade do serviço Ollama.

## 4. API e Conectividade

**Status: ✅ Concluído**

- **Daemon de Serviço:** Processo de fundo (`dictate2me-daemon`) que mantém os modelos carregados e prontos para uso.
- **REST API:** Endpoints HTTP padrão para controle (start/stop) e status.
- **WebSocket Streaming:** Canal bidirecional para envio de áudio e recebimento de texto transcrito/corrigido em tempo real com latência mínima.
- **Segurança:** Autenticação via Token Bearer gerado localmente.

## 5. Integração com Obsidian

**Status: ✅ Concluído**

- **Plugin Nativo:** Extensão desenvolvida em TypeScript para o ecossistema Obsidian.
- **Inserção Direta:** O texto ditado aparece diretamente na nota ativa, na posição do cursor.
- **Controle Integrado:** Comandos (Command Palette) e atalhos de teclado para iniciar/parar o ditado.
- **Feedback Visual:** Barra de status e animações indicando que a gravação está ativa.
- **Gerenciamento do Daemon:** O plugin tenta iniciar automaticamente o daemon se ele não estiver rodando.

## 6. Ferramentas de Linha de Comando (CLI)

**Status: ✅ Concluído**

- **Comando `dictate2me`:** Utilitário para testar e usar o sistema diretamente do terminal.
- **Flags de Configuração:** Opções para escolher modelos, desativar correção, definir dispositivos de áudio, etc.
- **Output Formatado:** Saída visualmente clara com distinção entre resultados parciais, finais e corrigidos.
