package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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
	var aiTemp float32 = 1

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel(defaultModel)
	model.Temperature = &aiTemp

	sysInfo := ""
	sysErr := error(nil)
	if sysInfo, sysErr = GetSystemInfo(); err != nil {
		log.Fatal(sysErr)
	}

	finalPrompt := sysInfo + SystemInstruction

	SystemInstruction := []genai.Part{
		genai.Text(finalPrompt),
	}

	attachMents := []genai.Part{}

	attachMents = append(attachMents, genai.Text(query))

	for _, imgBytes := range imageBytes {
		if imageBytes != nil {

			attachMents = append(attachMents, genai.ImageData("image/jpeg", imgBytes))
		}
	}

	model.SystemInstruction = &genai.Content{Parts: SystemInstruction}

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"response": {
				Type: genai.TypeString,
			},
			"command": {
				Type: genai.TypeString,
			},
		},
		Required: []string{"response"},
	}

	resp, err := model.GenerateContent(ctx, attachMents...)
	if err != nil {
		log.Fatal(err)
	}

	jsonResponse := responseString(resp)

	if !json.Valid([]byte(jsonResponse)) {
		log.Fatal("Invalid JSON response")
	}

	var result AiResponse
	if err := json.Unmarshal([]byte(jsonResponse), &result); err != nil {
		log.Fatal(err)
	}

	if result.Response == "" {
		log.Fatal("Response field is missing or not a string")
	}

	StoreCommandHistory(query, result.Response)

	return result
}

func responseString(resp *genai.GenerateContentResponse) string {
	var b strings.Builder
	for i, cand := range resp.Candidates {
		if len(resp.Candidates) > 1 {
			fmt.Fprintf(&b, "%d:", i+1)
		}
		b.WriteString(contentString(cand.Content))
	}
	return b.String()
}

func contentString(c *genai.Content) string {
	var b strings.Builder
	if c == nil || c.Parts == nil {
		return ""
	}
	for i, part := range c.Parts {
		if i > 0 {
			fmt.Fprintf(&b, ";")
		}
		fmt.Fprintf(&b, "%v", part)
	}
	return b.String()
}
