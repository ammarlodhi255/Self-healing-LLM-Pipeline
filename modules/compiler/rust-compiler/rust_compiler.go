package rust_compiler

import (
	"compiler"
	"runtime"
)

// Deprecated: Use `rust_compiler_v2.NewRustCompiler` instead
//
// CompileStringToRust compiles a string of go code to a rust executable
func CompileStringToRust(code string, filename string, dependencies ...string) (string, error) {

	// Get the platform
	OS := runtime.GOOS

	// SetupEnvironment
	return compiler.InitCompiler(compiler.OS(OS), compiler.Rust, code, filename, dependencies...).Compile()
}
