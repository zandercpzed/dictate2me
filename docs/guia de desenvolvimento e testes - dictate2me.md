# dictate2me — Guia de Desenvolvimento AI-Driven

> Desenvolvimento 100% via IA com validação ao vivo em cada fase

---

## 1. Metodologia: AI-Driven Development (ADD)

### 1.1 Ciclo de Desenvolvimento

```
┌─────────────────────────────────────────────────────────────────┐
│                     CICLO POR MÓDULO                            │
│                                                                 │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐  │
│   │  PROMPT  │───▶│  GERAR   │───▶│  TESTAR  │───▶│ VALIDAR  │  │
│   │          │    │  CÓDIGO  │    │ AO VIVO  │    │ OUTPUT   │  │
│   └──────────┘    └──────────┘    └──────────┘    └────┬─────┘  │
│        ▲                                               │        │
│        │              ┌──────────┐                     │        │
│        └──────────────│ CORRIGIR │◀────────────────────┘        │
│           (se falha)  └──────────┘      (se sucesso: próximo)   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Princípios

| Princípio             | Implementação                                   |
| --------------------- | ----------------------------------------------- |
| **Zero Human Coding** | Todo código via LLM no VS Code                  |
| **Live Validation**   | Cada módulo testado ao vivo antes de prosseguir |
| **Self-Correcting**   | IA analisa erros e corrige automaticamente      |
| **Observable Output** | Logs estruturados para IA interpretar           |
| **Fail Fast**         | Testes executados imediatamente após geração    |

### 1.3 Estrutura do Projeto

```
dictate2me/
├── src/
│   ├── main.ts                 # Entry point do plugin
│   ├── settings.ts             # Configurações do usuário
│   ├── core/
│   │   ├── AudioCapture.ts     # Web Audio API wrapper
│   │   ├── VADProcessor.ts     # Silero VAD integration
│   │   ├── Transcriber.ts      # Whisper via Transformers.js
│   │   └── TextProcessor.ts    # Pipeline de pós-processamento
│   ├── ui/
│   │   ├── StatusBar.ts        # Indicador de gravação
│   │   └── SettingsTab.ts      # Painel de configurações
│   └── utils/
│       ├── AudioUtils.ts       # Resample, conversão de formatos
│       └── Logger.ts           # Logging estruturado para IA
├── tests/
│   ├── unit/                   # Testes unitários
│   ├── live/                   # Testes ao vivo (executáveis)
│   └── fixtures/               # Áudio e texto de teste
├── scripts/
│   ├── validate-phase.ts       # Validador por fase
│   ├── live-test.ts            # Executor de testes ao vivo
│   └── download-models.ts      # Download de modelos
├── .ai/
│   ├── prompts/                # Prompts por fase
│   ├── context.md              # Contexto do projeto para IA
│   └── validation-rules.json   # Regras de validação
├── manifest.json
├── package.json
├── tsconfig.json
├── esbuild.config.mjs
└── vitest.config.ts
```

---

## 2. Sistema de Prompts por Fase

### 2.0 Contexto Base (sempre incluir)

```markdown
# CONTEXTO DO PROJETO

## Projeto: dictate2me

Plugin Obsidian para transcrição de voz 100% local.

## Stack Técnico

- TypeScript + esbuild (bundle único)
- @ricky0123/vad-web (Silero VAD para detecção de fala)
- @huggingface/transformers (Whisper para STT)
- Obsidian Plugin API

## Restrições

- Sem módulos nativos (cross-platform)
- Bundle único main.js
- RAM máxima ~500MB durante transcrição
- Funcionar offline 100%

## Padrões de Código

- Funções puras quando possível
- Async/await (não callbacks)
- Erros descritivos com contexto
- Logs estruturados: console.log('[MÓDULO]', { dados })
```

---

### FASE 0: Setup do Projeto

#### Prompt 0.1 — Estrutura Base

```markdown
# TAREFA: Criar estrutura base do plugin dictate2me

## Arquivos a criar:

1. `package.json` com:
    - name: "dictate2me"
    - scripts: dev, build, test
    - dependencies: @huggingface/transformers, @ricky0123/vad-web
    - devDependencies: typescript, esbuild, vitest, obsidian

2. `tsconfig.json` para Obsidian plugin

3. `esbuild.config.mjs` com:
    - Bundle para main.js
    - External: obsidian, electron
    - Watch mode para dev

