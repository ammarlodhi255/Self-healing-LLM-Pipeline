package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Rust code to compile
	rustCode := `
	fn main() {
		println!(a"Hello from Rust!");
	}
	`

	// Path to the temporary Rust source file
	sourceFile := "main.rs"
	// Path to the binary output
	binaryFile := "main"

	// Write Rust code to a file
	if err := os.WriteFile(sourceFile, []byte(rustCode), 0644); err != nil {
		fmt.Println("Error writing Rust code to file:", err)
		return
	}
	defer os.Remove(sourceFile) // Clean up the file

	// Compile Rust code
	cmd := exec.Command("rustc", sourceFile, "-o", binaryFile)
	var out, stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("Compilation failed with errors:")
		fmt.Println(stderr.String())
		return
	}

	// If compilation succeeded, run the binary and capture output
	cmd = exec.Command("./" + binaryFile)
	var result strings.Builder
	cmd.Stdout = &result
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		fmt.Println("Error running the compiled binary:")
		fmt.Println(stderr.String())
		return
	}

	// Output from the Rust program
	fmt.Println("Rust program output:")
	fmt.Println(result.String())
}
