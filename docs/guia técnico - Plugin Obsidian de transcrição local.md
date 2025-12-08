A melhor arquitetura para um plugin Obsidian de transcrição local combina **Transformers.js** para STT (sem módulos nativos), **@ricky0123/vad-web** para detecção de voz, e uma pipeline híbrida de pós-processamento com regex + modelo de pontuação + LLM opcional via Ollama. Esta abordagem funciona nativamente em macOS/Windows/Linux sem dependências externas, consome ~**400-500MB** de RAM com o modelo Whisper tiny, e processa áudio em **2-3x tempo real** no M1.

## Arquitetura de plugins Obsidian para processamento pesado

Plugins Obsidian são módulos **TypeScript/JavaScript** que executam no processo renderer do Electron. O ciclo de vida básico inclui `onload()` para inicialização e `onunload()` para limpeza, com métodos como `registerEvent()`, `registerDomEvent()` e `registerInterval()` garantindo cleanup automático.

A API do Obsidian expõe três superfícies principais: **app.vault** (sistema de arquivos), **app.workspace** (UI/editor) e **app.metadataCache** (indexação). Plugins podem ler/escrever arquivos, criar comandos, adicionar itens à status bar e responder a eventos do vault.

Para processamento ML, existem **três restrições críticas** do ambiente Electron:

- Todo código deve compilar em um único `main.js` (bundle requirement)
- Web Workers **não podem usar módulos nativos** Node.js
- Processamento pesado no thread principal bloqueia a UI

A solução recomendada para transcriçãoé **WASM-based** usando Transformers.js, que não requer módulos nativos. A alternativa é um **servidor local externo** (whisper.cpp rodando separadamente) comunicando via HTTP. Plugins como Datacore demonstram o padrão: usar `esbuild-plugin-inline-worker` para bundle de workers, rate-limiting para evitar lag, e processamento incremental em background.

```typescript
// Estrutura padrão de plugin
import { Plugin } from 'obsidian'

export default class TranscriptionPlugin extends Plugin {
    async onload() {
        this.addCommand({
            id: 'start-transcription',
            name: 'Iniciar Transcrição',
            callback: () => this.startTranscription(),
        })
    }
}
```

O **hot reload** durante desenvolvimento usa o plugin `pjeby/hot-reload`: basta criar um arquivo `.hotreload` no diretório do plugin. Para testes, `jest-environment-obsidian` fornece mocks da API, mas a estratégia mais robusta é dependency inversion—separar lógica de negócio do código específico do Obsidian.

## Captura de áudio e detecção de pausas com Silero VAD

Para captura de áudio, a **Web Audio API** é a escolha mais simples no contexto Electron, funcionando em todas as plataformas via `navigator.mediaDevices.getUserMedia()`. No macOS, requer `systemPreferences.askForMediaAccess('microphone')` para permissões.

| Biblioteca          | Plataformas | Prós                           | Contras                                |
| ------------------- | ----------- | ------------------------------ | -------------------------------------- |
| **Web Audio API**   | Todas       | Nativo no Electron, sem deps   | Sample rate 44.1kHz (precisa resample) |
| **naudiodon**       | Todas       | PortAudio nativo, Node streams | Requer node-gyp build                  |
| **speech-recorder** | Todas       | VAD integrado (WebRTC+Silero)  | Precisa electron-rebuild               |

A recomendação principal é **@ricky0123/vad-web** (Silero VAD), que roda inteiramente no renderer process com Web Audio API. Silero VAD é state-of-the-art em detecção de fala, usando uma rede neural ONNX de ~**2MB** que consome ~**50-100MB** de RAM com o runtime ONNX. O output é Float32Array a **16kHz**—exatamente o formato que Whisper espera.

