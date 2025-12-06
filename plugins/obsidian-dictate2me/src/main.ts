import { App, Editor, MarkdownView, Modal, Notice, Plugin, PluginSettingTab, Setting } from 'obsidian';
import { Dictate2MeClient } from './client';
import { Dictate2MeSettings, DEFAULT_SETTINGS } from './settings';
import { spawn } from 'child_process';
import { readFileSync } from 'fs';

/**
 * Main plugin class for Dictate2Me
 */
export default class Dictate2MePlugin extends Plugin {
	settings: Dictate2MeSettings;
	client: Dictate2MeClient | null = null;
	statusBarItem: HTMLElement;
	ribbonIconEl: HTMLElement | null = null;
	isRecording = false;
	
	// Auto-detected configuration
	private readonly API_URL = 'http://localhost:8765/api/v1';
	private readonly PROJECT_PATH = '/Users/zander/Library/CloudStorage/GoogleDrive-zander.cattapreta@zedicoes.com/My Drive/_ programação/_ dictate2me/dictate2me';
	private readonly DAEMON_START_SCRIPT = '/Users/zander/Library/CloudStorage/GoogleDrive-zander.cattapreta@zedicoes.com/My Drive/_ programação/_ dictate2me/dictate2me/scripts/start-daemon.sh';

	async onload() {
		await this.loadSettings();

		// Add ribbon icon
		this.ribbonIconEl = this.addRibbonIcon(
			'microphone',
			'Start Dictation',
			async (evt: MouseEvent) => {
				await this.toggleDictation();
			}
		);

		// Add status bar item
		this.statusBarItem = this.addStatusBarItem();
		this.updateStatusBar('Ready');

		// Add command
		this.addCommand({
			id: 'start-dictation',
			name: 'Start/Stop Dictation',
			editorCallback: async (editor: Editor, view: MarkdownView) => {
				await this.toggleDictation();
			},
			hotkeys: [
				{
					modifiers: ['Mod', 'Shift'],
					key: 'D',
				},
			],
		});

		// Add settings tab
		this.addSettingTab(new Dictate2MeSettingTab(this.app, this));

		console.log('Dictate2Me plugin loaded');
	}

	onunload() {
		if (this.client) {
			this.client.disconnect();
		}
		console.log('Dictate2Me plugin unloaded');
	}

	async loadSettings() {
		this.settings = Object.assign({}, DEFAULT_SETTINGS, await this.loadData());
	}

	async saveSettings() {
		await this.saveData(this.settings);
	}


	/**
	 * Get API token from file system
	 */
	private async getApiToken(): Promise<string> {
		try {
			const homeDir = process.env.HOME || process.env.USERPROFILE;
			const tokenPath = `${homeDir}/.dictate2me/api-token`;
			const token = readFileSync(tokenPath, 'utf8');
			return token.trim();
		} catch (error) {
			console.error('Failed to read API token:', error);
			new Notice('Failed to read API token: ' + error.message);
		}
		return '';
	}

	/**
	 * Start daemon automatically using child_process
	 */
	private async startDaemonAutomatically(): Promise<boolean> {
		try {
			console.log('Attempting to auto-start daemon...');
			
			// Copy command to clipboard as a backup
			const projectRoot = this.DAEMON_START_SCRIPT.replace('/scripts/start-daemon.sh', '');
			const command = `cd '${projectRoot}' && ./scripts/start-daemon.sh`;
			await navigator.clipboard.writeText(command);

			// Attempt to spawn the process
			const subprocess = spawn(this.DAEMON_START_SCRIPT, [], {
				detached: true,
				stdio: 'ignore'
			});

			subprocess.unref();
			
			// Wait a bit to see if it starts
			await new Promise(resolve => setTimeout(resolve, 3000));
			
			// Check health
			const isHealthy = await this.checkDaemonHealth();
			return isHealthy;

		} catch (error) {
			console.error('Failed to auto-start daemon:', error);
			return false;
		}
	}

	/**
	 * Toggle dictation on/off
	 */
	async toggleDictation() {
		if (this.isRecording) {
			await this.stopDictation();
		} else {
			await this.startDictation();
		}
	}

