package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kbinani/screenshot"
)

func Screenshot() ([]byte, error) {

	n := 0
	if screenshot.NumActiveDisplays() != 0 {
		n = 0
	}
	bounds := screenshot.GetDisplayBounds(n)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		return nil, err
	}

	imageBytes := buffer.Bytes()

	return imageBytes, nil

}

func SpeakMessage(message string, stopSpeech chan struct{}) {
	// Use appropriate text-to-speech command based on OS
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("say", message)
	case "windows":
		cmd = exec.Command("powershell", "-c", "Add-Type -AssemblyName System.Speech; $speak = New-Object System.Speech.Synthesis.SpeechSynthesizer; $speak.Speak('"+message+"')")
	case "linux":
		cmd = exec.Command("spd-say", message)
	default:
		fmt.Println("Text-to-speech not supported on this platform")
		return
	}

	// Start command asynchronously
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting speech:", err)
		return
	}

	// Create a channel to signal command completion
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either command completion or stop signal
	select {
	case err := <-done:
		if err != nil {
			fmt.Println("Error speaking message:", err)
		}
	case <-stopSpeech:
		// Kill the process if stop signal received
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Failed to stop speech:", err)
		}
	}
}

func PutTextOnClipboard(data string) error {

	// Use appropriate clipboard command based on OS
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "c")
	default:
		fmt.Println("Clipboard not supported on this platform")
		return nil
	}

	cmd.Stdin = strings.NewReader(data)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil

}

func RecordAudio(seconds int) ([]byte, error) {

	audioBytes := []byte{}

	return audioBytes, nil

}
