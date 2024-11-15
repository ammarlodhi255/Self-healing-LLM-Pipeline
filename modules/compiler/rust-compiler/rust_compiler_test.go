package rust_compiler

import (
	"os"
	"testing"
)

func TestCompileStringToRust(t *testing.T) {

	tests := []struct {
		filename      string
		shouldCompile bool
		dependencies  []string
	}{
		{
			filename:      "should_compile",
			shouldCompile: true,
			dependencies:  nil,
		},
		{
			filename:      "should_not_compile",
			shouldCompile: false,
			dependencies:  nil,
		},
		{
			filename:      "should_compile_with_dependencies",
			shouldCompile: true,
			dependencies:  []string{"rand", "colored"},
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			// Read the code from the file
			code, err := os.ReadFile(test.filename)

			output, err := CompileStringToRust(string(code), "main.rs", test.dependencies...)

			if err != nil && test.shouldCompile {
				t.Errorf("Expected the code to compile, but got an error: %v", err)
			} else if err == nil && !test.shouldCompile {
				t.Errorf("Expected the code to not compile, but got no error")
			}
			// Check if the output is empty when the code shouldn't compile
			if output == "" && !test.shouldCompile {
				t.Errorf("Expected compiler error output, but got none")
			}
		})
	}
}
