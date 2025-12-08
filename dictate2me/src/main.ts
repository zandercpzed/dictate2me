// Dictate2Me — Plugin Obsidian para transcrição de voz 100% local
// MVP: captura, transcrição, correção simuladas, integração UI e inserção no editor

import { Plugin, MarkdownView } from 'obsidian';
import { WebAudioCapture } from './audio';
import { ApiTranscriber, SimulatedTranscriber } from './transcription';
import { ApiCorrector, SimulatedCorrector } from './correction';

export default class Dictate2MePlugin extends Plugin {
  statusBarItem: HTMLElement | null = null;

  // Inicialização do plugin
  onload() {
    console.log('[Dictate2MePlugin]', { status: 'onload', time: new Date().toISOString() });
    // Botão na ribbon
    this.addRibbonIcon('microphone', 'Iniciar transcrição de voz', () => {
      this.startTranscription();
    });
    // Status bar
    this.statusBarItem = this.addStatusBarItem();
    this.statusBarItem.setText('Dictate2Me: Pronto');
    // Comando na paleta
    this.addCommand({
      id: 'start-transcription',
      name: 'Iniciar transcrição de voz',
      callback: () => {
        this.startTranscription();
      }
    });
    this.addCommand({
      id: 'repeat-last-transcription',
      name: 'Repetir última transcrição',
      callback: () => {
        this.repeatLastTranscription();
      }
    });
  }

  // Inicia ciclo de captura, transcrição e correção
  async startTranscription() {
    if (this.statusBarItem) this.statusBarItem.setText('🎤 Dictate2Me: Gravando...');
    console.log('[COMMAND]', { id: 'start-transcription', status: 'executed', time: new Date().toISOString() });
    const audio = new WebAudioCapture();
    const transcriber = new ApiTranscriber();
    const fallbackTranscriber = new SimulatedTranscriber();
    const corrector = new ApiCorrector();
    const fallbackCorrector = new SimulatedCorrector();
    let lastCorrected = '';
    let errorOccurred = false;
    audio.onData(async (chunk) => {
      let transcript = '';
      let corrected = '';
      try {
        transcript = await transcriber.transcribe(chunk);
        if (!transcript) throw new Error('Transcrição vazia');
      } catch (err) {
        console.log('[TRANSCRIBE][ERRO]', { err });
        transcript = fallbackTranscriber.transcribe(chunk);
        errorOccurred = true;
      }
      try {
        corrected = await corrector.correct(transcript);
        if (!corrected) throw new Error('Correção vazia');
      } catch (err) {
        console.log('[CORRECT][ERRO]', { err });
        corrected = fallbackCorrector.correct(transcript);
        errorOccurred = true;
      }
      lastCorrected = corrected;
      if (this.statusBarItem) {
        this.statusBarItem.setText(errorOccurred
          ? `⚠️ Dictate2Me: ${corrected}`
          : `✅ Dictate2Me: ${corrected}`);
      }
      console.log('[RESULT]', { transcript, corrected, errorOccurred });
    });
    audio.start();
    setTimeout(() => {
      audio.stop();
      if (this.statusBarItem) this.statusBarItem.setText('📝 Dictate2Me: Pronto');
      // Inserir texto no editor ativo
      const view = this.app.workspace.getActiveViewOfType(MarkdownView);
      if (view && lastCorrected) {
        view.editor.replaceSelection(lastCorrected + '\n');
        console.log('[MAIN]', { event: 'inserted', text: lastCorrected });
      } else {
        console.log('[MAIN]', { event: 'insert-failed', reason: 'No active editor or no text' });
        if (this.statusBarItem) this.statusBarItem.setText('Dictate2Me: Falha ao inserir');
      }
    }, 3000);
  }

  repeatLastTranscription() {
    const view = this.app.workspace.getActiveViewOfType(MarkdownView);
    if (view && this.statusBarItem) {
      const text = this.statusBarItem.innerText.replace(/^[^:]+: /, '');
      if (text && text !== 'Pronto' && text !== 'Desativado' && text !== 'Falha ao inserir') {
        view.editor.replaceSelection(text + '\n');
        this.statusBarItem.setText('🔁 Dictate2Me: Repetido');
        console.log('[MAIN]', { event: 'repeat', text });
      } else {
        this.statusBarItem.setText('Dictate2Me: Nada para repetir');
      }
    }
  }

  // Limpeza ao descarregar plugin
  onunload() {
    console.log('[Dictate2MePlugin]', { status: 'onunload', time: new Date().toISOString() });
    if (this.statusBarItem) this.statusBarItem.setText('Dictate2Me: Desativado');
  }
}
