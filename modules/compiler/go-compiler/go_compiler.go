package go_compiler

import (
	"compiler"
	"runtime"
)

// Deprecated: Use `go_compiler_v2.NewGoCompiler` instead
//
// CompileStringToGo tries to compile a string of go code to a go executable, and returns the compiler output and an error.
// The function does not produce any executables, since they are deleted after the function ends.
func CompileStringToGo(code string, filename string) (string, error) {

	// Get the platform
	OS := runtime.GOOS

	// SetupEnvironment
	return compiler.InitCompiler(compiler.OS(OS), compiler.Go, code, filename).Compile()
}