```javascript
import { MicVAD } from '@ricky0123/vad-web'

const vad = await MicVAD.new({
    positiveSpeechThreshold: 0.5, // threshold de detecção
    negativeSpeechThreshold: 0.35, // threshold de silêncio
    redemptionFrames: 8, // ~240ms grace period
    preSpeechPadFrames: 3, // captura lead-in
    minSpeechFrames: 3, // duração mínima

    onSpeechEnd: async (audio) => {
        // audio já está em Float32 @ 16kHz - enviar para Whisper
        const transcript = await transcribeWithWhisper(audio)
        insertIntoEditor(transcript)
    },
})

await vad.start()
```

A configuração `redemptionFrames: 8` (~240ms) define a duração mínima de silêncio para considerar fim de frase. Para pausas mais longas entre parágrafos, ajuste para 15-20 frames (~500-600ms).

Comparativo de VAD: **Silero VAD** oferece a melhor precisão (estado da arte) mas usa mais memória (~80MB). **WebRTC VAD** é muito mais leve (<10MB) e rápido, mas tem mais falsos positivos. **RNNoise** é bom se você também precisa de supressão de ruído.

## Opções de transcrição local: Transformers.js vs whisper.cpp

### Comparativo de implementações Whisper

| Biblioteca             | Performance (M1)  | Memória (tiny) | Integração           | Licença    |
| ---------------------- | ----------------- | -------------- | -------------------- | ---------- |
| **Transformers.js**    | Boa (WASM/ONNX)   | ~300MB         | ⭐⭐⭐⭐⭐ Excelente | Apache 2.0 |
| **whisper.cpp (Node)** | ⭐⭐⭐⭐⭐ Melhor | ~273MB         | ⭐⭐⭐ Moderada      | MIT        |
| **Vosk**               | Boa               | ~50MB          | ⭐⭐⭐ Moderada      | Apache 2.0 |

Para um plugin Obsidian, **Transformers.js** (`@huggingface/transformers`) é a escolha primária porque não requer módulos nativos—funciona identicamente em todas as plataformas sem builds específicos. A integração é trivial:

```typescript
import { pipeline } from '@huggingface/transformers'

class WhisperTranscriber {
    private transcriber: any

    async initialize() {
        this.transcriber = await pipeline(
            'automatic-speech-recognition',
            'Xenova/whisper-tiny',
            { quantized: true } // modelo quantizado para menor memória
        )
    }

    async transcribe(audioFloat32: Float32Array): Promise<string> {
        const result = await this.transcriber(audioFloat32, {
            language: 'pt', // Português
            return_timestamps: true,
        })
        return result.text
    }
}
```

O modelo `Xenova/whisper-tiny` tem **39M parâmetros**, ocupa ~75MB em disco e usa ~**273-300MB** de RAM durante inferência. No M1, processa áudio em ~**2-3x tempo real** (1 minuto de áudio = 20-30 segundos de processamento).

### Tamanhos e requisitos dos modelos Whisper

| Modelo     | Parâmetros | Tamanho | RAM     | Velocidade | Recomendado para               |
| ---------- | ---------- | ------- | ------- | ---------- | ------------------------------ |
| **tiny**   | 39M        | 75 MB   | ~273 MB | ~32x       | MVP, recursos limitados        |
| **base**   | 74M        | 142 MB  | ~388 MB | ~16x       | Equilíbrio velocidade/precisão |
| **small**  | 244M       | 466 MB  | ~852 MB | ~6x        | Boa precisão diária            |
| **medium** | 769M       | 1.5 GB  | ~2.1 GB | ~2x        | Alta precisão                  |
| **turbo**  | 809M       | 1.5 GB  | ~2.5 GB | ~8x        | Large otimizado                |

Para **MacBook Air M1 16GB**, comece com **tiny** e upgrade para **base** ou **small** se a precisão for insuficiente. Modelos até **medium** rodam confortavelmente em 16GB.

Se performance máxima for crítica, a alternativa é **nodejs-whisper** ou **smart-whisper** (bindings para whisper.cpp). No M1, whisper.cpp com ARM NEON atinge ~**10x tempo real** com o modelo tiny. Porém, requer `electron-rebuild` e compilação de módulos nativos—adicionando complexidade de distribuição.

