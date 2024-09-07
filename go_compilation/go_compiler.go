package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Go code to compile
	goCode := `
package main

import "fmt"

func main() {
	fmt.Println("Hello from Go!")
}
`

	// Path to the temporary Go source file
	sourceFile := "main.go"
	// Path to the binary output
	binaryFile := "main"

	// Write Go code to a file
	if err := os.WriteFile(sourceFile, []byte(goCode), 0644); err != nil {
		fmt.Println("Error writing Go code to file:", err)
		return
	}
	defer func() {
		os.Remove(sourceFile) // Clean up the source file
		os.Remove(binaryFile) // Clean up the binary file
	}()

	fmt.Println("Go code written to", sourceFile)

	// Compile Go code
	cmd := exec.Command("go", "build", "-o", binaryFile, sourceFile)
	var out, stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("Compilation failed with errors:")
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("Compilation succeeded. Output:")
	fmt.Println(out.String())

	// If compilation succeeded, run the binary and capture output
	cmd = exec.Command("./" + binaryFile)
	var result, runStderr strings.Builder
	cmd.Stdout = &result
	cmd.Stderr = &runStderr
	err = cmd.Run()

	if err != nil {
		fmt.Println("Error running the compiled binary:")
		fmt.Println(runStderr.String())
		return
	}

	// Output from the Go program
	fmt.Println("Go program output:")
	fmt.Println(result.String())
}
