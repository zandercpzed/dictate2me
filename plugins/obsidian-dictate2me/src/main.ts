import { App, Editor, MarkdownView, Notice, Plugin, PluginSettingTab, Setting } from 'obsidian';
import { Dictate2MeClient } from './client';
import { Dictate2MeSettings, DEFAULT_SETTINGS } from './settings';

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
			// Try to read from standard location
			const response = await fetch('file:///Users/zander/.dictate2me/api-token');
			if (response.ok) {
				const token = await response.text();
				return token.trim();
			}
		} catch (error) {
			console.error('Failed to read API token:', error);
		}
		return '';
	}

	/**
	 * Start daemon automatically
	 */
	private async startDaemonAutomatically(): Promise<boolean> {
		try {
			new Notice('🚀 Starting dictate2me daemon...');
			
			// Open Terminal and run the start script
			const command = `cd '${this.PROJECT_PATH}' && ./scripts/start-daemon.sh`;
			
			// Use AppleScript to execute in Terminal
			const  script = `tell application "Terminal"\n  do script "${command.replace(/"/g, '\\"')}"\n  activate\nend tell`;
			const encodedScript = encodeURIComponent(script);
			
			// Try to open using URL scheme (may not work in all contexts)
			window.open(`data:text/plain;charset=utf-8,${encodedScript}`, '_blank');
			
			// Give daemon time to start
			await new Promise(resolve => setTimeout(resolve, 5000));
			
			// Check if it started
			const isRunning = await this.checkDaemonHealth();
			if (isRunning) {
				new Notice('✅ Daemon started successfully!');
				return true;
			} else {
				new Notice('⚠️ Daemon may be starting... Please wait a moment and try again.');
				return false;
			}
		} catch (error) {
			console.error('Failed to start daemon:', error);
			new Notice('❌ Couldn\'t start daemon automatically. Please start it manually in Terminal:\ncd ~/dictate2me && ./scripts/start-daemon.sh');
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
	 * Start dictation
	 */
	/**
	 * Start dictation
	 */
	async startDictation() {
		// Check if daemon is running
		let isRunning = await this.checkDaemonHealth();
		
		if (!isRunning) {
			// Try to start automatically
			const started = await this.startDaemonAutomatically();
			if (!started) {
				return;
			}
			// Re-check health
			isRunning = await this.checkDaemonHealth();
			if (!isRunning) {
				new Notice('❌ Daemon started but not responding yet. Please try again in a few seconds.');
				return;
			}
		}

		// Get active editor
		const view = this.app.workspace.getActiveViewOfType(MarkdownView);
		if (!view) {
			new Notice('❌ No active note found');
			return;
		}

		const editor = view.editor;

		try {
			// Create client if not exists
			if (!this.client) {
				const token = await this.getApiToken();
				
				this.client = new Dictate2MeClient(
					this.API_URL,
					token
				);

				// Setup event handlers
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
						// Insert text at cursor
						const cursor = editor.getCursor();
						editor.replaceRange(textToInsert + ' ', cursor);

						// Move cursor to end of inserted text
						const newCursor = {
							line: cursor.line,
							ch: cursor.ch + textToInsert.length + 1,
						};
						editor.setCursor(newCursor);

						// Show confidence if enabled
						if (this.settings.showConfidence) {
							const confidencePercent = Math.round(data.confidence * 100);
							this.updateStatusBar(`✓ Inserted (${confidencePercent}% confidence)`);
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

			// Start recording
			await this.client.connect({
				language: this.settings.language,
				enableCorrection: this.settings.enableCorrection,
			});

			this.isRecording = true;
			this.updateStatusBar('🎤 Recording...');
			this.updateRibbonIcon();
			new Notice('🎤 Dictation started');

		} catch (error) {
			new Notice(`❌ Failed to start dictation: ${error}`);
			console.error('Dictation error:', error);
			this.isRecording = false;
			this.updateRibbonIcon();
		}
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
