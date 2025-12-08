"use strict";
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/main.ts
var main_exports = {};
__export(main_exports, {
  default: () => Dictate2MePlugin
});
module.exports = __toCommonJS(main_exports);
var import_obsidian = require("obsidian");

// src/audio.ts
var WebAudioCapture = class {
  constructor() {
    this.running = false;
  }
  async start() {
    this.running = true;
    console.log("[AUDIO]", { status: "start", time: (/* @__PURE__ */ new Date()).toISOString() });
    try {
      this.mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true });
      this.audioContext = new window.AudioContext();
      const source = this.audioContext.createMediaStreamSource(this.mediaStream);
      this.processor = this.audioContext.createScriptProcessor(4096, 1, 1);
      source.connect(this.processor);
      this.processor.connect(this.audioContext.destination);
      this.processor.onaudioprocess = (event) => {
        if (!this.running || !this.callback)
          return;
        const input = event.inputBuffer.getChannelData(0);
        const pcm = new Int16Array(input.length);
        for (let i = 0; i < input.length; i++) {
          let s = Math.max(-1, Math.min(1, input[i]));
          pcm[i] = s < 0 ? s * 32768 : s * 32767;
        }
        this.callback(new Uint8Array(pcm.buffer));
        console.log("[AUDIO]", { event: "data", len: pcm.length });
      };
    } catch (err) {
      console.log("[AUDIO]", { event: "error", err });
    }
  }
  stop() {
    this.running = false;
    if (this.processor)
      this.processor.disconnect();
    if (this.audioContext)
      this.audioContext.close();
    if (this.mediaStream)
      this.mediaStream.getTracks().forEach((t) => t.stop());
    console.log("[AUDIO]", { status: "stop", time: (/* @__PURE__ */ new Date()).toISOString() });
  }
  onData(callback) {
    this.callback = callback;
    console.log("[AUDIO]", { status: "onData", time: (/* @__PURE__ */ new Date()).toISOString() });
  }
};

// src/transcription.ts
var SimulatedTranscriber = class {
  // Simula transcrição: converte bytes em string
  transcribe(chunk) {
    const text = `Texto simulado (${chunk.join(",")})`;
    console.log("[TRANSCRIPTION]", { event: "transcribe", chunk, text });
    return text;
  }
};
var ApiTranscriber = class {
  async transcribe(audioChunk) {
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 4e3);
      const response = await fetch("http://localhost:8080/api/v1/transcribe", {
        method: "POST",
        body: audioChunk,
        headers: { "Content-Type": "application/octet-stream" },
        signal: controller.signal
      });
      clearTimeout(timeout);
      if (!response.ok)
        throw new Error(`HTTP ${response.status}`);
      const { text } = await response.json();
      console.log("[API][TRANSCRIBE]", { text });
      return text || "";
    } catch (err) {
      console.log("[API][TRANSCRIBE][ERRO]", { err });
      return "";
    }
  }
};

// src/correction.ts
var SimulatedCorrector = class {
  // Simula correção: adiciona pontuação
  correct(text) {
    const corrected = text + ".";
    console.log("[CORRECTION]", { event: "correct", original: text, corrected });
    return corrected;
  }
};
var ApiCorrector = class {
  async correct(text) {
    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 4e3);
      const response = await fetch("http://localhost:8080/api/v1/correct", {
        method: "POST",
        body: JSON.stringify({ text }),
        headers: { "Content-Type": "application/json" },
        signal: controller.signal
      });
      clearTimeout(timeout);
      if (!response.ok)
        throw new Error(`HTTP ${response.status}`);
      const { corrected } = await response.json();
      console.log("[API][CORRECT]", { corrected });
      return corrected || "";
    } catch (err) {
      console.log("[API][CORRECT][ERRO]", { err });
      return "";
    }
  }
};

