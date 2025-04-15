package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/4nkitd/sapiens"
)

type AiResponse struct {
	Response string `json:"response"`
	Command  string `json:"command"`
}

func AskQuery(query string, imageBytes [][]byte) AiResponse {
	apiKey := os.Getenv("GENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("GENAI_API_KEY environment variable is not set")
	}

	defaultModel := os.Getenv("GENAI_DEFAULT_MODEL")
	if defaultModel == "" {
		log.Fatal("GENAI_DEFAULT_MODEL environment variable is not set")
	}

	// Initialize the Sapiens LLM client
	llm := sapiens.NewGoogleGenAI(apiKey, defaultModel)

	if err := llm.Initialize(); err != nil {
		log.Fatalf("Failed to initialize LLM: %v", err)
	}

	// Create a Sapiens agent
	agent := sapiens.NewAgent("GemaCLI", llm, apiKey, defaultModel, "google")

	// Define system info tool
	sysInfoTool := sapiens.Tool{
		Name:        "get_system_info",
		Description: "Get detailed system information including hardware, OS, memory usage, and running processes",
	}

	// Add the system info tool to the agent
	agent.AddTools(sysInfoTool)
	// Register tool implementation
	agent.RegisterToolImplementation("get_system_info", func(params map[string]interface{}) (interface{}, error) {
		return GetSystemInfo(params)
	})

	// Add system prompt (without system info directly embedded)
	agent.AddSystemPrompt(SystemInstruction, "1.0")

	// Define schema for structured output
	schema := sapiens.Schema{
		Type: "object",
		Properties: map[string]sapiens.Schema{
			"response": {
				Type:        "string",
				Description: "The response from the AI",
			},
			"command": {
				Type:        "string",
				Description: "Command to execute on the system",
			},
		},
		Required: []string{"response"},
	}

	// Set the schema on your agent
	agent.SetStructuredResponseSchema(schema)

	// Create context and run the query
	ctx := context.Background()

	// Handle image attachments if present
	if len(imageBytes) > 0 {
		for _, imgBytes := range imageBytes {
			if imgBytes != nil {
				agent.AddImageContent(imgBytes, "image/jpeg")
			}
		}
	}

	// Run the agent with the query
	response, err := agent.Run(ctx, query)
	if err != nil {
		log.Fatalf("Error from agent: %v", err)
	}

	// fmt.Println("Response:", response.Content)

	// Create a result object with default values
	result := AiResponse{
		Response: "",
		Command:  "",
	}

	// Access the structured data from the response
	if response.Structured != nil {
		structuredData, ok := response.Structured.(map[string]interface{})
		if ok {
			// Extract response field
			if responseVal, exists := structuredData["response"]; exists {
				if responseStr, ok := responseVal.(string); ok {
					result.Response = responseStr
				}
			}

			// Extract command field
			if commandVal, exists := structuredData["command"]; exists {
				if commandStr, ok := commandVal.(string); ok {
					result.Command = commandStr
				}
			}
		} else {
			log.Printf("Failed to cast structured data to map[string]interface{}")
		}
	} else if response.Content != "" {
		// As fallback, try to parse the content as JSON
		var jsonResult AiResponse
		if json.Valid([]byte(response.Content)) {
			if err := json.Unmarshal([]byte(response.Content), &jsonResult); err == nil {
				result = jsonResult
			}
		} else {
			result.Response = response.Content
		}
	}

	// Validate that we have at least a response
	if result.Response == "" {
		log.Fatal("Response field is missing or not a string")
	}

	errDb := StoreCommandHistory(query, result.Response)
	if errDb != nil {
		log.Fatal(errDb)
	}

	return result
}
