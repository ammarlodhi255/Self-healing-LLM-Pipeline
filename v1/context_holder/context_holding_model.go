package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func runCode(rustCode string) (bool, string) {
	projectDir := "rust_project"
	os.Mkdir(projectDir, 0755)

	// Create the Cargo.toml file with necessary dependencies
	cargoToml := `
	[package]
	name = "rust_project"
	version = "0.1.0"
	edition = "2021"

	[dependencies]
	rand = "0.8"
	`
	cargoTomlPath := projectDir + "/Cargo.toml"
	if err := os.WriteFile(cargoTomlPath, []byte(cargoToml), 0644); err != nil {
		fmt.Println("Error writing Cargo.toml:", err)
		return false, ""
	}

	// Create the src directory and main.rs file
	os.Mkdir(projectDir+"/src", 0755)
	mainRsPath := projectDir + "/src/main.rs"
	if err := os.WriteFile(mainRsPath, []byte(rustCode), 0644); err != nil {
		fmt.Println("Error writing Rust code to file:", err)
		return false, ""
	}

	// Compile the Rust project using Cargo
	cmd := exec.Command("cargo", "build", "--release")
	cmd.Dir = projectDir
	var out, stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("Compilation failed with errors:")
		return false, stderr.String()
	}

	// If compilation succeeded, run the binary and capture output
	cmd = exec.Command("./target/release/rust_project")
	cmd.Dir = projectDir
	var result strings.Builder
	cmd.Stdout = &result
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println("Error running the compiled binary:")
		return false, stderr.String()
	}

	// Output from the Rust program
	// fmt.Println("Rust program output:")
	// fmt.Println(result.String())

	// Clean up: Remove the rust_project directory
	if err := os.RemoveAll(projectDir); err != nil {
		fmt.Println("Error cleaning up the project directory:", err)
	}

	return true, result.String()
}

const ollamaEndpoint = "http://localhost:11434/api/generate" // The local endpoint for the Ollama API

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

func main() {
	reader := bufio.NewReader(os.Stdin)
	var conversationContext []int // Variable to store conversation context

	fmt.Print("Enter your prompt (or type 'exit' to quit): ")
	userPrompt, _ := reader.ReadString('\n')
	userPrompt = strings.TrimSpace(userPrompt)
	var firstIter bool = true

	for {

		if userPrompt == "exit" {
			fmt.Println("Exiting the program.")
			break
		} else if firstIter {
			suffix := "The code should be in the Rust programming language. There should also be 3 robust test cases within the same code. There should also be a main function inside of which all the execution takes place. Please only provide the source code and no further explanation, The format should be ```rust <yourcode + testcases> ```"
			userPrompt = userPrompt + " " + suffix
			firstIter = false
		}

		// Generate response using Ollama API, passing the context
		response, updatedContext, err := getOllamaResponse(userPrompt, conversationContext)
		if err != nil {
			fmt.Println("Error generating response:", err)
			continue
		}

		// Update the conversation context with the response
		conversationContext = updatedContext

		// fmt.Println("Ollama's response:", response)

		var code string = strings.ReplaceAll(response, "```rust", "")
		code = strings.ReplaceAll(code, "```", "")

		fmt.Println("Code:", code)

		var success, output = runCode(code)

		if success {
			fmt.Println("Output: ", output)
			userPrompt = "exit"
		} else {
			fmt.Println(success, output)
			userPrompt = output + "\nFollowing are the errors, please fix the code. Write it again, and write only source code along with same test cases with no further explanation. The format should be ```rust <yourcode + testcases> ```"
		}
	}
}

// Function to make a POST request to Ollama API
func getOllamaResponse(prompt string, context []int) (string, []int, error) {
	// Create request payload with the model specified and context
	requestBody, err := json.Marshal(OllamaRequest{
		Prompt:  prompt,
		Model:   "llama3.1", // Use your downloaded model
		Context: context,    // Pass the conversation context
	})
	if err != nil {
		return "", nil, err
	}

	// Send HTTP POST request to Ollama API
	resp, err := http.Post(ollamaEndpoint, "application/json", bytes.NewBuffer(requestBody))
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

		// If the response is complete, break the loop
		if chunk.Done {
			break
		}
	}

	return completeResponse, updatedContext, nil
}