4. `manifest.json` com:
    - id: "dictate2me"
    - name: "Dictate2Me"
    - version: "0.1.0"
    - minAppVersion: "1.4.0"

5. `src/main.ts` com:
    - Classe DictatePlugin extends Plugin
    - onload(): console.log('[dictate2me] Plugin loaded')
    - Comando básico registrado

## Após criar, execute:

npm install && npm run build

## Validação esperada:

- ✓ main.js gerado sem erros
- ✓ Tamanho do bundle < 100KB (sem modelos)
```

#### Teste ao Vivo 0.1

```bash
# Executar após gerar código:
npm run build 2>&1 | tee build.log

# Validar output (IA deve verificar):
# - "Build completed" ou similar
# - Arquivo main.js existe
# - Sem "error" no log
```

#### Prompt 0.2 — Validação de Carga no Obsidian

```markdown
# TAREFA: Criar script de validação de plugin

Criar `scripts/validate-plugin.ts`:

1. Verifica se main.js existe e é válido
2. Verifica se manifest.json tem campos obrigatórios
3. Simula carga do módulo (import dinâmico)
4. Retorna JSON estruturado:

{
"phase": 0,
"status": "pass" | "fail",
"checks": [
{ "name": "main.js exists", "pass": true },
{ "name": "manifest valid", "pass": true },
{ "name": "module loads", "pass": true }
],
"errors": []
}
```

---

### FASE 1: Captura de Áudio

#### Prompt 1.1 — AudioCapture

````markdown
# TAREFA: Implementar captura de áudio

Criar `src/core/AudioCapture.ts`:

## Requisitos:

- Classe AudioCapture com métodos: start(), stop(), getStream()
- Usar navigator.mediaDevices.getUserMedia
- Configurar AudioContext com sampleRate 16000 (ou resample)
- Emitir eventos: 'start', 'stop', 'error'
- Tratar permissão negada graciosamente

## Interface esperada:

```typescript
class AudioCapture extends EventEmitter {
    async start(): Promise<void>
    stop(): void
    isRecording(): boolean
}
```
````

## Logs obrigatórios:

console.log('[AudioCapture]', { event: 'start' | 'stop' | 'error', details })

````

#### Teste ao Vivo 1.1

```typescript
// scripts/live-test-audio.ts
import { AudioCapture } from '../src/core/AudioCapture';

async function testAudioCapture() {
  const results = {
    phase: 1,
    module: 'AudioCapture',
    tests: [] as any[]
  };

  const capture = new AudioCapture();

  // Teste 1: Inicialização
  try {
    await capture.start();
    results.tests.push({ name: 'start()', pass: true });
  } catch (e) {
    results.tests.push({ name: 'start()', pass: false, error: e.message });
  }

  // Teste 2: Estado
  const isRecording = capture.isRecording();
  results.tests.push({
    name: 'isRecording()',
    pass: isRecording === true,
    actual: isRecording
  });

  // Teste 3: Stop
  capture.stop();
  results.tests.push({
    name: 'stop()',
    pass: capture.isRecording() === false
  });

  console.log(JSON.stringify(results, null, 2));
  return results;
}

testAudioCapture();
````

#### Prompt 1.2 — VAD Integration

````markdown
# TAREFA: Integrar Silero VAD

Criar `src/core/VADProcessor.ts`:

## Requisitos:

- Importar @ricky0123/vad-web
- Configurar MicVAD com thresholds:
    - positiveSpeechThreshold: 0.5
    - negativeSpeechThreshold: 0.35
    - redemptionFrames: 8 (pausa ~240ms)
- Emitir evento 'speechEnd' com Float32Array do áudio
- Emitir evento 'speechStart' para feedback visual

## Interface:

```typescript
class VADProcessor extends EventEmitter {
    async initialize(): Promise<void>
    start(): void
    stop(): void
    // Eventos: 'speechStart', 'speechEnd', 'error'
}
```
````

## Callback onSpeechEnd deve:

1. Logar duração do chunk: console.log('[VAD]', { event: 'speechEnd', durationMs })
2. Emitir evento com audio data

````

#### Teste ao Vivo 1.2

