# AI Models Directory

This directory stores the AI models used by dictate2me for transcription and text correction.

## Models Required

### Whisper (Transcription)

- **Model**: `whisper-small` or `whisper-medium`
- **Format**: GGML/GGUF
- **Quantization**: Q5_K_M (recommended) or Q4_K_M (faster, less accurate)
- **Download**: Run `./scripts/download-models.sh` or `make models`

### LLM (Text Correction)

- **Model**: `Phi-3-mini-4k-instruct` or `Gemma-2B-it`
- **Format**: GGUF
- **Quantization**: Q4_K_M (recommended)
- **Download**: Run `./scripts/download-models.sh` or `make models`

## Directory Structure

```
models/
├── whisper-small.bin         # Whisper model
├── phi-3-mini-4k-q4.gguf     # LLM model
└── .gitkeep
```

## Model Sizes

| Model                 | Size   | RAM Usage | Speed     |
| --------------------- | ------ | --------- | --------- |
| whisper-small Q5_K_M  | ~500MB | ~1.5GB    | Fast      |
| whisper-medium Q5_K_M | ~1.5GB | ~3GB      | Medium    |
| Phi-3-mini Q4_K_M     | ~2GB   | ~4GB      | Fast      |
| Gemma-2B Q4_K_M       | ~1.5GB | ~3GB      | Very Fast |

## Download Instructions

### Automatic (Recommended)

```bash
make models
```

Or:

```bash
./scripts/download-models.sh
```

### Manual Download

See the [models download guide](../docs/MODELS.md) for manual download instructions.

## Note

⚠️ **Models are NOT included in git** due to their size. They must be downloaded separately.

The `.gitignore` file is configured to ignore `*.gguf` and `*.bin` files in this directory.
