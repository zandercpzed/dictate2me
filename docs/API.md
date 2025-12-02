# dictate2me API Documentation

**Version:** 0.2.0  
**Base URL:** `http://localhost:8765/api/v1`  
**Authentication:** Bearer Token

---

## 📋 Table of Contents

- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Health Check](#get-health)
  - [Transcribe Audio](#post-transcribe)
  - [Correct Text](#post-correct)
  - [Stream (WebSocket)](#ws-stream)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Examples](#examples)

---

## 🔐 Authentication

All endpoints (except `/health`) require Bearer token authentication.

### Getting the Token

The token is automatically generated on first daemon startup and saved to:

```
~/.dictate2me/api-token
```

### Using the Token

Include the token in the `Authorization` header:

```http
Authorization: Bearer <your-token-here>
```

### Example

```bash
export TOKEN=$(cat ~/.dictate2me/api-token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8765/api/v1/health
```

---

## 📡 Endpoints

### GET /health

Health check endpoint to verify API and service availability.

**Authentication:** Not required

**Response:**

```json
{
  "status": "healthy",
  "services": {
    "transcription": "ready",
    "correction": "ready"
  },
  "uptime": 3600
}
```

**Fields:**

- `status` (string): Overall status (`healthy` or `unhealthy`)
- `services` (object): Status of each service
  - `transcription` (string): Transcription service status
  - `correction` (string): Correction service status (or `disabled`)
- `uptime` (integer): Server uptime in seconds

**Status Codes:**

- `200 OK`: Server is healthy

**Example:**

```bash
curl http://localhost:8765/api/v1/health
```

---

### POST /transcribe

Transcribe audio data to text.

**Authentication:** Required

**Request Body:**

```json
{
  "audio": "base64_encoded_audio_data",
  "sampleRate": 16000,
  "language": "pt"
}
```

**Fields:**

- `audio` (string, required): Base64-encoded WAV audio data (16-bit PCM)
- `sampleRate` (number, optional): Sample rate in Hz (default: 16000)
- `language` (string, optional): Language code (default: "pt")

**Response:**

```json
{
  "text": "Olá mundo como vai você",
  "confidence": 0.95,
  "segments": [
    {
      "text": "Olá mundo",
      "start": 0.0,
      "end": 1.2,
      "confidence": 0.96
    },
    {
      "text": "como vai você",
      "start": 1.2,
      "end": 2.5,
      "confidence": 0.94
    }
  ]
}
```

**Response Fields:**

- `text` (string): Complete transcribed text
- `confidence` (number): Average confidence score (0.0-1.0)
- `segments` (array): Individual transcription segments
  - `text` (string): Segment text
  - `start` (number): Start time in seconds
  - `end` (number): End time in seconds
  - `confidence` (number): Segment confidence (0.0-1.0)

**Status Codes:**

- `200 OK`: Transcription successful
- `400 Bad Request`: Invalid request (missing/invalid audio data)
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Transcription failed

**Example:**

```bash
# Record audio
sox -d -r 16000 -c 1 -b 16 audio.wav trim 0 5

# Convert to base64
AUDIO_BASE64=$(base64 -i audio.wav)

# Send request
curl -X POST http://localhost:8765/api/v1/transcribe \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"audio\": \"$AUDIO_BASE64\"}"
```

---

### POST /correct

Correct text using LLM (grammar, punctuation, capitalization).

**Authentication:** Required

**Request Body:**

```json
{
  "text": "olá mundo como vai você"
}
```

**Fields:**

- `text` (string, required): Text to correct

**Response:**

```json
{
  "original": "olá mundo como vai você",
  "corrected": "Olá, mundo! Como vai você?",
  "model": "gemma2:2b"
}
```

**Response Fields:**

- `original` (string): Original input text
- `corrected` (string): Corrected text
- `model` (string): LLM model used for correction

**Status Codes:**

- `200 OK`: Correction successful
- `400 Bad Request`: Invalid request (missing text)
- `401 Unauthorized`: Missing or invalid token
- `503 Service Unavailable`: Correction service not available
- `500 Internal Server Error`: Correction failed

**Example:**

```bash
curl -X POST http://localhost:8765/api/v1/correct \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"text": "olá mundo como vai você"}'
```

---

### WS /stream

Real-time transcription and correction via WebSocket.

**Authentication:** Required (via query parameter or initial message)

**Connection:**

```javascript
const ws = new WebSocket("ws://localhost:8765/api/v1/stream");
```

#### Message Types

##### Client → Server

**1. Start Stream**

```json
{
  "type": "start",
  "data": {
    "language": "pt",
    "enableCorrection": true
  }
}
```

**2. Send Audio Chunk**

```json
{
  "type": "audio",
  "data": {
    "data": "base64_audio_chunk"
  }
}
```

**3. Stop Stream**

```json
{
  "type": "stop"
}
```

##### Server → Client

**1. Partial Result**

```json
{
  "type": "partial",
  "data": {
    "text": "resultado parcial em tempo real..."
  }
}
```

**2. Final Result**

```json
{
  "type": "final",
  "data": {
    "transcript": "Texto transcrito completo.",
    "corrected": "Texto corrigido completo.",
    "confidence": 0.95
  }
}
```

**3. Error**

```json
{
  "type": "error",
  "data": {
    "message": "Error message"
  }
}
```

**Example (JavaScript):**

```javascript
// Read token
const token = await fetch("file://~/.dictate2me/api-token").then((r) =>
  r.text()
);

// Connect
const ws = new WebSocket("ws://localhost:8765/api/v1/stream");

// Authenticate and start
ws.onopen = () => {
  // Send auth (if needed)
  ws.send(
    JSON.stringify({
      type: "start",
      data: {
        language: "pt",
        enableCorrection: true,
      },
    })
  );
};

// Handle messages
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  switch (msg.type) {
    case "partial":
      console.log("Partial:", msg.data.text);
      break;
    case "final":
      console.log("Final:", msg.data.corrected);
      insertTextInEditor(msg.data.corrected);
      break;
    case "error":
      console.error("Error:", msg.data.message);
      break;
  }
};

// Send audio chunks
function sendAudioChunk(audioBuffer) {
  const base64 = btoa(String.fromCharCode(...new Uint8Array(audioBuffer)));
  ws.send(
    JSON.stringify({
      type: "audio",
      data: { data: base64 },
    })
  );
}

// Stop
function stop() {
  ws.send(JSON.stringify({ type: "stop" }));
  ws.close();
}
```

---

## 📊 Data Models

### HealthResponse

```typescript
interface HealthResponse {
  status: string; // "healthy" | "unhealthy"
  services: {
    transcription: string; // "ready" | "error: ..."
    correction: string; // "ready" | "disabled" | "unavailable: ..."
  };
  uptime: number; // seconds
}
```

### TranscribeRequest

```typescript
interface TranscribeRequest {
  audio: string; // Base64-encoded WAV (16-bit PCM)
  sampleRate?: number; // Default: 16000
  language?: string; // Default: "pt"
}
```

### TranscribeResponse

```typescript
interface TranscribeResponse {
  text: string; // Full transcribed text
  confidence: number; // 0.0 - 1.0
  segments: Segment[];
}

interface Segment {
  text: string;
  start: number; // seconds
  end: number; // seconds
  confidence: number; // 0.0 - 1.0
}
```

### CorrectRequest

```typescript
interface CorrectRequest {
  text: string; // Text to correct
}
```

### CorrectResponse

```typescript
interface CorrectResponse {
  original: string;
  corrected: string;
  model: string; // e.g., "gemma2:2b"
}
```

### StreamMessage

```typescript
type StreamMessage =
  | { type: "start"; data: StreamStartConfig }
  | { type: "audio"; data: { data: string } }
  | { type: "stop" }
  | { type: "partial"; data: { text: string } }
  | { type: "final"; data: StreamFinalData }
  | { type: "error"; data: { message: string } };

interface StreamStartConfig {
  language: string; // e.g., "pt"
  enableCorrection: boolean;
}

interface StreamFinalData {
  transcript: string;
  corrected?: string;
  confidence: number;
}
```

---

## ⚠️ Error Handling

### Error Response Format

```json
{
  "error": "Error message describing what went wrong"
}
```

### Common Error Codes

| Code | Description           | Common Causes                         |
| ---- | --------------------- | ------------------------------------- |
| 400  | Bad Request           | Invalid JSON, missing required fields |
| 401  | Unauthorized          | Missing or invalid token              |
| 429  | Too Many Requests     | Rate limit exceeded (>100 req/min)    |
| 500  | Internal Server Error | Server error during processing        |
| 503  | Service Unavailable   | Correction service not available      |

### Example Error Response

```json
{
  "error": "missing Authorization header"
}
```

---

## 💡 Examples

### Complete Transcription + Correction Workflow

```bash
#!/bin/bash

# 1. Get token
TOKEN=$(cat ~/.dictate2me/api-token)

# 2. Check if daemon is running
if ! curl -s http://localhost:8765/api/v1/health > /dev/null; then
  echo "Starting daemon..."
  dictate2me-daemon &
  sleep 3
fi

# 3. Record audio (5 seconds)
sox -d -r 16000 -c 1 -b 16 audio.wav trim 0 5

# 4. Transcribe
AUDIO_BASE64=$(base64 -i audio.wav)
TRANSCRIPT=$(curl -s -X POST http://localhost:8765/api/v1/transcribe \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"audio\": \"$AUDIO_BASE64\"}" \
  | jq -r '.text')

echo "Transcript: $TRANSCRIPT"

# 5. Correct
CORRECTED=$(curl -s -X POST http://localhost:8765/api/v1/correct \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"$TRANSCRIPT\"}" \
  | jq -r '.corrected')

echo "Corrected: $CORRECTED"
```

### TypeScript Client Class

```typescript
class Dictate2MeClient {
  private baseUrl: string;
  private token: string;

  constructor(baseUrl = "http://localhost:8765/api/v1", token: string) {
    this.baseUrl = baseUrl;
    this.token = token;
  }

  async health(): Promise<HealthResponse> {
    const response = await fetch(`${this.baseUrl}/health`);
    return response.json();
  }

  async transcribe(audioBase64: string): Promise<TranscribeResponse> {
    const response = await fetch(`${this.baseUrl}/transcribe`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${this.token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ audio: audioBase64 }),
    });

    if (!response.ok) {
      throw new Error(`Transcription failed: ${response.statusText}`);
    }

    return response.json();
  }

  async correct(text: string): Promise<CorrectResponse> {
    const response = await fetch(`${this.baseUrl}/correct`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${this.token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ text }),
    });

    if (!response.ok) {
      throw new Error(`Correction failed: ${response.statusText}`);
    }

    return response.json();
  }

  createStream(): WebSocket {
    return new WebSocket(`ws://localhost:8765/api/v1/stream`);
  }
}

// Usage
const client = new Dictate2MeClient(
  "http://localhost:8765/api/v1",
  "your-token-here"
);

const health = await client.health();
console.log(health);
```

---

## 🔒 Security Considerations

1. **Localhost Only**: API binds to `127.0.0.1` and only accepts local connections
2. **Token Required**: All endpoints (except `/health`) require valid token
3. **Rate Limiting**: 100 requests per minute per IP
4. **CORS**: Only `localhost` origins allowed
5. **No External Access**: Never expose this API to the internet

---

## 📝 Notes

- The API uses JSON for all request/response bodies
- All timestamps are in seconds with decimal precision
- Audio must be 16-bit PCM WAV format
- Sample rate should be 16000 Hz for best results
- WebSocket connections timeout after 60 seconds of inactivity
- Maximum audio chunk size: 10MB (recommended: 1-2 seconds of audio per chunk)

---

**Last Updated:** 2025-12-01