// src/main.ts
var Dictate2MePlugin = class extends import_obsidian.Plugin {
  constructor() {
    super(...arguments);
    this.statusBarItem = null;
  }
  // Inicialização do plugin
  onload() {
    console.log("[Dictate2MePlugin]", { status: "onload", time: (/* @__PURE__ */ new Date()).toISOString() });
    this.addRibbonIcon("microphone", "Iniciar transcri\xE7\xE3o de voz", () => {
      this.startTranscription();
    });
    this.statusBarItem = this.addStatusBarItem();
    this.statusBarItem.setText("Dictate2Me: Pronto");
    this.addCommand({
      id: "start-transcription",
      name: "Iniciar transcri\xE7\xE3o de voz",
      callback: () => {
        this.startTranscription();
      }
    });
    this.addCommand({
      id: "repeat-last-transcription",
      name: "Repetir \xFAltima transcri\xE7\xE3o",
      callback: () => {
        this.repeatLastTranscription();
      }
    });
  }
  // Inicia ciclo de captura, transcrição e correção
  async startTranscription() {
    if (this.statusBarItem)
      this.statusBarItem.setText("\u{1F3A4} Dictate2Me: Gravando...");
    console.log("[COMMAND]", { id: "start-transcription", status: "executed", time: (/* @__PURE__ */ new Date()).toISOString() });
    const audio = new WebAudioCapture();
    const transcriber = new ApiTranscriber();
    const fallbackTranscriber = new SimulatedTranscriber();
    const corrector = new ApiCorrector();
    const fallbackCorrector = new SimulatedCorrector();
    let lastCorrected = "";
    let errorOccurred = false;
    audio.onData(async (chunk) => {
      let transcript = "";
      let corrected = "";
      try {
        transcript = await transcriber.transcribe(chunk);
        if (!transcript)
          throw new Error("Transcri\xE7\xE3o vazia");
      } catch (err) {
        console.log("[TRANSCRIBE][ERRO]", { err });
        transcript = fallbackTranscriber.transcribe(chunk);
        errorOccurred = true;
      }
      try {
        corrected = await corrector.correct(transcript);
        if (!corrected)
          throw new Error("Corre\xE7\xE3o vazia");
      } catch (err) {
        console.log("[CORRECT][ERRO]", { err });
        corrected = fallbackCorrector.correct(transcript);
        errorOccurred = true;
      }
      lastCorrected = corrected;
      if (this.statusBarItem) {
        this.statusBarItem.setText(errorOccurred ? `\u26A0\uFE0F Dictate2Me: ${corrected}` : `\u2705 Dictate2Me: ${corrected}`);
      }
      console.log("[RESULT]", { transcript, corrected, errorOccurred });
    });
    audio.start();
    setTimeout(() => {
      audio.stop();
      if (this.statusBarItem)
        this.statusBarItem.setText("\u{1F4DD} Dictate2Me: Pronto");
      const view = this.app.workspace.getActiveViewOfType(import_obsidian.MarkdownView);
      if (view && lastCorrected) {
        view.editor.replaceSelection(lastCorrected + "\n");
        console.log("[MAIN]", { event: "inserted", text: lastCorrected });
      } else {
        console.log("[MAIN]", { event: "insert-failed", reason: "No active editor or no text" });
        if (this.statusBarItem)
          this.statusBarItem.setText("Dictate2Me: Falha ao inserir");
      }
    }, 3e3);
  }
  repeatLastTranscription() {
    const view = this.app.workspace.getActiveViewOfType(import_obsidian.MarkdownView);
    if (view && this.statusBarItem) {
      const text = this.statusBarItem.innerText.replace(/^[^:]+: /, "");
      if (text && text !== "Pronto" && text !== "Desativado" && text !== "Falha ao inserir") {
        view.editor.replaceSelection(text + "\n");
        this.statusBarItem.setText("\u{1F501} Dictate2Me: Repetido");
        console.log("[MAIN]", { event: "repeat", text });
      } else {
        this.statusBarItem.setText("Dictate2Me: Nada para repetir");
      }
    }
  }
  // Limpeza ao descarregar plugin
  onunload() {
    console.log("[Dictate2MePlugin]", { status: "onunload", time: (/* @__PURE__ */ new Date()).toISOString() });
    if (this.statusBarItem)
      this.statusBarItem.setText("Dictate2Me: Desativado");
  }
};
