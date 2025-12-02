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
	async startDictation() {
		// Check if daemon is running
		if (this.settings.autoCheckDaemon) {
			const isRunning = await this.checkDaemonHealth();
			if (!isRunning) {
				new Notice('❌ Dictate2Me daemon is not running. Please start it first.');
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
				this.client = new Dictate2MeClient(
					this.settings.apiUrl,
					this.settings.apiToken
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
			const response = await fetch(`${this.settings.apiUrl}/health`);
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

		// API URL
		new Setting(containerEl)
			.setName('API URL')
			.setDesc('URL of the Dictate2Me daemon API')
			.addText((text) =>
				text
					.setPlaceholder('http://localhost:8765/api/v1')
					.setValue(this.plugin.settings.apiUrl)
					.onChange(async (value) => {
						this.plugin.settings.apiUrl = value;
						await this.plugin.saveSettings();
					})
			);

		// API Token
		new Setting(containerEl)
			.setName('API Token')
			.setDesc('Authentication token from ~/.dictate2me/api-token')
			.addText((text) =>
				text
					.setPlaceholder('Enter your API token')
					.setValue(this.plugin.settings.apiToken)
					.onChange(async (value) => {
						this.plugin.settings.apiToken = value;
						await this.plugin.saveSettings();
					})
			);

		// Language
		new Setting(containerEl)
			.setName('Language')
			.setDesc('Transcription language (e.g., pt, en, es)')
			.addText((text) =>
				text
					.setPlaceholder('pt')
					.setValue(this.plugin.settings.language)
					.onChange(async (value) => {
						this.plugin.settings.language = value;
						await this.plugin.saveSettings();
					})
			);

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

		// Test connection button
		new Setting(containerEl)
			.setName('Test connection')
			.setDesc('Test connection to Dictate2Me daemon')
			.addButton((button) =>
				button
					.setButtonText('Test')
					.setCta()
					.onClick(async () => {
						button.setButtonText('Testing...');
						button.setDisabled(true);

						const isHealthy = await this.plugin.checkDaemonHealth();

						if (isHealthy) {
							new Notice('✅ Connection successful!');
							button.setButtonText('Success ✓');
						} else {
							new Notice('❌ Connection failed. Is the daemon running?');
							button.setButtonText('Failed ✗');
						}

						setTimeout(() => {
							button.setButtonText('Test');
							button.setDisabled(false);
						}, 2000);
					})
			);
	}
}
