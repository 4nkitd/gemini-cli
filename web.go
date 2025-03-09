package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "Start a web server with chat interface and API",
	Long: `Starts a web server that provides both a web interface and REST API for AI interactions.
The server runs on port 8080 by default (configurable via PORT environment variable).
Available endpoints:
- Web UI: http://localhost:8080/
- API: POST to http://localhost:8080/answer with JSON body {"message": "your question", "history": {}}`,
	Run: executeWebCommand,
}

func executeWebCommand(cmd *cobra.Command, args []string) {
	port := getPort()
	r := mux.NewRouter()

	// Open browser after a short delay to ensure server is running
	go func() {
		time.Sleep(500 * time.Millisecond)
		url := fmt.Sprintf("http://localhost:%d", port)
		log.Printf("Opening %s in your browser", url)
		if err := exec.Command("open", url).Run(); err != nil {
			// Try with xdg-open for Linux
			if err := exec.Command("xdg-open", url).Run(); err != nil {
				// Try with start for Windows
				if err := exec.Command("cmd", "/c", "start", url).Run(); err != nil {
					log.Printf("Could not open browser: %v", err)
				}
			}
		}
	}()

	r.HandleFunc("/answer", answerHandler).Methods("POST")

	webFS := getEmbeddedWebFS()
	r.PathPrefix("/").Handler(http.FileServer(http.FS(webFS)))

	log.Printf("Starting web server on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

func answerHandler(w http.ResponseWriter, r *http.Request) {

	var requestBody struct {
		Message string            `json:"message"`
		History map[string]string `json:"history"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var formattedHistory string
	if len(requestBody.History) > 0 {
		formattedHistory = "Previous conversation:\n"
		for question, answer := range requestBody.History {
			formattedHistory += "Question: " + question + "\nAnswer: " + answer + "\n\n"
		}
	}

	query := fmt.Sprintf("You are an AI assistant. %sNew question: %s",
		formattedHistory,
		requestBody.Message)

	os.Setenv("SKIP_SYS_INFO", "true")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	os.Setenv("SKIP_SYS_INFO", "true")
	ai := AskQuery(query, nil)
	os.Setenv("SKIP_SYS_INFO", "")

	response := `{"message": " ` + ai.Response + `"}`
	w.Write([]byte(response))
}

func getPort() int {
	portStr := os.Getenv("GENAI_PORT")
	if portStr == "" {
		color.New(color.FgYellow).Println("PORT environment variable not set, defaulting to 8080")
		return 8080
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		color.New(color.FgRed).Printf("Invalid PORT environment variable: %v, defaulting to 8080\n", err)
		return 8080
	}
	return port
}
