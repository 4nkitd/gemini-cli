package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// WriterCmd represents the writer command
var MakeCmd = &cobra.Command{
	Use:     "cli [text]",
	Aliases: []string{"ask"},
	Short:   "ask questions about cli tool and other things",
	Long:    `ask questions about cli tool and other things.`,
	Args:    cobra.ExactArgs(1),
	RunE:    executeMakeCommand,
}

func executeMakeCommand(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	m := waitForResponse(model{query: query, loading: true})

	fmt.Println(m.View())

	if m.command != "" {
		var input string
		color.New(color.FgYellow).Print("Run command (y for yes, n for no): ")
		fmt.Scanln(&input)
		if input == "y" {
			runCommand(m.command)
		}
	}

	return nil
}

func waitForResponse(m model) model {
	// Show animated loading dots
	done := make(chan bool)
	go func() {
		loadingChars := []string{"|", "/", "-", "\\"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				color.New(color.FgCyan).Printf("\rLoading %s", loadingChars[i%len(loadingChars)])
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()

	time.Sleep(3 * time.Second)
	genaiResponse := AskQuery(m.query)
	done <- true
	fmt.Print("\r") // Clear the loading line

	m.loading = false
	m.response = genaiResponse.Response
	m.command = genaiResponse.Command
	return m
}

func (m model) View() string {
	if m.loading {
		return color.CyanString("Loading...")
	}

	responseHeader := color.New(color.FgGreen, color.Bold).Sprint("\nResponse:")
	commandHeader := color.New(color.FgYellow, color.Bold).Sprint("\nSuggested Command to RUN: ")
	formattedResponse := formatResponse(m.response)
	commandText := color.New(color.FgHiYellow).Sprint(m.command)

	return fmt.Sprintf("%s\n%s\n%s%s\n", responseHeader, formattedResponse, commandHeader, commandText)
}

func runCommand(command string) {
	color.New(color.FgBlue).Printf("Executing: %s\n", command)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		color.New(color.FgRed).Printf("Error executing command: %v\n", err)
		return
	}
	color.New(color.FgGreen).Println("Output:")
	fmt.Printf("=> %s\n", strings.ReplaceAll(string(out), "dump.rdb", ""))
}

func formatResponse(response string) string {
	words := strings.Fields(response)
	var formattedResponse strings.Builder

	for i, word := range words {
		formattedResponse.WriteString(word + " ")
		if (i+1)%15 == 0 {
			formattedResponse.WriteString("\n")
		}
	}

	return formattedResponse.String()
}
