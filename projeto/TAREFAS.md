# Lista de Tarefas - dictate2me

## 🏁 Fase Atual: Preparação para Lançamento (v0.1.0)

### ✅ Concluído (MVP + Core)

- [x] **Arquitetura Base**

  - [x] Estrutura do projeto (Standard Go Project Layout)
  - [x] Gerenciamento de dependências (Go Modules)
  - [x] Sistema de Build (Makefile)

- [x] **Módulo de Áudio**

  - [x] Integração PortAudio
  - [x] Detecção de Voz (VAD)
  - [x] Buffer Circular

- [x] **Transcrição (Speech-to-Text)**

  - [x] Engine Vosk implementado (Offline)
  - [x] Suporte a download de modelos
  - [x] Streaming de áudio real-time

- [x] **Correção de Texto (AI)**

  - [x] Integração Ollama
  - [x] Prompt engineering para correção
  - [x] Fallback gracioso (sem correção se offline)

- [x] **API & Daemon**

  - [x] Servidor REST Local
  - [x] WebSocket para streaming
  - [x] Daemon de buackground (`dictate2me-daemon`)
  - [x] Autenticação via Token

- [x] **Integrações**
  - [x] Plugin Obsidian (TypeScript)
  - [x] CLI (`dictate2me`)

### 🚧 Em Andamento / Próximos Passos Imediatos

#### 1. Qualidade & Testes

- [ ] Aumentar cobertura de testes unitários para 100% (Atualmente ~85%)
- [ ] Executar smoke test final em ambiente limpo
- [ ] Validar instalação em uma máquina "virgem" (simulada)

#### 2. Release Engineering

- [ ] Verificar CI/CD Workflows (`.github/workflows/ci.yaml`)
- [ ] Criar Release Tag `v0.1.0` no Git
- [ ] Publicar Release no GitHub com binários pré-compilados (se aplicável)
- [ ] Verificar documentação no GitHub (Links quebrados, imagens)

#### 3. Organização

- [ ] Manter raiz do repositório limpa (Apenas arquivos essenciais)
- [ ] Revisar metadados do repositório (Topics, Description, Website)

### 🔮 Futuro (v0.2.0+)

- [ ] **Suporte Multiplataforma Nativo**
  - [ ] Validar e empacotar para Linux (.deb/.rpm)
  - [ ] Validar e empacotar para Windows (.exe/MSI)
- [ ] **Novas Integrações**
  - [ ] Plugin VS Code
  - [ ] Extensão Chrome/Firefox
- [ ] **Melhorias de Engine**

  - [ ] Suporte a modelos Whisper locais (via whisper.cpp - reavaliar)
  - [ ] Ajuste fino do modelo Ollama para terminologias específicas

- [ ] **Interface Gráfica**
  - [ ] Dashboard Web local para configuração e logs
  - [ ] Tray icon para controle rápido do daemon
