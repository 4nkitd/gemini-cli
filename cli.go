package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// WriterCmd represents the writer command
var MakeCmd = &cobra.Command{
	Use:   "cli [text]",
	Short: "ask questions about cli tool and other things",
	Long:  `ask questions about cli tool and other things.`,
	Args:  cobra.ExactArgs(1),
	RunE:  executeMakeCommand,
}

func executeMakeCommand(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	m := waitForResponse(model{query: query, loading: true})

	fmt.Println(m.View())

	if m.command != "" {
		var input string
		fmt.Print("Run command (y for yes, n for no): ")
		fmt.Scanln(&input)
		if input == "y" {
			runCommand(m.command)
		}
	}

	return nil
}

func waitForResponse(m model) model {
	time.Sleep(3 * time.Second)
	genaiResponse := AskQuery(m.query)
	m.loading = false
	m.response = genaiResponse.Response
	m.command = genaiResponse.Command
	return m
}

func (m model) View() string {
	if m.loading {
		return "Loading..."
	}
	return fmt.Sprintf("\nResponse:\n%s\n\nSuggested Command to RUN: %s\n", formatResponse(m.response), m.command)
}

func runCommand(command string) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Printf("Error executing command: %v\n", err)
		return
	}
	fmt.Printf("Output:\n=> %s\n", strings.ReplaceAll(string(out), "dump.rdb", ""))
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
