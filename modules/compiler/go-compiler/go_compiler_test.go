package go_compiler

import (
	"os"
	"testing"
)

func TestCompileStringToGo(t *testing.T) {

	tests := []struct {
		filename      string
		shouldCompile bool
	}{
		{
			filename:      "should_compile",
			shouldCompile: true,
		},
		{
			filename:      "should_not_compile",
			shouldCompile: false,
		},
		{
			filename:      "should_compile_with_standard_library_dependencies",
			shouldCompile: true,
		},
		{
			filename:      "should_compile_with_external_dependencies",
			shouldCompile: true,
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			// Read the code from the file
			code, err := os.ReadFile(test.filename)

			output, err := CompileStringToGo(string(code), test.filename)

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