	/**
	 * Start dictation with visual feedback
	 */
	async startDictation() {
		// Create and open the startup modal
		const modal = new StartupModal(this.app, this);
		modal.open();

		try {
			// STEP 0: Quick daemon health check (in background)
			modal.setMessage('Checking daemon status...');
			let isHealthy = await this.checkDaemonHealth();

			if (!isHealthy) {
				// Try auto-start
				modal.setMessage('Daemon not running. Attempting auto-start...');
				isHealthy = await this.startDaemonAutomatically();
				
				if (!isHealthy) {
					// Daemon is not running and auto-start failed - show manual instructions
					modal.updateStep('daemon', 'error');
					modal.setMessage('Failed to start daemon automatically.');
					modal.showManualInstructions(this.DAEMON_START_SCRIPT);
					return; // Keep modal open to show instructions
				}
			}

			// STEP 1: Check Transcriber Engine (highest priority)
			modal.updateStep('transcriber', 'loading');
			modal.setMessage('Checking transcription engine...');
			
			// STEP 2: Check Ollama (can run in parallel with other checks)
			modal.updateStep('ollama', 'loading');
			
			// Parallel fetch of daemon status (includes both services)
			const health = await this.getDaemonStatus();
			
			// Evaluate Transcriber Status
			if (health.services.transcription === 'ready') {
				modal.updateStep('transcriber', 'done');
			} else {
				modal.updateStep('transcriber', 'error');
				modal.setMessage('⚠️ Transcriber issue: ' + health.services.transcription);
				new Notice('⚠️ Transcriber not ready: ' + health.services.transcription);
				return; // Don't proceed if transcriber is not ready
			}

			// Evaluate Ollama Status
			if (health.services.correction === 'ready') {
				modal.updateStep('ollama', 'done');
			} else {
				// It's okay if correction is disabled, just warn
				modal.updateStep('ollama', 'warning'); 
				modal.setMessage('Ollama not ready. Correction disabled.');
			}

			// STEP 3: Mark Daemon as healthy (already confirmed at start)
			modal.updateStep('daemon', 'done');
			modal.setMessage('All services ready. Connecting audio...');

			// STEP 4: Connect Audio Stream
			modal.updateStep('audio', 'loading');
			
			// Get active editor
			const view = this.app.workspace.getActiveViewOfType(MarkdownView);
			if (!view) {
				modal.updateStep('audio', 'error');
				modal.setMessage('No active note found to transcribe to.');
				return;
			}

			const editor = view.editor;

			// Initialize Client if needed
			if (!this.client) {
				const token = await this.getApiToken();
				this.client = new Dictate2MeClient(this.API_URL, token);
				this.setupClientHandlers(editor);
			}

			// Connect
			await this.client.connect({
				language: this.settings.language,
				enableCorrection: this.settings.enableCorrection,
			});

			modal.updateStep('audio', 'done');
			
			// Success! Close modal and start
			setTimeout(() => {
				modal.close();
				this.isRecording = true;
				this.updateStatusBar('🎤 Recording...');
				this.updateRibbonIcon();
				new Notice('🎙️ Listening...');
			}, 500);

		} catch (error) {
			console.error(error);
            const errorMessage = error instanceof Error ? error.message : String(error);
			modal.setMessage('Error: ' + errorMessage);
			modal.updateStep('audio', 'error');
		}
	}

	/**
	 * Get detailed daemon status
	 */
	async getDaemonStatus(): Promise<any> {
		try {
			const response = await fetch(`${this.API_URL}/health`);
			if (response.ok) {
				return await response.json();
			}
		} catch (e) { /* ignore */ }
		return { services: { transcription: 'unknown', correction: 'unknown' } };
	}

