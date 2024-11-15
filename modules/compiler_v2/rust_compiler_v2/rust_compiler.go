package rust_compiler_v2

import (
	"llama/modules/compiler_v2/consts"
	"llama/modules/compiler_v2/utils"
	"os"
	"strings"
)

const fileName = "main.rs"

type RustCompiler struct{}

// NewRustCompiler creates a new RustCompiler
func NewRustCompiler() *RustCompiler {
	return &RustCompiler{}
}

// CheckCompileErrors takes Rust source code and the dependencies it requires and checks for compile errors.
//
// The dependencies are optional, and should be name only, not version.
// For instance "rand" and not "rand:0.8.3". Cargo will automatically fetch the latest version.
//
// Returns the output of the compilation and an error if any
func (gb *RustCompiler) CheckCompileErrors(srcCode []byte, dependencies ...string) ([]byte, error) {
	// Make temp folders
	utils.SetupTempFolders(consts.TempOutputDir)
	defer utils.RemoveTempFolders(consts.TempOutputDir)

	// Init cargo
	if err := initCargo(); err != nil {
		return nil, err
	}

	// Write code to file
	if err := os.WriteFile(consts.TempOutputDir+"src/"+fileName, srcCode, 0644); err != nil {
		return nil, err
	}

	cmdString := ""
	// Add dependencies
	if dependencies != nil {
		cmdString = "cargo add " + strings.Join(dependencies, " ") + " &&"
	}

	// Run go build
	cmdString += " cargo build"

	cmdString += " && cargo test"

	//cmdSlice := strings.Fields(cmdString)
	cmd := utils.MakeCommand(cmdString)
	cmd.Dir = consts.TempOutputDir
	return cmd.CombinedOutput()
}

// initCargo initializes a cargo project
func initCargo() error {
	// Init cargo
	cmd := utils.MakeCommand("cargo init --bin")
	cmd.Dir = consts.TempOutputDir
	return cmd.Run()
}