## Pós-processamento: pipeline híbrida recomendada

O output bruto do Whisper frequentemente precisa de correções: pontuação inconsistente, disfluências ("éh", "hum"), e formatação para notas. A abordagem mais eficiente é uma **pipeline híbrida**:

```
Whisper Output → Regex (disfluências) → Modelo ONNX (pontuação) → LLM opcional (parágrafos)
```

### Etapa 1: Remoção de disfluências com regex (instantâneo)

```javascript
const disfluencias =
    /\b(é|éh|ã|hum|uhm|uh|tipo|assim|né|então|basicamente|na verdade)\b,?\s*/gi
const limpo = texto.replace(disfluencias, '')
```

Esta etapa é **instantânea** (<1ms) e remove os fillers mais comuns em português e inglês.

### Etapa 2: Restauração de pontuação com modelo dedicado

Bibliotecas rule-based como **compromise.js** e **natural.js** **não conseguem adicionar pontuação** em texto sem pontuação—apenas funcionam quando a pontuação já existe. Para restauração, use **deepmultilingualpunctuation** (modelo `oliverguhr/fullstop-punctuation-multilang-large`):

- **F1 scores**: Período 94.8%, Vírgula 81.9%, Interrogação 89.0%
- **Tamanho**: ~1.2GB (0.6B parâmetros)
- **Licença**: MIT
- **Velocidade**: ~500-1000 palavras/segundo no M1

O modelo está disponível em formato ONNX e pode ser executado via Transformers.js:

```typescript
// Usando @xenova/transformers para punctuation
const punctuator = await pipeline(
    'token-classification',
    'oliverguhr/fullstop-punctuation-multilang-large'
)
const result = await punctuator(textoSemPontuacao)
```

### Etapa 3: Detecção de parágrafos e formatação opcional (LLM)

Para inserir quebras de parágrafo e formatação avançada, três opções:

1. **Heurística baseada em timestamps**: Whisper fornece timestamps por palavra; quebre em pausas >1.5s
2. **Contagem de sentenças**: Inserir quebra a cada 4-5 sentenças
3. **LLM via Ollama**: Mais preciso, detecta mudanças de tópico

Para LLM local, **Ollama** é a solução mais simples (API REST em `localhost:11434`). Modelos recomendados para M1 16GB:

| Modelo            | RAM    | Velocidade  | Qualidade |
| ----------------- | ------ | ----------- | --------- |
| **Qwen 2.5 1.5B** | ~1.5GB | 60-80 tok/s | Boa       |
| **Llama 3.2 1B**  | ~1.8GB | 150+ tok/s  | Boa       |
| **Llama 3.2 3B**  | ~3.4GB | 50-70 tok/s | Muito boa |
| **Phi-3 Mini**    | ~6-8GB | 30-50 tok/s | Excelente |

**Llama 3.2 3B** oferece o melhor equilíbrio para o M1 16GB. Prompt sugerido:

```
Limpe esta transcrição: adicione pontuação, remova repetições,
insira quebras de parágrafo em mudanças de tópico.
Mantenha todo o conteúdo significativo intacto.

Transcrição: {texto}
```

Para evitar overhead do LLM em transcrições curtas, use apenas para textos >200 palavras ou quando o usuário solicitar "formatação avançada".

## Projetos similares e componentes reutilizáveis

**Não existe nenhum plugin Obsidian verdadeiramente local**—todos os existentes usam APIs cloud ou requerem servidor Docker separado. Os principais projetos de referência:

- **whisper-obsidian-plugin** (nikdanilov): Usa OpenAI API, mas demonstra boa arquitetura de plugin
- **obsidian-transcription** (djmango): Transcreve arquivos de áudio via servidor Whisper local (Docker)
- **Buzz** (15.6k stars): Melhor referência para gerenciamento de modelos e multi-backend
- **whisper-writer** (GPL-3.0): Melhor referência para workflow de ditado com VAD

