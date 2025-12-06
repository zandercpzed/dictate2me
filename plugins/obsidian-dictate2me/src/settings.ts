/**
 * Plugin settings interface
 */
export interface Dictate2MeSettings {
	language: string;
	enableCorrection: boolean;
	showPartialResults: boolean;
	showConfidence: boolean;
	autoCheckDaemon: boolean;
	groqApiKey: string; // Groq API Key for transcription
}

/**
 * Default settings
 */
export const DEFAULT_SETTINGS: Dictate2MeSettings = {
	language: 'pt',
	enableCorrection: true,
	showPartialResults: true,
	showConfidence: true,
	autoCheckDaemon: true,
	groqApiKey: '', // User must configure
};
