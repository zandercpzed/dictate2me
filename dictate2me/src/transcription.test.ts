import { SimulatedTranscriber } from './transcription';

console.log('[TEST][TRANSCRIPTION]', { status: 'start' });
const transcriber = new SimulatedTranscriber();
const chunk = new Uint8Array([1,2,3]);
const result = transcriber.transcribe(chunk);
console.log('[TEST][TRANSCRIPTION]', { result });
