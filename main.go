package main

import (
	"log"

	"github.com/spf13/cobra"
)

type model struct {
	query    string
	loading  bool
	response string
	command  string
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "ai",
		Short: "A CLI tool to execute commands",
	}

	rootCmd.AddCommand(MakeCmd)

	rootCmd.AddCommand(WriterCmd)

	rootCmd.AddCommand(GitCommitCmd)

	rootCmd.AddCommand(CoPilotCmd)

	rootCmd.AddCommand(WebCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}
