import { SimulatedAudioCapture } from './audio';

console.log('[TEST][AUDIO]', { status: 'start' });
const audio = new SimulatedAudioCapture();
let dataCalled = false;
audio.onData(() => { dataCalled = true; });
audio.start();
setTimeout(() => {
  audio.stop();
  console.log('[TEST][AUDIO]', { dataCalled });
}, 1000);