	/**
	 * Setup client event handlers
	 */
	private setupClientHandlers(editor: Editor) {
		if (!this.client) return;

		this.client.on('partial', (text: string) => {
			if (this.settings.showPartialResults) {
				this.updateStatusBar(`💭 ${text}`);
			}
		});

		this.client.on('final', (data: { transcript: string; corrected: string; confidence: number }) => {
			const textToInsert = this.settings.enableCorrection && data.corrected
				? data.corrected
				: data.transcript;

			if (textToInsert) {
				const cursor = editor.getCursor();
				editor.replaceRange(textToInsert + ' ', cursor);
				
				const newCursor = {
					line: cursor.line,
					ch: cursor.ch + textToInsert.length + 1,
				};
				editor.setCursor(newCursor);

				if (this.settings.showConfidence) {
					const confidencePercent = Math.round(data.confidence * 100);
					this.updateStatusBar(`✓ Inserted (${confidencePercent}%)`);
				} else {
					this.updateStatusBar('✓ Inserted');
				}
			}
		});

		this.client.on('error', (error: string) => {
			new Notice(`❌ Error: ${error}`);
			this.updateStatusBar('Error');
			this.isRecording = false;
			this.updateRibbonIcon();
		});
	}

	/**
	 * Stop dictation
	 */
	async stopDictation() {
		if (this.client) {
			try {
				await this.client.disconnect();
				this.isRecording = false;
				this.updateStatusBar('Ready');
				this.updateRibbonIcon();
				new Notice('⏹️ Dictation stopped');
			} catch (error) {
				new Notice(`❌ Failed to stop dictation: ${error}`);
				console.error('Stop dictation error:', error);
			}
		}
	}

	/**
	 * Check if daemon is healthy
	 */
	async checkDaemonHealth(): Promise<boolean> {
		try {
			const response = await fetch(`${this.API_URL}/health`);
			if (response.ok) {
				const data = await response.json();
				return data.status === 'healthy';
			}
			return false;
		} catch (error) {
			return false;
		}
	}

	/**
	 * Update status bar text
	 */
	private updateStatusBar(text: string) {
		this.statusBarItem.setText(`Dictate2Me: ${text}`);
	}

	/**
	 * Update ribbon icon based on recording status
	 */
	private updateRibbonIcon() {
		if (!this.ribbonIconEl) return;

		if (this.isRecording) {
			this.ribbonIconEl.addClass('dictate2me-recording');
		} else {
			this.ribbonIconEl.removeClass('dictate2me-recording');
		}
	}
}

/**
 * Startup Modal Class
 */
class StartupModal extends Modal {
	plugin: Dictate2MePlugin;
	steps: { [key: string]: { bar: HTMLElement, label: HTMLElement } } = {};
	messageEl: HTMLElement;
	manualContainer: HTMLElement;

	constructor(app: App, plugin: Dictate2MePlugin) {
		super(app);
		this.plugin = plugin;
	}

	onOpen() {
		const { contentEl } = this;
		contentEl.empty();
		contentEl.addClass('dictate2me-startup-modal');

		contentEl.createEl('h2', { text: '🚀 Starting Dictate2Me' });

		const stepList = contentEl.createDiv({ cls: 'step-list' });

		this.createStep(stepList, 'transcriber', 'Transcriber Engine');
		this.createStep(stepList, 'ollama', 'Ollama (LLM)');
		this.createStep(stepList, 'daemon', 'Daemon Services');
		this.createStep(stepList, 'audio', 'Audio Stream');

		this.messageEl = contentEl.createDiv({ cls: 'startup-message' });
		this.messageEl.setText('Initializing...');
		
		this.manualContainer = contentEl.createDiv({ cls: 'manual-instructions' });
		this.manualContainer.hide();

		this.addStyles();
	}

	createStep(container: HTMLElement, id: string, text: string) {
		const stepContainer = container.createDiv({ cls: 'startup-step-container' });
		
		// Label
		const label = stepContainer.createDiv({ cls: 'step-label', text: text });
		
		// Segmented Bar
		const barContainer = stepContainer.createDiv({ cls: 'segmented-bar' });
		// Create 10 segments
		for (let i = 0; i < 10; i++) {
			barContainer.createDiv({ cls: 'bar-segment' });
		}

		this.steps[id] = { bar: barContainer, label: label };
	}

