package main

import (
	"bufio"
	"fmt"
	"llama/modules/compiler_v2/go_compiler_v2"
	displayindicator "llama/modules/display-indicator"
	"llama/modules/extraction"
	ollamaimplementation "llama/modules/ollama-implementation"
	"os"
	"strings"
	"time"
)

var startTime time.Time = time.Now()
var model string = "codellama"

type compilerFunc func([]byte, ...string) ([]byte, error)

func main() {

	reader := bufio.NewReader(os.Stdin)
	var conversationContext []int // Variable to store conversation context

	fmt.Print("Enter your prompt (or type 'exit' to quit): ")
	userPrompt, _ := reader.ReadString('\n')
	userPrompt = strings.TrimSpace(userPrompt)

	RunProgram(userPrompt, conversationContext, extraction.GoPrompt, go_compiler_v2.NewGoCompiler().CheckCompileErrors)

	// Go routines dont seem to be executing at the moment, so will figure this out later.
	// go RunProgram(userPrompt, conversationContext, extraction.GoPrompt, go_compiler_v2.NewGoCompiler().CheckCompileErrors)
	// go RunProgram(userPrompt, conversationContext, extraction.RustPrompt, rust_compiler_v2.NewRustCompiler().CheckCompileErrors)


	// Compute total execution time.
	endtime := time.Now()
    diff := endtime.Sub(startTime)
    fmt.Println("Total Execution Time:", diff)
}

func RunProgram(userPrompt string, conversationContext []int, languagePrompt string, compiler compilerFunc) {
	currentConversationContext := conversationContext
	var numOfIterations uint = 1

	for {
		var modifiedPrompt = userPrompt + languagePrompt
		fmt.Println("Prompt received. Generating response...")

		// Start a go routine to display a waiting indicator while the response is being generated
		done := make(chan bool)
		go displayindicator.DisplayLoadingIndicator(done)

		// Generate response using Ollama API, passing the context
		response, updatedContext, err := ollamaimplementation.GetOllamaResponse(modifiedPrompt, currentConversationContext, model)

		// Signal the waiting indicator to stop
		done <- true

		if err != nil {
			fmt.Println("Error generating response:", err)
			continue
		}

		// Update the conversation context with the response
		currentConversationContext = updatedContext

		generatedCode, errExtract := extraction.Extract(response) // Handle error with string
		if errExtract != nil {
			fmt.Printf("The LLM gave a improper string in response: %v", response)
			userPrompt = "exit"
			continue
		}

		fmt.Println("Current Iteration:", numOfIterations)
		fmt.Println("Ollama's response:\n\n", generatedCode)

		output, err := compiler([]byte(generatedCode))

		if err != nil {
			fmt.Printf("The code did not compile and contains the following errors: %v\n", string(output))
			userPrompt = "Following are the errors, please fix the code. Write it again, and write only source code along with same test cases with no further explanation. The format should be ```rust <yourcode + testcases> ```  :\n" +  string(output) 

			numOfIterations += 1
		} else {
			fmt.Printf("Compiled successfully. Here is the output: %v", string(output))
			fmt.Println("Total Number of Iterations:", numOfIterations)
			break
		}
	}
}
