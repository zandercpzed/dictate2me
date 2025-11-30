package audio

import (
	"math"
	"time"
)

// VAD (Voice Activity Detection) detects speech in audio streams
type VAD struct {
	sampleRate      int
	energyThreshold float64

	// Hysteresis state
	speechDuration  time.Duration
	silenceDuration time.Duration
	minSpeechDur    time.Duration
	minSilenceDur   time.Duration

	isSpeech bool
}

// NewVAD creates a new VAD instance
func NewVAD(sampleRate int) *VAD {
	return &VAD{
		sampleRate:      sampleRate,
		energyThreshold: 500.0, // Conservative default for 16-bit PCM
		minSpeechDur:    100 * time.Millisecond,
		minSilenceDur:   500 * time.Millisecond,
	}
}

// SetThreshold sets the energy threshold for speech detection
// Typical values for 16-bit PCM: 100-1000 for quiet environments
func (v *VAD) SetThreshold(threshold float64) {
	v.energyThreshold = threshold
}

// Process analyzes an audio chunk and returns true if speech is detected
// It maintains internal state for hysteresis
func (v *VAD) Process(chunk []int16) bool {
	if len(chunk) == 0 {
		return v.isSpeech
	}

	// Calculate RMS energy
	var sumSquares float64
	for _, sample := range chunk {
		val := float64(sample)
		sumSquares += val * val
	}
	rms := math.Sqrt(sumSquares / float64(len(chunk)))

	chunkDuration := time.Duration(float64(len(chunk)) / float64(v.sampleRate) * float64(time.Second))

	if rms > v.energyThreshold {
		// Energy detected
		v.speechDuration += chunkDuration
		v.silenceDuration = 0

		if v.speechDuration >= v.minSpeechDur {
			v.isSpeech = true
		}
	} else {
		// Silence detected
		v.silenceDuration += chunkDuration
		v.speechDuration = 0

		if v.silenceDuration >= v.minSilenceDur {
			v.isSpeech = false
		}
	}

	return v.isSpeech
}

// Reset clears the VAD state
func (v *VAD) Reset() {
	v.speechDuration = 0
	v.silenceDuration = 0
	v.isSpeech = false
}
