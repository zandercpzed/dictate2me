/**
 * Plugin settings interface
 */
export interface Dictate2MeSettings {
	apiUrl: string;
	apiToken: string;
	language: string;
	enableCorrection: boolean;
	showPartialResults: boolean;
	showConfidence: boolean;
	autoCheckDaemon: boolean;
}

/**
 * Default settings
 */
export const DEFAULT_SETTINGS: Dictate2MeSettings = {
	apiUrl: 'http://localhost:8765/api/v1',
	apiToken: '',
	language: 'pt',
	enableCorrection: true,
	showPartialResults: true,
	showConfidence: true,
	autoCheckDaemon: true,
};
