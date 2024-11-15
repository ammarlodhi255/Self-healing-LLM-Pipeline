package ollamaimplementation

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var OllamaEndpoint = "http://localhost:11434/api/generate" // The local endpoint for the Ollama API
// var OllamaEndpoint = "http://host.docker.internal:11434" // The local endpoint for the Ollama API

// Struct for request to Ollama API
type OllamaRequest struct {
	Prompt  string `json:"prompt"`
	Model   string `json:"model"`
	Context []int  `json:"context,omitempty"` // Context to maintain conversation
}

// Struct for response from Ollama API
type OllamaResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason,omitempty"`
	Context            []int  `json:"context,omitempty"` // Updated context
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
}

func GetOllamaResponse(prompt string, context []int, model string) (string, []int, error) {
	// Create request payload with the model specified and context
	requestBody, err := json.Marshal(OllamaRequest{
		Prompt:  prompt,
		// Model:   "llama3.1",
		Model: model,
		Context: context, // Pass the conversation context
	})
	if err != nil {
		return "", nil, err
	}

	// Send HTTP POST request to Ollama API
	resp, err := http.Post(OllamaEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// Read and accumulate response body in chunks
	var completeResponse string
	var updatedContext []int
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var chunk OllamaResponse
		if err := decoder.Decode(&chunk); err != nil {
			return "", nil, err
		}
		completeResponse += chunk.Response

		// Capture the updated context from the response
		updatedContext = chunk.Context

		if chunk.Done {
			break
		}
	}

	return completeResponse, updatedContext, nil
}
