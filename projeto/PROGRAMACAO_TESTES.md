# Programação de Testes - dictate2me

Este documento define a estratégia e o cronograma de testes para garantir a robustez antes e após o lançamento da v0.1.0.

## 🧪 Estratégia de Testes

### 1. Testes Unitários (Automated)

- **Frequência:** A cada commit/push (via CI) e durante o desenvolvimento local.
- **Escopo:** Lógica interna de pacotes (`internal/audio`, `internal/transcription`, `internal/correction`, `internal/api`).
- **Meta de Cobertura:** > 90% para lógica crítica.
- **Comando:** `make test`

### 2. Testes de Integração (Automated/Manual)

- **Frequência:** Antes de cada release ou merge na branch `main`.
- **Escopo:**
  - Fluxo completo: Audio -> Vosk -> Ollama -> Saída.
  - Comunicação Daemon <-> API Client.
- **Script de Validação:** `scripts/test-full.sh`

### 3. Testes End-to-End (Manual - Plugin)

- **Frequência:** Antes de releases que tocam o plugin ou API.
- **Cenários:**
  - Instalar plugin em um vault vazio.
  - Iniciar ditado via hotkey.
  - Ditar frases longas e curtas.
  - Testar desconexão de internet (não deve afetar, pois é offline).
  - Testar desligamento do Ollama durante o uso.

---

## 📅 Plano de Execução (v0.1.0 Release)

| Tipo      | Cenário                                            | Status      | Responsável |
| :-------- | :------------------------------------------------- | :---------- | :---------- |
| **Unit**  | `internal/audio`: Buffer overflow e recover        | ⏳ Pendente | Dev         |
| **Unit**  | `internal/api`: Autenticação e tokens inválidos    | ✅ Pronto   | Dev         |
| **Integ** | Script `test-full.sh` rodando limpo                | ⏳ Pendente | Dev         |
| **E2E**   | Verificar permissões de microfone no macOS         | ✅ Pronto   | User        |
| **E2E**   | Plugin: Inserção de texto com caracteres especiais | ⏳ Pendente | User        |

---

## 🐛 Casos de Teste Críticos (Smoke Tests)

1.  **Cold Start:** O sistema inicia corretamente na primeira vez (_fresh install_) sem travar baixando modelos?
2.  **No Microphone:** O sistema falha graciosamente se nenhum microfone for detectado?
3.  **Ollama Down:** O sistema continua transcrevendo (sem correção) se o Ollama estiver desligado?
4.  **Concurrency:** O que acontece se duas instâncias do plugin tentarem conectar ao mesmo daemon?
5.  **Long Session:** Ditar por mais de 5 minutos contínuos causa vazamento de memória?

## 📊 Matriz de Compatibilidade Alvo

- **OS:** macOS (Arm64/Intel) - **Prioridade 1**
- **OS:** Linux (Ubuntu/Fedora) - Prioridade 2
- **OS:** Windows 10/11 - Prioridade 3
- **Obsidian:** v1.5+