Para integração Node.js/Electron, **@kutalia/whisper-node-addon** oferece binários pré-compilados para Windows, Linux e macOS (incluindo arm64). Alternativa mais simples é usar Transformers.js que não requer binários nativos.

A pipeline de áudio padrão observada nos projetos maduros:

```
Captura → Resample 16kHz → Ring Buffer → VAD → Chunk → STT → Pós-proc → Output
```

Todos os projetos bem-sucedidos implementam múltiplos modos de gravação: push-to-talk, toggle, VAD automático, e contínuo.

## Arquitetura recomendada para o plugin

```
┌─────────────────────────────────────────────────────────────┐
│                    OBSIDIAN PLUGIN                          │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │ Web Audio API│───▶│@ricky0123/   │───▶│ Float32Array │  │
│  │  (Microfone) │    │  vad-web     │    │   @ 16kHz    │  │
│  └──────────────┘    │ (Silero VAD) │    └──────────────┘  │
│                      └──────────────┘           │          │
│                                                 ▼          │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │ Editor Insert│◀───│Pós-processa- │◀───│Transformers.js│ │
│  │              │    │  mento       │    │(Whisper tiny)│  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│                             │                              │
│                    ┌────────┴────────┐                     │
│                    │ Opcional: Ollama │                    │
│                    │  (localhost API) │                    │
│                    └─────────────────┘                     │
└─────────────────────────────────────────────────────────────┘
```

### Consumo estimado de recursos (M1 16GB)

| Componente                | CPU    | RAM    | Latência       |
| ------------------------- | ------ | ------ | -------------- |
| Web Audio capture         | <1%    | ~5MB   | <10ms          |
| @ricky0123/vad-web        | 2-5%   | ~80MB  | 10-20ms/frame  |
| Transformers.js (tiny)    | 30-50% | ~300MB | 2-3x real-time |
| Pós-processamento (regex) | <1%    | <1MB   | <1ms           |
| Ollama (Llama 3.2 3B)     | 50-80% | ~3.4GB | ~50 tok/s      |

**Total sem LLM**: ~400MB RAM, ~35% CPU durante transcrição **Total com LLM**: ~4GB RAM pico, ~80% CPU durante formatação

### Tooling de desenvolvimento

```json
// package.json
{
    "scripts": {
        "dev": "node esbuild.config.mjs",
        "build": "tsc -noEmit && node esbuild.config.mjs production",
        "test": "vitest"
    },
    "devDependencies": {
        "@types/node": "^20.10.0",
        "esbuild": "^0.19.9",
        "obsidian": "latest",
        "typescript": "^5.3.3",
        "vitest": "^1.0.0"
    },
    "dependencies": {
        "@huggingface/transformers": "^3.0.0",
        "@ricky0123/vad-web": "^0.0.29"
    }
}
```

Use **Vitest** ou **Bun test** para testes unitários (~200ms vs 6-8s do Jest). Para hot reload, instale `pjeby/hot-reload` no vault e crie `.hotreload` no diretório do plugin.

## Conclusão

O stack técnico recomendado para um plugin de transcrição 100% local em Obsidian é:

1. **Audio**: Web Audio API (nativo) + **@ricky0123/vad-web** (Silero VAD)
2. **STT**: **Transformers.js** com `Xenova/whisper-tiny` (upgrade para base/small se necessário)
3. **Pós-proc**: Regex para disfluências + **fullstop-punctuation** (ONNX) + Ollama opcional
4. **Dev**: esbuild + TypeScript + hot-reload + Vitest

Esta arquitetura evita módulos nativos (simplificando distribuição cross-platform), consome ~400MB de RAM em operação normal, e processa áudio em 2-3x tempo real no M1. A principal lacuna nos projetos existentes—um plugin Obsidian verdadeiramente local com ditado em tempo real—representa uma oportunidade clara de diferenciação.
