package main

import (
	"context"
	"fmt"
	"llama/modules/compiler_v2/go_compiler_v2"
	displayindicator "llama/modules/display-indicator"
	"llama/modules/extraction"
	ollamaimplementation "llama/modules/ollama-implementation"

	// "log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

var model string = "llama3.1"
// var model string = "codellama:13b"
// var model string = "codellama"

// // MongoDB setup
// var client *mongo.Client
// var collection *mongo.Collection

// func initMongoDB() {
//     var err error
//     client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
//     if err != nil {
//         log.Fatal("Failed to connect to MongoDB:", err)
//     }
//     collection = client.Database("proj/local").Collection("startup_log")
// }

// func insertGeneratedCode(responseData ResponseData) {
//     _, err := collection.InsertOne(context.TODO(), bson.M{
//         "iteration":           responseData.Iteration,
//         "generated_code":      responseData.GeneratedCode,
//         "generated_code1":     responseData.GeneratedCode1,
//         "generated_code2":     responseData.GeneratedCode2,
//         "compiler_output":     responseData.CompilerOutput,
//         "compiled_successfully": responseData.CompiledSuccessfully,
//         "total_execution_time": responseData.TotalExecutionTime,
//     })
//     if err != nil {
//         fmt.Println("Failed to insert document:", err)
//     } else {
//         fmt.Println("Document inserted successfully")
//     }
// }

type compilerFunc func(string, string) ([]byte, error)

type ResponseData struct {
	Iteration           uint   `json:"iteration"`
	GeneratedCode       string `json:"generated_code"`
	GeneratedCode1      string `json:"generatedCode1"`  // Add this
	GeneratedCode2      string `json:"generatedCode2"`  // Add this
	CompilerOutput      string `json:"compiler_output"`
	CompiledSuccessfully bool   `json:"compiled_successfully"`
	TotalExecutionTime  string `json:"total_execution_time,omitempty"`
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func removeUnwantedLines(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	for _, line := range lines {
		// Skip lines starting with "go: " or "# command-line-arguments"
		if strings.HasPrefix(line, "go: ") || strings.HasPrefix(line, "# command-line-arguments") {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// Function to remove lines containing "go mod tidy" anywhere in the line
func removeGoModTidyLines(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	for _, line := range lines {
		if strings.Contains(line, "go mod tidy") {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

func removeLinesContaining(input string) string {
	keyword := "tempOutput"
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if !strings.Contains(line, keyword) {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func main() {
	// initMongoDB()  // Initialize MongoDB connection
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", handleWebSocket) // WebSocket route
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "index.html"))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Error upgrading to WebSocket:", err)
        return
    }
    defer conn.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var msg map[string]string
    if err := conn.ReadJSON(&msg); err != nil {
        fmt.Println("Error reading JSON:", err)
        return
    }
    userPrompt := msg["prompt"]

    // Update model based on user selection
    if selectedModel, ok := msg["model"]; ok {
        model = selectedModel
        fmt.Println("Model updated to:", model)
    }

    go RunProgram(ctx, conn, userPrompt, extraction.GoPrompt, go_compiler_v2.NewGoCompiler().CheckCompileErrors)

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            cancel()
            break
        }
    }
}

func RunProgram(ctx context.Context, conn *websocket.Conn, userPrompt string, languagePrompt string, compiler compilerFunc) {
	currentConversationContext := []int{}
	var numOfIterations uint = 1
	startTime := time.Now() // Track start time for execution duration

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Process canceled by client")
			return // Exit the loop if context is canceled
		default:
			// Prepare the prompt for LLM
			modifiedPrompt := userPrompt + languagePrompt
			fmt.Println("Prompt received. Generating response...")

			// Display loading indicator (for testing purposes)
			done := make(chan bool)
			go displayindicator.DisplayLoadingIndicator(done)

			// Generate LLM response
			fmt.Println("Model:", model)
			response, updatedContext, err := ollamaimplementation.GetOllamaResponse(modifiedPrompt, currentConversationContext, model)
			done <- true

			if err != nil {
				fmt.Println("Error generating response:", err)
				conn.WriteJSON(ResponseData{CompilerOutput: "Error generating response"})
				return
			}

			fmt.Println("LLM Response", response)
			currentConversationContext = updatedContext
			generatedCode1, generatedCode2, errExtract := extraction.Extract(response)

			fmt.Println("\n\nExtraction", generatedCode1)
			fmt.Println("\n\nExtraction", generatedCode2)

			// Handle extraction errors
			if errExtract != nil {
				fmt.Println("Improper LLM response")
				conn.WriteJSON(ResponseData{CompilerOutput: "Improper LLM response"})
				continue
			} 

			// Compile the generated code
			output, err := compiler(generatedCode1, generatedCode2)
			fmt.Println(err)
			compilationSuccess := err == nil

			// Prepare response data with generatedCode1 and generatedCode2
			responseData := ResponseData{
				Iteration:           numOfIterations,
				GeneratedCode:       response,
				GeneratedCode1:      generatedCode1,   // Include here
				GeneratedCode2:      generatedCode2,   // Include here
				CompilerOutput:      removeLinesContaining(removeGoModTidyLines(removeUnwantedLines(string(output)))),
				CompiledSuccessfully: compilationSuccess,
				TotalExecutionTime:  time.Since(startTime).String(),
			}

			// Insert the generated code into MongoDB
			// insertGeneratedCode(responseData)


			// Send iteration data back to the client
			conn.WriteJSON(responseData)

			// Check for successful compilation
			if compilationSuccess {
				return
			}

			// Update prompt with errors for next iteration
			userPrompt = "Please fix the error(s). If the error is not related to test cases, write the previous test cases again. The format should be ```go main code``` and ```go testcode```. Following is the error:\n" + (removeLinesContaining(removeGoModTidyLines(removeUnwantedLines(string(output)))))
			numOfIterations++
		}
	}
}
