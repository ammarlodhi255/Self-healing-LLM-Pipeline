package ollamaimplementation

import (
	"encoding/json"
	"io/ioutil"
	"llama/modules/extraction"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock OllamaResponse for testing
var mockResponse = OllamaResponse{
	Model:    "llama3.1",
	Response: "This is a mock response.",
	Done:     true,
	Context:  []int{1, 2, 3},
}

var model = "llama3.1"

// Mock function to replace the actual HTTP request
func TestGetOllamaResponse(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Read and verify the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		// Verify the request payload contains the expected prompt and model
		if !strings.Contains(string(body), `"prompt":"Test prompt"`) || !strings.Contains(string(body), `"model":"llama3.1"`) {
			t.Errorf("Request body is not as expected: %s", string(body))
		}

		// Send a mock response
		responseData, _ := json.Marshal(mockResponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	}))
	defer mockServer.Close()

	// Temporarily replace the ollamaEndpoint with the mock server's URL
	originalEndpoint := OllamaEndpoint
	OllamaEndpoint = mockServer.URL
	defer func() { OllamaEndpoint = originalEndpoint }() // Restore the original endpoint after the test

	// Call the function to be tested
	prompt := "Test prompt"
	context := []int{}
	response, updatedContext, err := GetOllamaResponse(prompt, context, model)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the response content
	expectedResponse := "This is a mock response."
	if response != expectedResponse {
		t.Errorf("Expected response %q, got %q", expectedResponse, response)
	}

	// Verify the updated context
	expectedContext := []int{1, 2, 3}
	for i, val := range updatedContext {
		if val != expectedContext[i] {
			t.Errorf("Expected context %v, got %v", expectedContext, updatedContext)
		}
	}
}

// Test for prompts.
var promptTestCases = []struct {
	name          string
	prompt        string
	suffixStr     string
	shouldContain []string
}{
	{"5 Even Integers GO", "Write a program that generates 5 random integers.", extraction.GoPrompt, []string{"```go", "```"}},
	{"Sort the array using mergesort GO", "Write a program that sorts the array [23, 2, 0, -1, 89, 500] using mergesort.", extraction.GoPrompt, []string{"```go", "```"}},
	{"Reverse the string GO.", "Reverse the string 'ammar'", extraction.GoPrompt, []string{"```go", "```"}},

	{"5 Even Integers rust", "Write a program that generates 5 random integers.", extraction.RustPrompt, []string{"```rust", "```"}},
	{"Sort the array using mergesort rust", "Write a program that sorts the array [23, 2, 0, -1, 89, 500] using mergesort.", extraction.RustPrompt, []string{"```rust", "```"}},
	{"Reverse the string rust.", "Reverse the string 'ammar'", extraction.RustPrompt, []string{"```rust", "```"}},
}

// Test function to verify the prefix and suffix of responses from GetOllamaResponse
func TestGetOllamaResponsePrompts(t *testing.T) {
	for _, tc := range promptTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function to get the response
			response, _, _ := GetOllamaResponse(tc.prompt+tc.suffixStr, []int{}, model)

			// Check if the response starts with the expected prefix
			if strings.HasPrefix(response, tc.shouldContain[0]) {
				// Check if the response ends with "```"
				if !strings.HasSuffix(response, "```") {
					t.Errorf("Test %s failed: expected response to end with ```; got %q", tc.name, response)
				}
			} else {
				t.Errorf("Test %s failed: expected response to start with %q; got %q", tc.name, tc.shouldContain[0], response)
			}
		})
	}
}
