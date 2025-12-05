package main

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

func main() {
	if err := portaudio.Initialize(); err != nil {
		fmt.Printf("PortAudio Initialize error: %v\n", err)
		return
	}
	defer portaudio.Terminate()

	// List devices as a basic smoke test
	devices, err := portaudio.Devices()
	if err != nil {
		fmt.Printf("PortAudio Devices error: %v\n", err)
		return
	}
	fmt.Printf("PortAudio initialized, found %d devices\n", len(devices))
}
