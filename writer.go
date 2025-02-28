package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// WriterCmd represents the writer command
var WriterCmd = &cobra.Command{
	Use:     "writer [text]",
	Aliases: []string{"revise", "edit", "improve", "refine", "w"},
	Short:   "Revises text to be more professional using Gemini AI",
	Long: `A command that uses the Gemini API to revise input text and make it more professional.
It maintains the original length unless instructed otherwise with [length=X] or [type=email] tags.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selectedText := args[0]

		// Indicate processing
		processingMsg := color.New(color.FgYellow).PrintFunc()
		processingMsg("Processing your text with Gemini AI...\n")

		extractedText, err := extractGeminiText(selectedText)
		if err != nil {
			errorMsg := color.New(color.FgRed, color.Bold).PrintFunc()
			errorMsg("Error: " + err.Error())
			return err
		}

		// Print success indicator and the resulting text
		successMsg := color.New(color.FgGreen, color.Bold).PrintFunc()
		successMsg("✓ Professional revision complete:\n\n")
		fmt.Println(extractedText)
		return nil
	},
}

func extractGeminiText(selectedText string) (string, error) {
	// Rest of the function remains the same...
	systemPrompt := fmt.Sprintf(`
	Please revise the following text to be more professional. Maintain the original length as closely as possible, unless a length constraint or output type is specified.

	**Input Text:**
	%s

	**Instructions:**

	-   **Default Behavior:** If no explicit instructions are provided, maintain the original length of the input text (within a reasonable tolerance, a few words longer or shorter is acceptable) while making it sound more professional.
	-   **Length Constraint:** If [length=X] is present (where X is a number), adjust the output to be approximately X words long.
	-   **Output Type:** If [type=email] is present, format the output as a professional email, including a subject line, a greeting, and a closing.
	-   **Professionalism:**  Focus on clarity, concise wording, proper grammar, and a tone that is appropriate for business or formal communication. Avoid slang, colloquialisms, and overly informal language.

	**Example Usage:**

	-   **Basic:** Input Text: "Hey, wanna chat later?"  (Output would be something like: "I was wondering if you would be available to connect later?")
	-   **With Length Constraint:** Input Text: "This is a really long message that talks about a lot of stuff, so you should probably read it carefully. Blah blah blah blah blah." [length=50]
	(Output would be a refined, professional version of the input, limited to around 50 words.)
	-   **Email Format:** Input Text: "Hi, can we talk about this important thing? Thx." [type=email] (Output would be a complete email with a subject, greeting, closing)

	Please provide the improved text according to the given instructions.
	`, selectedText)

	// Use the AskQuery function from gemini.go
	response := AskQuery(systemPrompt)

	// Return the response text, trimming any whitespace
	return strings.TrimSpace(response.Response), nil
}
