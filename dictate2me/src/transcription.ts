// Módulo de transcrição (simulado para MVP)
// Interface e implementação fake para testes e integração

export interface Transcriber {
  transcribe(chunk: Uint8Array): string;
}

export class SimulatedTranscriber implements Transcriber {
  // Simula transcrição: converte bytes em string
  transcribe(chunk: Uint8Array): string {
    const text = `Texto simulado (${chunk.join(',')})`;
    console.log('[TRANSCRIPTION]', { event: 'transcribe', chunk, text });
    return text;
  }
}

// Transcrição real via API local (HTTP)
export class ApiTranscriber implements Transcriber {
  async transcribe(audioChunk: ArrayBuffer): Promise<string> {
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 4000);
      const response = await fetch('http://localhost:8080/api/v1/transcribe', {
        method: 'POST',
        body: audioChunk,
        headers: { 'Content-Type': 'application/octet-stream' },
        signal: controller.signal
      });
      clearTimeout(timeout);
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      const { text } = await response.json();
      console.log('[API][TRANSCRIBE]', { text });
      return text || '';
    } catch (err) {
      console.log('[API][TRANSCRIBE][ERRO]', { err });
      return '';
    }
  }
}