	updateStep(id: string, status: 'pending' | 'loading' | 'done' | 'error' | 'warning') {
		const step = this.steps[id];
		if (!step) return;

		const segments = step.bar.querySelectorAll('.bar-segment');
		step.bar.removeClass('status-loading', 'status-done', 'status-error');

		switch (status) {
			case 'pending':
				segments.forEach(el => (el as HTMLElement).style.backgroundColor = 'var(--background-modifier-border)');
				break;
			case 'loading':
				step.bar.addClass('status-loading');
				// Animation is handled by CSS
				break;
			case 'done':
				step.bar.addClass('status-done');
				segments.forEach(el => (el as HTMLElement).style.backgroundColor = '#4CAF50'); // Green
				break;
			case 'error':
				step.bar.addClass('status-error');
				segments.forEach(el => (el as HTMLElement).style.backgroundColor = '#F44336'); // Red
				break;
			case 'warning':
				segments.forEach(el => (el as HTMLElement).style.backgroundColor = '#FFC107'); // Yellow
				break;
		}
	}

	setMessage(text: string) {
		if (this.messageEl) this.messageEl.setText(text);
	}

	showManualInstructions(scriptPath: string) {
		this.manualContainer.empty();
		this.manualContainer.show();
		
		this.manualContainer.createEl('h3', { text: '⚠️ Manual Start Required' });
		this.manualContainer.createEl('p', { text: 'Please open Terminal and run this command (already copied to clipboard):' });
		
		const projectRoot = scriptPath.replace('/scripts/start-daemon.sh', '');
		const fullCommand = `cd '${projectRoot}' && ./scripts/start-daemon.sh`;
		
		const codeBlock = this.manualContainer.createEl('pre');
		codeBlock.createEl('code', { text: fullCommand });
		
		const buttonContainer = this.manualContainer.createDiv({ cls: 'button-container' });
		
		const copyBtn = buttonContainer.createEl('button', { text: '📋 Copy Command' });
		copyBtn.onclick = () => {
			navigator.clipboard.writeText(fullCommand);
			copyBtn.setText('✅ Copied!');
			setTimeout(() => copyBtn.setText('📋 Copy Command'), 2000);
		};
		
		const retryBtn = buttonContainer.createEl('button', { text: '🔄 Check Again' });
		retryBtn.onclick = async () => {
			retryBtn.setText('⏳ Checking...');
			const isHealthy = await this.plugin.checkDaemonHealth();
			if (isHealthy) {
				this.close();
				new Notice('✅ Daemon is running! Starting dictation...');
				// Restart the dictation process
				this.plugin.startDictation();
			} else {
				retryBtn.setText('❌ Still not running');
				setTimeout(() => retryBtn.setText('🔄 Check Again'), 2000);
			}
		};
	}

	addStyles() {
		const style = document.createElement('style');
		style.textContent = `
			.dictate2me-startup-modal .step-list { margin: 20px 0; }
			.dictate2me-startup-modal .startup-step-container { margin-bottom: 15px; }
			.dictate2me-startup-modal .step-label { margin-bottom: 5px; font-weight: bold; font-size: 0.9em; }
			
			/* Segmented Bar Styles */
			.dictate2me-startup-modal .segmented-bar {
				display: flex;
				gap: 4px;
				padding: 4px;
				border: 2px solid var(--text-muted);
				border-radius: 4px;
				height: 24px;
				background: var(--background-primary);
			}

			.dictate2me-startup-modal .bar-segment {
				flex: 1;
				background-color: var(--background-modifier-border);
				border-radius: 1px;
				transition: background-color 0.2s ease;
			}

			/* Loading Animation */
			@keyframes segment-load {
				0% { background-color: var(--background-modifier-border); }
				50% { background-color: #4CAF50; }
				100% { background-color: var(--background-modifier-border); }
			}

			.dictate2me-startup-modal .status-loading .bar-segment {
				animation: segment-load 1s infinite;
			}
			
			/* Stagger animations for wave effect */
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(1) { animation-delay: 0.0s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(2) { animation-delay: 0.1s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(3) { animation-delay: 0.2s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(4) { animation-delay: 0.3s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(5) { animation-delay: 0.4s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(6) { animation-delay: 0.5s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(7) { animation-delay: 0.6s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(8) { animation-delay: 0.7s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(9) { animation-delay: 0.8s; }
			.dictate2me-startup-modal .status-loading .bar-segment:nth-child(10) { animation-delay: 0.9s; }

			.dictate2me-startup-modal .status-done .segmented-bar {
				border-color: #4CAF50;
			}
			
			.dictate2me-startup-modal .startup-message { 
				margin-top: 20px; 
				font-style: italic; 
				color: var(--text-muted); 
				text-align: center;
			}
			
			.dictate2me-startup-modal .manual-instructions {
				margin-top: 20px;
				padding: 15px;
				background-color: var(--background-secondary);
				border-radius: 8px;
				border: 1px solid var(--interactive-error);
			}
			.dictate2me-startup-modal .manual-instructions h3 { margin-top: 0; color: var(--interactive-error); }
			.dictate2me-startup-modal .manual-instructions pre { background: var(--background-primary); padding: 10px; overflow-x: auto; }
			.dictate2me-startup-modal .manual-instructions .button-container {
				display: flex;
				gap: 10px;
				margin-top: 15px;
				justify-content: center;
			}
			.dictate2me-startup-modal .manual-instructions button {
				padding: 8px 16px;
				border-radius: 4px;
				font-size: 14px;
				cursor: pointer;
				transition: opacity 0.2s;
			}
			.dictate2me-startup-modal .manual-instructions button:hover {
				opacity: 0.8;
			}
		`;
		document.head.appendChild(style);
	}