```typescript
// scripts/live-test-vad.ts
// REQUER: Microfone ativo e fala real

async function testVAD() {
  const results = { phase: 1, module: 'VADProcessor', tests: [] };
  const vad = new VADProcessor();

  let speechDetected = false;
  let audioReceived = false;

  vad.on('speechStart', () => { speechDetected = true; });
  vad.on('speechEnd', (audio) => {
    audioReceived = audio instanceof Float32Array && audio.length > 0;
  });

  await vad.initialize();
  vad.start();

  console.log('[TEST] Fale algo por 2 segundos...');

  await new Promise(r => setTimeout(r, 5000)); // Aguarda 5s

  vad.stop();

  results.tests.push({ name: 'speechStart detected', pass: speechDetected });
  results.tests.push({ name: 'speechEnd with audio', pass: audioReceived });

  console.log(JSON.stringify(results, null, 2));
}
````

---

### FASE 2: Transcrição

#### Prompt 2.1 — Transcriber

````markdown
# TAREFA: Implementar transcrição com Whisper

Criar `src/core/Transcriber.ts`:

## Requisitos:

- Usar @huggingface/transformers pipeline
- Modelo: Xenova/whisper-tiny (ou configurável)
- Carregar modelo lazy (primeira transcrição)
- Progress callback durante carregamento

## Interface:

```typescript
interface TranscribeResult {
    text: string
    language: string
    duration: number // tempo de processamento em ms
}

class Transcriber {
    async initialize(model?: string): Promise<void>
    async transcribe(
        audio: Float32Array,
        lang?: string
    ): Promise<TranscribeResult>
    isReady(): boolean
}
```
````

## Logs obrigatórios:

- '[Transcriber] Loading model...' com progresso
- '[Transcriber] Transcription complete' com { text, durationMs, audioLengthMs }

````

#### Teste ao Vivo 2.1

```typescript
// scripts/live-test-transcriber.ts
import { Transcriber } from '../src/core/Transcriber';
import * as fs from 'fs';

async function testTranscriber() {
  const results = { phase: 2, module: 'Transcriber', tests: [] };

  const transcriber = new Transcriber();

  // Teste 1: Inicialização
  console.log('[TEST] Carregando modelo (pode demorar ~30s)...');
  const startLoad = Date.now();
  await transcriber.initialize('Xenova/whisper-tiny');
  const loadTime = Date.now() - startLoad;

  results.tests.push({
    name: 'Model loaded',
    pass: transcriber.isReady(),
    loadTimeMs: loadTime
  });

  // Teste 2: Transcrição de áudio sintético (silêncio)
  const silentAudio = new Float32Array(16000); // 1s silêncio
  const silentResult = await transcriber.transcribe(silentAudio, 'pt');

  results.tests.push({
    name: 'Transcribe silent audio',
    pass: silentResult.text !== undefined,
    result: silentResult
  });

  // Teste 3: Transcrição de fixture (se existir)
  if (fs.existsSync('./tests/fixtures/audio/test-pt-3s.wav')) {
    // Carregar e transcrever fixture
    // ... código para carregar WAV
    results.tests.push({ name: 'Transcribe fixture', pass: true, text: '...' });
  }

  console.log(JSON.stringify(results, null, 2));
}
````

---

### FASE 3: Pós-Processamento

#### Prompt 3.1 — TextProcessor

````markdown
# TAREFA: Implementar processamento de texto

Criar `src/core/TextProcessor.ts`:

## Requisitos:

### Nível 1 (Regex - sempre ativo):

- Remover disfluências: éh, hum, ã, tipo, né, então (início)
- Remover repetições: "eu eu" → "eu", "que que" → "que"
- Normalizar espaços múltiplos
- Capitalizar primeira letra

### Nível 2 (Pontuação - opcional):

- Integrar com Ollama se disponível
- Fallback: heurística básica (ponto após pausas longas)

## Interface:

```typescript
interface ProcessOptions {
    removeDisfluencies: boolean
    addPunctuation: boolean
    formatParagraphs: boolean
}

class TextProcessor {
    process(text: string, options?: ProcessOptions): string
    async processWithLLM(text: string): Promise<string> // opcional
}
```
````

## Casos de teste inline:

// Input: "éh então eu eu acho que que sim" // Output: "Eu acho que sim"

// Input: "hum deixa eu pensar né tipo assim" // Output: "Deixa eu pensar"

````

#### Teste ao Vivo 3.1

