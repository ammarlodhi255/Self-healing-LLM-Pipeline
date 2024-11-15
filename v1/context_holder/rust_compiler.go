package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"
// )

// func runCode(rustCode string) (bool, string) {
// 	projectDir := "rust_project"
// 	os.Mkdir(projectDir, 0755)

// 	// Create the Cargo.toml file with necessary dependencies
// 	cargoToml := `
// 	[package]
// 	name = "rust_project"
// 	version = "0.1.0"
// 	edition = "2021"

// 	[dependencies]
// 	rand = "0.8"
// 	`
// 	cargoTomlPath := projectDir + "/Cargo.toml"
// 	if err := os.WriteFile(cargoTomlPath, []byte(cargoToml), 0644); err != nil {
// 		fmt.Println("Error writing Cargo.toml:", err)
// 		return false, ""
// 	}

// 	// Create the src directory and main.rs file
// 	os.Mkdir(projectDir+"/src", 0755)
// 	mainRsPath := projectDir + "/src/main.rs"
// 	if err := os.WriteFile(mainRsPath, []byte(rustCode), 0644); err != nil {
// 		fmt.Println("Error writing Rust code to file:", err)
// 		return false, ""
// 	}

// 	// Compile the Rust project using Cargo
// 	cmd := exec.Command("cargo", "build", "--release")
// 	cmd.Dir = projectDir
// 	var out, stderr strings.Builder
// 	cmd.Stdout = &out
// 	cmd.Stderr = &stderr
// 	err := cmd.Run()

// 	if err != nil {
// 		fmt.Println("Compilation failed with errors:")
// 		return false, stderr.String()
// 	}

// 	// If compilation succeeded, run the binary and capture output
// 	cmd = exec.Command("./target/release/rust_project")
// 	cmd.Dir = projectDir
// 	var result strings.Builder
// 	cmd.Stdout = &result
// 	cmd.Stderr = &stderr
// 	err = cmd.Run()

// 	if err != nil {
// 		fmt.Println("Error running the compiled binary:")
// 		return false, stderr.String()
// 	}

// 	// Output from the Rust program
// 	fmt.Println("Rust program output:")
// 	fmt.Println(result.String())

// 	// Clean up: Remove the rust_project directory
// 	if err := os.RemoveAll(projectDir); err != nil {
// 		fmt.Println("Error cleaning up the project directory:", err)
// 	}

// 	return true, result.String()
// }
