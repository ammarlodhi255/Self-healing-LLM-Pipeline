package extraction

import (
	"testing"
)

// var GoInputs []string = []string{"```go\nfunc main() {}\n```", "```go\nfmt.Println('Hello World')\n```"}
// var expectedGoOuputs []string = []string{"\nfunc main() {}\n", "\nfmt.Println('Hello World')\n"}

// var RustInputs []string = []string{"```go\nfn main() {}\n```", "```go\nprintln!('Hello World')\n```"}
// var expectedRustOuputs []string = []string{"\nfn main() {}\n", "\nprintln!('Hello World')\n"}

// func TestExtraction(t *testing.T) {

// 	t.Run("Golang Extraction 1", func(t *testing.T) {
// 		var output = Extract(GoInputs[0])
// 		if output != expectedGoOuputs[0] {
// 			t.Error(output)
// 		}
// 	})

// 	t.Run("Golang Extraction 2", func(t *testing.T) {
// 		var output = Extract(GoInputs[1])
// 		if output != expectedGoOuputs[1] {
// 			t.Error(output)
// 		}
// 	})

// 	t.Run("Rust Extraction 1", func(t *testing.T) {
// 		var output = Extract(RustInputs[0])
// 		if output != expectedRustOuputs[0] {
// 			t.Error(output)
// 		}
// 	})

// 	t.Run("Rust Extraction 2", func(t *testing.T) {
// 		var output = Extract(RustInputs[1])
// 		if output != expectedRustOuputs[1] {
// 			t.Error(output)
// 		}
// 	})

// }

// Inputs and Expected Outputs for the Test Cases
// This can be considered a table-driven test or equivalence partitioning
var testCases = []struct {
	name     string
	input    string
	expected string
}{
	// Go Test Cases
	{"Go Extraction 1 - Main", "```go\nfunc main() {}\n```", "\nfunc main() {}\n"},
	{"Go Extraction 2 - Print", "```go\nfmt.Println('Hello World')\n```", "\nfmt.Println('Hello World')\n"},
	{"Go Extraction 3 - Loop", "```go\nfor i := 0; i < 10; i++ {\nfmt.Println(i)\n}\n```", "\nfor i := 0; i < 10; i++ {\nfmt.Println(i)\n}\n"},
	{"Go Extraction 4 - If Else", "```go\nif x > 10 {\nfmt.Println('Greater than 10')\n} else {\nfmt.Println('Less than or equal to 10')\n}\n```", "\nif x > 10 {\nfmt.Println('Greater than 10')\n} else {\nfmt.Println('Less than or equal to 10')\n}\n"},
	{"Go Extraction 5 - Function with Parameters", "```go\nfunc add(a int, b int) int {\nreturn a + b\n}\n```", "\nfunc add(a int, b int) int {\nreturn a + b\n}\n"},
	{"Go Extraction 6 - Nested Loops", "```go\nfor i := 0; i < 3; i++ {\nfor j := 0; j < 3; j++ {\nfmt.Printf('(%d, %d)', i, j)\n}\n}\n```", "\nfor i := 0; i < 3; i++ {\nfor j := 0; j < 3; j++ {\nfmt.Printf('(%d, %d)', i, j)\n}\n}\n"},
	{"Go Extraction 7 - Invalid", "```go```", ""},

	// Rust Test Cases
	{"Rust Extraction 1 - Main", "```rust\nfn main() {}\n```", "\nfn main() {}\n"},
	{"Rust Extraction 2 - Print", "```rust\nprintln!('Hello World')\n```", "\nprintln!('Hello World')\n"},
	{"Rust Extraction 3 - Loop", "```rust\nfor i in 0..10 {\nprintf!(\"{}\", i);\n}\n```", "\nfor i in 0..10 {\nprintf!(\"{}\", i);\n}\n"},
	{"Rust Extraction 4 - If Else", "```rust\nif x > 10 {\nprintln!(\"Greater than 10\");\n} else {\nprintln!(\"Less than or equal to 10\");\n}\n```", "\nif x > 10 {\nprintln!(\"Greater than 10\");\n} else {\nprintln!(\"Less than or equal to 10\");\n}\n"},
	{"Rust Extraction 5 - Function with Parameters", "```rust\nfn add(a: i32, b: i32) -> i32 {\nreturn a + b;\n}\n```", "\nfn add(a: i32, b: i32) -> i32 {\nreturn a + b;\n}\n"},
	{"Rust Extraction 6 - Nested Loops", "```rust\nfor i in 0..3 {\nfor j in 0..3 {\nprintf!(\"({},{})\", i, j);\n}\n}\n```", "\nfor i in 0..3 {\nfor j in 0..3 {\nprintf!(\"({},{})\", i, j);\n}\n}\n"},
}

// Refined Test Function using Table-Driven Approach
func TestExtraction(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := Extract(tc.input)
			if output != tc.expected {
				t.Errorf("Test %s failed: Expected %q, got %q", tc.name, tc.expected, output)
				t.Log(err.Error())

			}
		})
	}
}