```typescript
// scripts/live-test-textprocessor.ts
import { TextProcessor } from '../src/core/TextProcessor';

function testTextProcessor() {
  const results = { phase: 3, module: 'TextProcessor', tests: [] };
  const processor = new TextProcessor();

  const testCases = [
    {
      input: 'éh então eu eu acho que sim',
      expected: 'Eu acho que sim',
      name: 'Remove disfluencies + repetitions'
    },
    {
      input: 'hum deixa eu pensar né tipo assim',
      expected: 'Deixa eu pensar',
      name: 'Remove fillers'
    },
    {
      input: 'muito muito obrigado',
      expected: 'Muito muito obrigado', // preserva ênfase
      name: 'Preserve intentional repetition'
    },
    {
      input: '  texto   com   espaços  ',
      expected: 'Texto com espaços',
      name: 'Normalize spaces'
    }
  ];

  for (const tc of testCases) {
    const actual = processor.process(tc.input);
    const pass = actual.trim() === tc.expected;
    results.tests.push({
      name: tc.name,
      pass,
      input: tc.input,
      expected: tc.expected,
      actual
    });
  }

  console.log(JSON.stringify(results, null, 2));
}
````

---

### FASE 4: Integração com Editor

#### Prompt 4.1 — EditorIntegration

````markdown
# TAREFA: Integrar com editor Obsidian

Criar `src/core/EditorIntegration.ts`:

## Requisitos:

- Inserir texto na posição do cursor
- Suportar modos: insert, append, replace
- Criar transação única (undo funciona)
- Não travar se nenhum editor ativo

## Interface:

```typescript
type InsertMode = 'cursor' | 'append' | 'replace'

class EditorIntegration {
    constructor(app: App)
    insert(text: string, mode: InsertMode): boolean
    getActiveEditor(): Editor | null
}
```
````

## Logs:

console.log('[Editor]', { action: 'insert', mode, charCount, success })

````

---

### FASE 5: Pipeline Completo

#### Prompt 5.1 — Pipeline

```markdown
# TAREFA: Criar pipeline completo

Criar `src/core/Pipeline.ts`:

## Requisitos:
- Orquestrar: VAD → Transcriber → TextProcessor → Editor
- Estados: idle, listening, transcribing, processing
- Eventos para UI: stateChange, progress, complete, error

## Interface:
```typescript
type PipelineState = 'idle' | 'listening' | 'transcribing' | 'processing';

class Pipeline extends EventEmitter {
  async initialize(): Promise<void>
  start(): void
  stop(): void
  getState(): PipelineState
  // Eventos: 'stateChange', 'transcription', 'complete', 'error'
}
````

## Fluxo:

1. start() → estado 'listening' → VAD ativo
2. VAD detecta fala → estado 'transcribing'
3. Whisper transcreve → estado 'processing'
4. TextProcessor limpa → insere no editor → estado 'listening'

````

#### Teste ao Vivo 5.1 — End-to-End

```typescript
// scripts/live-test-e2e.ts
// REQUER: Obsidian aberto com vault de teste

async function testE2E() {
  const results = { phase: 5, module: 'Pipeline E2E', tests: [] };

  console.log('[E2E] Iniciando teste completo...');
  console.log('[E2E] INSTRUÇÃO: Fale "Olá mundo" quando solicitado');

  const pipeline = new Pipeline();

  let statesObserved: string[] = [];
  let finalText = '';

  pipeline.on('stateChange', (state) => {
    statesObserved.push(state);
    console.log(`[E2E] Estado: ${state}`);
  });

  pipeline.on('complete', (text) => {
    finalText = text;
    console.log(`[E2E] Texto final: "${text}"`);
  });

  await pipeline.initialize();

  console.log('[E2E] >>> FALE AGORA: "Olá mundo" <<<');
  pipeline.start();

  // Aguarda até 15s por resultado
  await new Promise(r => setTimeout(r, 15000));

  pipeline.stop();

  // Validações
  results.tests.push({
    name: 'States observed',
    pass: statesObserved.includes('listening') && statesObserved.includes('transcribing'),
    observed: statesObserved
  });

  results.tests.push({
    name: 'Text received',
    pass: finalText.length > 0,
    text: finalText
  });

  results.tests.push({
    name: 'Text contains expected words',
    pass: finalText.toLowerCase().includes('olá') || finalText.toLowerCase().includes('mundo'),
    text: finalText
  });

  console.log(JSON.stringify(results, null, 2));
}
````

---

## 3. Scripts de Validação

### validate-phase.ts

