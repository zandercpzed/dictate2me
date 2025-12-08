// Módulo de captura de áudio real (Web Audio API)
// Interface compatível com AudioCapture

export interface AudioCapture {
  start(): void;
  stop(): void;
  onData(callback: (chunk: Uint8Array) => void): void;
}

export class WebAudioCapture implements AudioCapture {
  private running = false;
  private callback?: (chunk: Uint8Array) => void;
  private mediaStream?: MediaStream;
  private audioContext?: AudioContext;
  private processor?: ScriptProcessorNode;

  async start() {
    this.running = true;
    console.log('[AUDIO]', { status: 'start', time: new Date().toISOString() });
    try {
      this.mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true });
      this.audioContext = new window.AudioContext();
      const source = this.audioContext.createMediaStreamSource(this.mediaStream);
      this.processor = this.audioContext.createScriptProcessor(4096, 1, 1);
      source.connect(this.processor);
      this.processor.connect(this.audioContext.destination);
      this.processor.onaudioprocess = (event) => {
        if (!this.running || !this.callback) return;
        const input = event.inputBuffer.getChannelData(0);
        // Converte Float32 para Int16 PCM
        const pcm = new Int16Array(input.length);
        for (let i = 0; i < input.length; i++) {
          let s = Math.max(-1, Math.min(1, input[i]));
          pcm[i] = s < 0 ? s * 0x8000 : s * 0x7FFF;
        }
        this.callback(new Uint8Array(pcm.buffer));
        console.log('[AUDIO]', { event: 'data', len: pcm.length });
      };
    } catch (err) {
      console.log('[AUDIO]', { event: 'error', err });
    }
  }

  stop() {
    this.running = false;
    if (this.processor) this.processor.disconnect();
    if (this.audioContext) this.audioContext.close();
    if (this.mediaStream) this.mediaStream.getTracks().forEach(t => t.stop());
    console.log('[AUDIO]', { status: 'stop', time: new Date().toISOString() });
  }

  onData(callback: (chunk: Uint8Array) => void) {
    this.callback = callback;
    console.log('[AUDIO]', { status: 'onData', time: new Date().toISOString() });
  }
}
