package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// GitCommitCmd represents the git commit command
var GitCommitCmd = &cobra.Command{
	Use:     "commit [path]",
	Aliases: []string{"c"},
	Short:   "Generate an AI commit message and optionally commit changes",
	Long:    "Generate a commit message using AI for uncommitted changes and optionally commit them",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		} else {
			// Use the directory where command is executed
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			path = currentDir
		}

		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("the path %s is not a valid directory", path)
		}

		if !IsGitRepo(path) {
			return fmt.Errorf("the path %s is not a git repository", path)
		}

		if !HasUncommittedChanges(path) {
			color.Yellow("There are no uncommitted changes in the repository at %s.", path)
			return nil
		}

		systemPrompt, _ := cmd.Flags().GetString("prompt")
		commitMessage, changedFiles := GenerateCommitMessage(path, systemPrompt)

		if len(changedFiles) == 0 {
			fmt.Println(commitMessage)
			return fmt.Errorf("failed to generate commit message")
		}

		color.Green("Commit message: %s", commitMessage)
		color.Blue("Files to be committed:")
		for _, file := range changedFiles {
			fmt.Printf("  - %s\n", file)
		}

		fmt.Print(color.CyanString("Do you want to commit these changes? (yes/no): "))
		var confirm string
		fmt.Scanln(&confirm)
		confirm = strings.ToLower(strings.TrimSpace(confirm))

		if confirm == "yes" {
			if err := CommitChanges(path, commitMessage); err != nil {
				return fmt.Errorf("error committing changes: %v", err)
			}
			color.Green("Changes have been committed.")
		} else {
			color.Yellow("Commit aborted.")
		}
		return nil
	},
}

func init() {
	GitCommitCmd.Flags().StringP("prompt", "p", "", "Custom system prompt for generating the commit message")
}

// IsGitRepo checks if the given path is a git repository
func IsGitRepo(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse")
	cmd.Stderr = nil
	cmd.Stdout = nil
	return cmd.Run() == nil
}

// HasUncommittedChanges checks if there are uncommitted changes in the git repository
func HasUncommittedChanges(path string) bool {
	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) != ""
}

// GenerateCommitMessage generates a commit message using git diff and the Gemini API
func GenerateCommitMessage(path, systemPrompt string) (string, []string) {
	if systemPrompt == "" {
		systemPrompt = "Generate a commit message in present tense and less than 50 words for the following changes:"
	}

	cmd := exec.Command("git", "-C", path, "diff", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return color.RedString("Error getting git diff: %v", err), nil
	}

	changedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(changedFiles) == 1 && changedFiles[0] == "" {
		return color.RedString("No changed files found"), nil
	}

	// Get changed files content for better context
	cmd = exec.Command("git", "-C", path, "diff")
	diffOutput, err := cmd.Output()
	if err != nil {
		return color.RedString("Error getting git diff content: %v", err), nil
	}

	// Prepare the prompt with file names and diff content
	query := fmt.Sprintf("%s\n\nChanged files:\n%s\n\nDiff:\n%s",
		systemPrompt,
		strings.Join(changedFiles, "\n"),
		limitDiffSize(string(diffOutput), 4000)) // Limit diff size to avoid token limits

	// Use AskQuery from gemini.go
	result := AskQuery(query, nil)
	return result.Response, changedFiles
}

// limitDiffSize limits the diff output size to avoid token limits
func limitDiffSize(diff string, maxSize int) string {
	if len(diff) <= maxSize {
		return diff
	}
	return diff[:maxSize] + "\n...(diff truncated due to size)"
}

// CommitChanges commits changes with the given commit message
func CommitChanges(path, commitMessage string) error {
	addCmd := exec.Command("git", "-C", path, "add", ".")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("error adding changes: %w", err)
	}

	commitCmd := exec.Command("git", "-C", path, "commit", "-m", commitMessage)
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("error committing changes: %w", err)
	}

	return nil
}