```typescript
#!/usr/bin/env tsx
// scripts/validate-phase.ts
// Uso: npx tsx scripts/validate-phase.ts 0|1|2|3|4|5

import { execSync } from 'child_process'
import * as fs from 'fs'

const phase = parseInt(process.argv[2] || '0')

interface ValidationResult {
    phase: number
    timestamp: string
    status: 'pass' | 'fail'
    checks: Array<{
        name: string
        pass: boolean
        details?: any
    }>
    nextPhase?: number
    blockers?: string[]
}

const validators: Record<number, () => ValidationResult> = {
    0: validatePhase0,
    1: validatePhase1,
    2: validatePhase2,
    3: validatePhase3,
    4: validatePhase4,
    5: validatePhase5,
}

function validatePhase0(): ValidationResult {
    const checks = []

    // Check 1: package.json existe
    checks.push({
        name: 'package.json exists',
        pass: fs.existsSync('./package.json'),
    })

    // Check 2: Build funciona
    try {
        execSync('npm run build', { stdio: 'pipe' })
        checks.push({ name: 'npm run build', pass: true })
    } catch (e) {
        checks.push({ name: 'npm run build', pass: false, details: e.message })
    }

    // Check 3: main.js gerado
    checks.push({
        name: 'main.js generated',
        pass: fs.existsSync('./main.js'),
    })

    // Check 4: manifest.json válido
    try {
        const manifest = JSON.parse(fs.readFileSync('./manifest.json', 'utf-8'))
        checks.push({
            name: 'manifest.json valid',
            pass: manifest.id === 'dictate2me' && !!manifest.version,
        })
    } catch {
        checks.push({ name: 'manifest.json valid', pass: false })
    }

    const allPass = checks.every((c) => c.pass)

    return {
        phase: 0,
        timestamp: new Date().toISOString(),
        status: allPass ? 'pass' : 'fail',
        checks,
        nextPhase: allPass ? 1 : undefined,
        blockers: checks.filter((c) => !c.pass).map((c) => c.name),
    }
}

function validatePhase1(): ValidationResult {
    const checks = []

    // Check: AudioCapture existe
    checks.push({
        name: 'AudioCapture.ts exists',
        pass: fs.existsSync('./src/core/AudioCapture.ts'),
    })

    // Check: VADProcessor existe
    checks.push({
        name: 'VADProcessor.ts exists',
        pass: fs.existsSync('./src/core/VADProcessor.ts'),
    })

    // Check: TypeScript compila
    try {
        execSync('npx tsc --noEmit', { stdio: 'pipe' })
        checks.push({ name: 'TypeScript compiles', pass: true })
    } catch (e) {
        checks.push({ name: 'TypeScript compiles', pass: false })
    }

    // Check: Testes unitários passam
    try {
        execSync('npm run test:unit -- --run', { stdio: 'pipe' })
        checks.push({ name: 'Unit tests pass', pass: true })
    } catch {
        checks.push({ name: 'Unit tests pass', pass: false })
    }

    const allPass = checks.every((c) => c.pass)

    return {
        phase: 1,
        timestamp: new Date().toISOString(),
        status: allPass ? 'pass' : 'fail',
        checks,
        nextPhase: allPass ? 2 : undefined,
        blockers: checks.filter((c) => !c.pass).map((c) => c.name),
    }
}

// ... validatePhase2, 3, 4, 5 seguem padrão similar

function validatePhase2(): ValidationResult {
    const checks = []

    checks.push({
        name: 'Transcriber.ts exists',
        pass: fs.existsSync('./src/core/Transcriber.ts'),
    })

    // Verificar se modelo pode ser carregado (teste mais demorado)
    // Pode ser pulado com flag --skip-model

    const allPass = checks.every((c) => c.pass)
    return {
        phase: 2,
        timestamp: new Date().toISOString(),
        status: allPass ? 'pass' : 'fail',
        checks,
        nextPhase: allPass ? 3 : undefined,
        blockers: checks.filter((c) => !c.pass).map((c) => c.name),
    }
}

function validatePhase3(): ValidationResult {
    const checks = []

    checks.push({
        name: 'TextProcessor.ts exists',
        pass: fs.existsSync('./src/core/TextProcessor.ts'),
    })

    // Executar testes de processamento
    try {
        const output = execSync('npx tsx scripts/live-test-textprocessor.ts', {
            encoding: 'utf-8',
        })
        const result = JSON.parse(output)
        const allTestsPass = result.tests.every((t: any) => t.pass)
        checks.push({
            name: 'TextProcessor tests',
            pass: allTestsPass,
            details: result.tests,
        })
    } catch (e) {
        checks.push({ name: 'TextProcessor tests', pass: false })
    }

    const allPass = checks.every((c) => c.pass)
    return {
        phase: 3,
        timestamp: new Date().toISOString(),
        status: allPass ? 'pass' : 'fail',
        checks,
        nextPhase: allPass ? 4 : undefined,
        blockers: checks.filter((c) => !c.pass).map((c) => c.name),
    }
}

function validatePhase4(): ValidationResult {
    const checks = []

    checks.push({
        name: 'EditorIntegration.ts exists',
        pass: fs.existsSync('./src/core/EditorIntegration.ts'),
    })

    checks.push({
        name: 'Pipeline.ts exists',
        pass: fs.existsSync('./src/core/Pipeline.ts'),
    })

    const allPass = checks.every((c) => c.pass)
    return {
        phase: 4,
        timestamp: new Date().toISOString(),
        status: allPass ? 'pass' : 'fail',
        checks,
        nextPhase: allPass ? 5 : undefined,
        blockers: checks.filter((c) => !c.pass).map((c) => c.name),
    }
}

function validatePhase5(): ValidationResult {
    const checks = []

    // Validações finais
    try {
        execSync('npm run build', { stdio: 'pipe' })
        const stats = fs.statSync('./main.js')
        checks.push({
            name: 'Bundle size < 5MB',
            pass: stats.size < 5 * 1024 * 1024,
            details: { sizeBytes: stats.size },
        })
    } catch {
        checks.push({ name: 'Final build', pass: false })
    }

    // Todos os testes
    try {
        execSync('npm test -- --run', { stdio: 'pipe' })
        checks.push({ name: 'All tests pass', pass: true })
    } catch {
        checks.push({ name: 'All tests pass', pass: false })
    }

    const allPass = checks.every((c) => c.pass)
    return {
        phase: 5,
        timestamp: new Date().toISOString(),
        status: allPass ? 'pass' : 'fail',
        checks,
        nextPhase: allPass ? undefined : undefined, // Fase final
        blockers: checks.filter((c) => !c.pass).map((c) => c.name),
    }
}

// Executar
const validator = validators[phase]
if (!validator) {
    console.error(`Fase ${phase} não existe. Use 0-5.`)
    process.exit(1)
}

const result = validator()
console.log(JSON.stringify(result, null, 2))

// Exit code baseado no resultado
process.exit(result.status === 'pass' ? 0 : 1)
```

