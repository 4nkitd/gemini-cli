package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/ncruces/zenity"
	"github.com/spf13/cobra"
)

// WriterCmd represents the writer command
var CoPilotCmd = &cobra.Command{
	Use:     "assist [text]",
	Aliases: []string{"a", "copilot", "assistant"},
	Short:   "Help user with anything on his/her screen.",
	Long:    `Help user with anything on his/her screen.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Indicate processing
		processingMsg := color.New(color.FgYellow).PrintFunc()
		processingMsg("Starting Gema tray application...\n")

		systray.Run(onReady, onExit)

		return nil
	},
}

func onReady() {
	fmt.Println("Initializing system tray...")
	systray.SetIcon(icon.Data)
	systray.SetTitle("Gema")
	mHelp := systray.AddMenuItem("Help", "Take Screensot and and audio.")
	fmt.Println("Gema is now running in your system tray")

	// Listen for clicks on the help button
	go func() {
		for {
			<-mHelp.ClickedCh
			fmt.Println("Help button clicked, launching Copilot...")
			RunCopilot()
		}
	}()
}

func onExit() {
	fmt.Println("Exiting Gema...")
	// clean up here
}

func RunCopilot() {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("[%s] Starting Gema Assistant\n", time.Now().Format("15:04:05"))

	fmt.Println("[INFO] Displaying input dialog...")
	// Ask the user for their question with a centered dialog
	userQuery, err := zenity.Entry(
		"What would you like help with?",
		zenity.Title("Gema Assistant"),
		zenity.Width(400),
	)
	if err != nil {
		// User canceled or there was an error
		fmt.Printf("[ERROR] Dialog error: %v\n", err)
		return
	}

	// If the user didn't enter anything, return early
	if userQuery == "" {
		fmt.Println("[INFO] No query provided, canceling request")
		return
	}

	fmt.Println("[INFO] Taking screenshot...")
	imgBytes, err := Screenshot()
	if err != nil {
		fmt.Printf("[ERROR] Failed to capture screenshot: %v\n", err)
	} else {
		fmt.Printf("[INFO] Screenshot captured successfully (%d bytes)\n", len(imgBytes))
	}

	fmt.Println("[INFO] Sending query to AI service...")
	aiResp := AskQuery(userQuery, [][]byte{imgBytes})

	fmt.Println("====================================")
	fmt.Printf("[%s] QUERY: %s\n", time.Now().Format("15:04:05"), userQuery)
	fmt.Println("====================================")
	fmt.Println("[RESPONSE]")
	fmt.Println(aiResp.Response)
	fmt.Println("====================================")

	fmt.Println("[INFO] Copying response to clipboard...")
	err = PutTextOnClipboard(aiResp.Response)
	if err != nil {
		fmt.Printf("[ERROR] Failed to copy to clipboard: %v\n", err)
	} else {
		fmt.Println("[INFO] Response copied to clipboard")
	}

	fmt.Println("[INFO] Converting response to speech...")
	// Create a channel to signal speech to stop
	stopSpeech := make(chan struct{})

	// Start speech in a goroutine
	go func() {
		SpeakMessage(aiResp.Response, stopSpeech)
	}()

	// Monitor CLI input in parallel
	go func() {
		fmt.Println("[INFO] Press any key to stop speech...")
		var input string
		fmt.Scanln(&input) // Wait for Enter, but any input will work
		close(stopSpeech)  // Signal to stop the speech
		fmt.Println("[INFO] Speech stopped by user input")
	}()

	// Wait for speech to complete or be stopped
	<-stopSpeech

	fmt.Printf("[%s] Gema Assistant completed\n", time.Now().Format("15:04:05"))
	fmt.Println("--------------------------------------------------")
}
