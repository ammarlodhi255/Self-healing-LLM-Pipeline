package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const ollamaEndpoint = "http://localhost:11434/api/generate" // The local endpoint for the Ollama API

// Struct for request to Ollama API
type OllamaRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

// Struct for response from Ollama API
type OllamaResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason,omitempty"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var conversationHistory strings.Builder // To keep track of the conversation

	for {
		fmt.Print("Enter your prompt (or type 'exit' to quit): ")
		userPrompt, _ := reader.ReadString('\n')
		userPrompt = strings.TrimSpace(userPrompt)

		if userPrompt == "exit" {
			fmt.Println("Exiting the program.")
			break
		}

		suffix := "The code should be in the Rust programming language. There should also be 3 robust test cases within the same code. Please only provide the source code and no further explanation. Thank you."
		userPrompt = userPrompt + " " + suffix

		// Append user prompt to conversation history
		conversationHistory.WriteString("User: " + userPrompt + "\n")

		// Generate response using Ollama API
		response, err := getOllamaResponse(conversationHistory.String())
		if err != nil {
			fmt.Println("Error generating response:", err)
			continue
		}

		// Print and update conversation history with model response
		fmt.Println("Ollama's response:", response)
		conversationHistory.WriteString("Ollama: " + response + "\n")

		var success, output = runCode(response)
		var newQuery string

		if success {
			fmt.Println(output)
		} else {
			newQuery = "Following are the errors, please fix the code. Write it again, and write only source code along with same test cases. \n" + output
			fmt.Println(newQuery)
		}

	}
}

// Function to make a POST request to Ollama API
func getOllamaResponse(conversation string) (string, error) {
	// Create request payload with the model specified
	requestBody, err := json.Marshal(OllamaRequest{
		Prompt: conversation,
		Model:  "llama3.1", // Use your downloaded model
	})
	if err != nil {
		return "", err
	}

	// Send HTTP POST request to Ollama API
	resp, err := http.Post(ollamaEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and accumulate response body in chunks
	var completeResponse string
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var chunk OllamaResponse
		if err := decoder.Decode(&chunk); err != nil {
			return "", err
		}
		completeResponse += chunk.Response

		// If the response is complete, break the loop
		if chunk.Done {
			break
		}
	}

	return completeResponse, nil
}