---

## 4. Workflow VS Code para IA

### 4.1 Sequência de Comandos

Cada fase segue este workflow:

```
1. CARREGAR CONTEXTO
   → Copiar contexto base + prompt da fase para IA

2. GERAR CÓDIGO
   → IA gera arquivo(s) solicitado(s)
   → Salvar em src/

3. EXECUTAR BUILD
   → Terminal: npm run build
   → Copiar output para IA se houver erro

4. EXECUTAR TESTE AO VIVO
   → Terminal: npx tsx scripts/live-test-{modulo}.ts
   → Copiar JSON de resultado para IA

5. VALIDAR FASE
   → Terminal: npx tsx scripts/validate-phase.ts {N}
   → Se pass: próxima fase
   → Se fail: IA corrige baseado nos blockers

6. COMMIT
   → git add . && git commit -m "Phase {N}: {descrição}"
```

### 4.2 Comandos VS Code (tasks.json)

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Build Plugin",
            "type": "shell",
            "command": "npm run build",
            "group": "build",
            "problemMatcher": ["$tsc"]
        },
        {
            "label": "Validate Phase 0",
            "type": "shell",
            "command": "npx tsx scripts/validate-phase.ts 0"
        },
        {
            "label": "Validate Phase 1",
            "type": "shell",
            "command": "npx tsx scripts/validate-phase.ts 1"
        },
        {
            "label": "Live Test Audio",
            "type": "shell",
            "command": "npx tsx scripts/live-test-audio.ts"
        },
        {
            "label": "Live Test VAD",
            "type": "shell",
            "command": "npx tsx scripts/live-test-vad.ts"
        },
        {
            "label": "Live Test Transcriber",
            "type": "shell",
            "command": "npx tsx scripts/live-test-transcriber.ts"
        },
        {
            "label": "Live Test E2E",
            "type": "shell",
            "command": "npx tsx scripts/live-test-e2e.ts"
        },
        {
            "label": "Run All Tests",
            "type": "shell",
            "command": "npm test -- --run"
        }
    ]
}
```

### 4.3 Template de Prompt para Correção

Quando um teste falha, usar este template:

````markdown
# CORREÇÃO NECESSÁRIA

## Contexto

Fase: {N}
Módulo: {nome}

## Resultado do Teste

```json
{colar JSON do resultado}
```
````

## Erro Específico

{descrição do erro ou blocker}

## Código Atual

```typescript
{código do arquivo com problema}
```

## Tarefa

Corrija o código para que o teste passe. Mantenha a interface existente. Explique brevemente a causa do erro antes de fornecer o código corrigido.

````

---

## 5. Testes Unitários (Vitest)

### vitest.config.ts

```typescript
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    include: ['tests/unit/**/*.test.ts'],
    testTimeout: 30000,
    reporters: ['json', 'default'],
    outputFile: './test-results.json'
  }
});
````

