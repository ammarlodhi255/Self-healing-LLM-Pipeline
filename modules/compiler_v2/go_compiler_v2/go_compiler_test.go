package go_compiler_v2

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
		{
			filename:      "should_compile_and_run_tests",
			shouldCompile: true,
		},
		{ // TODO might change name from should compile to should succeed
			filename:      "should_compile_with_faulty_test", // Code is syntactically correct, but the test is faulty
			shouldCompile: false,                             // Here the test is faulty, so it will get a compiler error
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			// Read the code from the file
			code, err := os.ReadFile(test.filename)

			output, err := NewGoCompiler().CheckCompileErrors(code)

			if err != nil && test.shouldCompile {
				t.Errorf("Expected the code to compile, but got an output: %v \n error: %v", string(output), err)
			} else if err == nil && !test.shouldCompile {
				t.Errorf("Expected the code to not compile, but got no error")
			}

			// Check if the output is empty when the code shouldn't compile
			if output == nil && !test.shouldCompile {
				t.Errorf("Expected compiler error output, but got none")
			}
		})
	}
}
