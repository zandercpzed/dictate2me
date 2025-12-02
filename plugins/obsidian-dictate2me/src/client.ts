/**
 * Client for Dictate2Me API
 * Handles WebSocket connection and audio streaming
 */

type EventHandler = (data: any) => void;

interface StreamConfig {
	language: string;
	enableCorrection: boolean;
}

interface StreamMessage {
	type: string;
	data?: any;
}

export class Dictate2MeClient {
	private apiUrl: string;
	private token: string;
	private ws: WebSocket | null = null;
	private mediaStream: MediaStream | null = null;
	private audioContext: AudioContext | null = null;
	private scriptProcessor: ScriptProcessorNode | null = null;
	private eventHandlers: Map<string, EventHandler[]> = new Map();

	constructor(apiUrl: string, token: string) {
		this.apiUrl = apiUrl;
		this.token = token;
	}

	/**
	 * Register event handler
	 */
	on(event: string, handler: EventHandler) {
		if (!this.eventHandlers.has(event)) {
			this.eventHandlers.set(event, []);
		}
		this.eventHandlers.get(event)!.push(handler);
	}

	/**
	 * Emit event to handlers
	 */
	private emit(event: string, data: any) {
		const handlers = this.eventHandlers.get(event);
		if (handlers) {
			handlers.forEach((handler) => handler(data));
		}
	}

	/**
	 * Connect to WebSocket and start streaming
	 */
	async connect(config: StreamConfig): Promise<void> {
		// Create WebSocket connection
		const wsUrl = this.apiUrl.replace('http://', 'ws://').replace('/api/v1', '/api/v1/stream');
		
		this.ws = new WebSocket(wsUrl);

		return new Promise((resolve, reject) => {
			if (!this.ws) {
				reject(new Error('WebSocket not initialized'));
				return;
			}

			this.ws.onopen = async () => {
				console.log('WebSocket connected');

				// Send start message
				this.sendMessage({
					type: 'start',
					data: config,
				});

				// Start audio capture
				try {
					await this.startAudioCapture();
					resolve();
				} catch (error) {
					reject(error);
				}
			};

			this.ws.onmessage = (event) => {
				try {
					const message: StreamMessage = JSON.parse(event.data);
					this.handleMessage(message);
				} catch (error) {
					console.error('Error parsing message:', error);
				}
			};

			this.ws.onerror = (error) => {
				console.error('WebSocket error:', error);
				this.emit('error', 'Connection error');
				reject(error);
			};

			this.ws.onclose = () => {
				console.log('WebSocket closed');
				this.cleanup();
			};
		});
	}

	/**
	 * Disconnect and stop streaming
	 */
	async disconnect(): Promise<void> {
		// Send stop message
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.sendMessage({ type: 'stop' });
		}

		// Cleanup
		this.cleanup();
	}

	/**
	 * Start capturing audio from microphone
	 */
	private async startAudioCapture(): Promise<void> {
		try {
			// Request microphone access
			this.mediaStream = await navigator.mediaDevices.getUserMedia({
				audio: {
					channelCount: 1,
					sampleRate: 16000,
					echoCancellation: true,
					noiseSuppression: true,
					autoGainControl: true,
				},
			});

			// Create audio context
			this.audioContext = new AudioContext({ sampleRate: 16000 });
			const source = this.audioContext.createMediaStreamSource(this.mediaStream);

			// Create script processor for audio data
			const bufferSize = 4096;
			this.scriptProcessor = this.audioContext.createScriptProcessor(bufferSize, 1, 1);

			this.scriptProcessor.onaudioprocess = (event) => {
				const inputData = event.inputBuffer.getChannelData(0);
				this.sendAudioChunk(inputData);
			};

			// Connect nodes
			source.connect(this.scriptProcessor);
			this.scriptProcessor.connect(this.audioContext.destination);

			console.log('Audio capture started');
		} catch (error) {
			console.error('Error starting audio capture:', error);
			throw new Error('Failed to access microphone');
		}
	}

	/**
	 * Send audio chunk to server
	 */
	private sendAudioChunk(audioData: Float32Array) {
		if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;

		// Convert Float32Array to Int16Array (16-bit PCM)
		const int16Data = new Int16Array(audioData.length);
		for (let i = 0; i < audioData.length; i++) {
			const sample = Math.max(-1, Math.min(1, audioData[i]));
			int16Data[i] = sample < 0 ? sample * 0x8000 : sample * 0x7fff;
		}

		// Convert to base64
		const bytes = new Uint8Array(int16Data.buffer);
		const base64 = btoa(String.fromCharCode(...bytes));

		// Send message
		this.sendMessage({
			type: 'audio',
			data: { data: base64 },
		});
	}

	/**
	 * Send message through WebSocket
	 */
	private sendMessage(message: StreamMessage) {
		if (this.ws && this.ws.readyState === WebSocket.OPEN) {
			this.ws.send(JSON.stringify(message));
		}
	}

	/**
	 * Handle incoming message from server
	 */
	private handleMessage(message: StreamMessage) {
		switch (message.type) {
			case 'partial':
				this.emit('partial', message.data?.text || '');
				break;

			case 'final':
				this.emit('final', message.data);
				break;

			case 'error':
				this.emit('error', message.data?.message || 'Unknown error');
				break;

			default:
				console.warn('Unknown message type:', message.type);
		}
	}

	/**
	 * Cleanup resources
	 */
	private cleanup() {
		// Stop audio processing
		if (this.scriptProcessor) {
			this.scriptProcessor.disconnect();
			this.scriptProcessor = null;
		}

		// Close audio context
		if (this.audioContext) {
			this.audioContext.close();
			this.audioContext = null;
		}

		// Stop media stream
		if (this.mediaStream) {
			this.mediaStream.getTracks().forEach((track) => track.stop());
			this.mediaStream = null;
		}

		// Close WebSocket
		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
	}
}