### Exemplo: TextProcessor.test.ts

```typescript
import { describe, it, expect } from 'vitest'
import { TextProcessor } from '../../src/core/TextProcessor'

describe('TextProcessor', () => {
    const processor = new TextProcessor()

    describe('removeDisfluencies', () => {
        const cases = [
            ['éh então sim', 'Então sim'],
            ['hum deixa ver', 'Deixa ver'],
            ['eu eu acho', 'Eu acho'],
            ['que que isso', 'Que isso'],
            ['muito muito bom', 'Muito muito bom'], // preserva ênfase
        ]

        it.each(cases)('"%s" → "%s"', (input, expected) => {
            expect(processor.process(input)).toBe(expected)
        })
    })
})
```

---

## 6. Checklist de Progresso

```markdown
## dictate2me - Progresso de Desenvolvimento

### Fase 0: Setup ⬜

- [ ] package.json criado
- [ ] tsconfig.json configurado
- [ ] esbuild.config.mjs funcionando
- [ ] manifest.json válido
- [ ] main.ts básico
- [ ] `npm run build` sem erros
- [ ] Validação: `npx tsx scripts/validate-phase.ts 0` ✓

### Fase 1: Áudio ⬜

- [ ] AudioCapture.ts implementado
- [ ] VADProcessor.ts implementado
- [ ] Teste ao vivo: captura funciona
- [ ] Teste ao vivo: VAD detecta fala
- [ ] Validação: `npx tsx scripts/validate-phase.ts 1` ✓

### Fase 2: Transcrição ⬜

- [ ] Transcriber.ts implementado
- [ ] Modelo carrega corretamente
- [ ] Teste ao vivo: transcrição funciona
- [ ] Benchmark: latência aceitável
- [ ] Validação: `npx tsx scripts/validate-phase.ts 2` ✓

### Fase 3: Processamento ⬜

- [ ] TextProcessor.ts implementado
- [ ] Testes unitários passam
- [ ] Teste ao vivo: limpeza funciona
- [ ] Validação: `npx tsx scripts/validate-phase.ts 3` ✓

### Fase 4: Integração ⬜

- [ ] EditorIntegration.ts implementado
- [ ] Pipeline.ts implementado
- [ ] Teste E2E básico funciona
- [ ] Validação: `npx tsx scripts/validate-phase.ts 4` ✓

### Fase 5: Polimento ⬜

- [ ] Settings completo
- [ ] UI feedback
- [ ] Documentação
- [ ] Bundle otimizado
- [ ] Validação: `npx tsx scripts/validate-phase.ts 5` ✓
```

---

## 7. Resumo: Comandos por Fase

| Fase | Prompt   | Teste ao Vivo                            | Validação             |
| ---- | -------- | ---------------------------------------- | --------------------- |
| 0    | 0.1, 0.2 | `npm run build`                          | `validate-phase.ts 0` |
| 1    | 1.1, 1.2 | `live-test-audio.ts`, `live-test-vad.ts` | `validate-phase.ts 1` |
| 2    | 2.1      | `live-test-transcriber.ts`               | `validate-phase.ts 2` |
| 3    | 3.1      | `live-test-textprocessor.ts`             | `validate-phase.ts 3` |
| 4    | 4.1, 5.1 | `live-test-e2e.ts`                       | `validate-phase.ts 4` |
| 5    | —        | Todos                                    | `validate-phase.ts 5` |

---

_Documento versão 2.0 — Dezembro 2024_ _Metodologia: AI-Driven Development com Validação ao Vivo_