	onClose() {
		const { contentEl } = this;
		contentEl.empty();
	}
}

/**
 * Settings tab
 */
class Dictate2MeSettingTab extends PluginSettingTab {
	plugin: Dictate2MePlugin;

	constructor(app: App, plugin: Dictate2MePlugin) {
		super(app, plugin);
		this.plugin = plugin;
	}

	display(): void {
		const { containerEl } = this;
		containerEl.empty();

		containerEl.createEl('h2', { text: 'Dictate2Me Settings' });

		// Transcription section
		containerEl.createEl('h3', { text: 'Transcription Settings' });

		// Language dropdown
		new Setting(containerEl)
			.setName('Language')
			.setDesc('Transcription language')
			.addDropdown((dropdown) =>
				dropdown
					.addOption('pt', '🇧🇷 Portuguese (pt)')
					.addOption('en', '🇺🇸 English (en)')
					.addOption('es', '🇪🇸 Spanish (es)')
					.addOption('fr', '🇫🇷 French (fr)')
					.addOption('de', '🇩🇪 German (de)')
					.addOption('it', '🇮🇹 Italian (it)')
					.addOption('ru', '🇷🇺 Russian (ru)')
					.addOption('zh', '🇨🇳 Chinese (zh)')
					.setValue(this.plugin.settings.language)
					.onChange(async (value) => {
						this.plugin.settings.language = value;
						await this.plugin.saveSettings();
					})
			);

		// Features section
		containerEl.createEl('h3', { text: 'Features' });

		// Enable correction
		new Setting(containerEl)
			.setName('Enable text correction')
			.setDesc('Use LLM to correct grammar, punctuation, and capitalization')
			.addToggle((toggle) =>
				toggle
					.setValue(this.plugin.settings.enableCorrection)
					.onChange(async (value) => {
						this.plugin.settings.enableCorrection = value;
						await this.plugin.saveSettings();
					})
			);

		// Show partial results
		new Setting(containerEl)
			.setName('Show partial results')
			.setDesc('Display partial transcription results in status bar')
			.addToggle((toggle) =>
				toggle
					.setValue(this.plugin.settings.showPartialResults)
					.onChange(async (value) => {
						this.plugin.settings.showPartialResults = value;
						await this.plugin.saveSettings();
					})
			);

		// Show confidence
		new Setting(containerEl)
			.setName('Show confidence score')
			.setDesc('Display transcription confidence percentage')
			.addToggle((toggle) =>
				toggle
					.setValue(this.plugin.settings.showConfidence)
					.onChange(async (value) => {
						this.plugin.settings.showConfidence = value;
						await this.plugin.saveSettings();
					})
			);

		// Auto-check daemon
		new Setting(containerEl)
			.setName('Auto-check daemon')
			.setDesc('Check if daemon is running before starting dictation')
			.addToggle((toggle) =>
				toggle
					.setValue(this.plugin.settings.autoCheckDaemon)
					.onChange(async (value) => {
						this.plugin.settings.autoCheckDaemon = value;
						await this.plugin.saveSettings();
					})
			);
	}
}
