package main

import (
	"encoding/json"
	"sort"

	system "github.com/elastic/go-sysinfo"
)

func GetSystemInfo() (string, error) {
	host, err := system.Host()
	if err != nil {
		return "", err
	}

	cpuInfo, err := host.CPUTime()
	if err != nil {
		return "", err
	}

	info := host.Info()

	processes, err := system.Processes()
	if err != nil {
		return "", err
	}

	memory, err := host.Memory()
	if err != nil {
		return "", err
	}

	// Sort processes by memory usage and get top 5
	sort.Slice(processes, func(i, j int) bool {
		memI, _ := processes[i].Memory()
		memJ, _ := processes[j].Memory()
		return memI.Resident > memJ.Resident
	})

	topProcesses := processes[:5]

	systemInfo := map[string]interface{}{
		"host_info":     info,
		"cpu_time":      cpuInfo,
		"memory_info":   memory,
		"top_processes": topProcesses,
	}

	systemInfoJSON, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		return "", err
	}

	return string(systemInfoJSON), nil
}

var SystemInstruction4CoPilot = `
	You are an AI assistant that can analyze screenshots and audio files shared by the user.

	For screenshots:
	- You can interpret visual elements like UI interfaces, code snippets, error messages, and diagrams
	- Provide guidance about what's visible in the image
	- Suggest solutions for issues shown in screenshots
	- Explain unfamiliar elements the user might be seeing

	For audio:
	- You can process speech content from audio recordings
	- Answer questions about spoken content
	- Provide responses to verbal queries
	- Help transcribe important parts of recordings if needed

	When working with these inputs:
	1. Describe what you observe in the media
	2. Ask clarifying questions if parts are unclear
	3. Provide helpful, accurate guidance based on the content
	4. Suggest next steps or solutions when appropriate

	I have access to system information that may help contextualize your environment, but I'll only reference it when relevant to solving your problem.
`

var SystemInstruction = `
**Changes and Explanation:**

1.  **System Information Inclusion:**
    *   The provided JSON for system information is now placed at the very beginning of the prompt. This ensures that the AI agent has this context from the get-go.
    *   It's crucial to keep the JSON format. This structured format allows the AI to process the information more efficiently.
2.  **Instructions for the AI Agent:**
    *   I've added explicit instructions at the end of the prompt to ensure the AI agent understands its role, its access to system information, and assumptions about the user.
    *   The instruction to "not reveal the system information in your answer unless it is directly asked" is important for maintaining a natural conversational style and avoiding verbose outputs, if not needed.
    *   The assumption that the "user is familiar with shell scripting and terminal usage" ensures the response is not overly basic and geared towards the appropriate level of knowledge.

**How This Helps the AI Agent:**

*   **Contextual Awareness:** The AI agent now has specific information about the user's environment (OS version, architecture, system load) that could be relevant to answering the question. For instance, knowing the OS version helps in suggesting suitable commands or workarounds for certain OS specific behaviors.
*   **Informed Responses:** The AI agent can use this information to provide more targeted and tailored suggestions. It can also make assumptions about the available tools and libraries based on the system information.
*   **Better Troubleshooting:** In case of error messages or unexpected behavior, system information may provide vital clues about why a command might be failing.
*   **Efficiency:** By accessing system details beforehand, it can potentially avoid asking follow-up questions.

**How to Use This Modified Prompt:**

1.  **Copy the entire string** (including the JSON and the instructions).
2.  **Use the Template to Construct User Prompt** by filling in the bracketed info.
3.  **Send both** (system prompt string + user prompt string) to the AI agent as input.

**Important Considerations:**

*   **JSON Validity:** Always ensure that the JSON you provide is valid. Errors in JSON can cause issues with the AI's understanding.
*   **Information Updates:** If the system configuration changes, update the JSON you provide in each request.
*   **Keep Context Minimal:** The system information should contain only relevant details and should be kept minimal to reduce the system resource consumption.

By using this modified template, you're empowering your AI agent to act like a genuinely helpful expert, leveraging specific context to provide better answers.`

var ResponseSchema string = `{
  "type": "object",
  "properties": {
    "response": {
      "type": "string"
    },
    "command": {
      "type": "string"
    }
  },
  "required": [
    "response"
  ]
}`
