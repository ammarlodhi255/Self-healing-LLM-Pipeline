package compiler

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const TempOutputDir = "tempOutput/"
const TempModuleName = "tempModule"

type Language string

// Supported languages
const (
	Go   Language = "go"
	Rust Language = "rust"
)

type OS string

// Supported OS
const (
	Windows OS = "windows"
	Linux   OS = "linux"
	MacOS   OS = "darwin" // Darwin is the kernel of macOS
)

// TODO: I want to make an interface for a compilable language, so that we can add more languages in the future
// TODO: The cmd might also be an interface or a struct, so that it can build itself based on the platform and language
// TODO: A cleanup and panic might be needed in setup because if it panics the temp folders should be removed
// TODO: I am not sure that the setup should panic, maybe it should return an error instead so its easier to clean up

type Compiler struct {
	OS            OS
	Language      Language
	languageEnv   ILanguageEnvironment
	SourceCode    string
	Filename      string
	cmdPrefix     string // For example "cmd /c" on Windows
	Dependencies  []string
	tempOutputDir string
}

type ICompiler interface {
	Compile() (string, error)
}

type ILanguageEnvironment interface {
	SetupEnvironment(cmdPrefix string, dependencies []string)
	CheckCompileErrors(filename string, language Language, cmdPrefix string) (string, error)
	WriteCodeToFile(filename, sourceCode string) error
	RunPipeline(c *Compiler) (string, error)
}

type GoEnvironment struct {
}

func (ge *GoEnvironment) RunPipeline(c *Compiler) (string, error) {
	srcCodeFilename := TempOutputDir + appendSuffix(c.Filename, c.Language)
	//compiledFilename := TempOutputDir + c.Filename

	// Write the source code to a file first, because it determines the dependencies for "go mod tidy"
	err := c.languageEnv.WriteCodeToFile(srcCodeFilename, c.SourceCode)
	if err != nil {
		log.Fatalf("Error writing source code to file: %v", err)
	}

	// Sets up go environment with go mod and go mod tidy
	c.languageEnv.SetupEnvironment(c.cmdPrefix, c.Dependencies)

	// CheckCompileErrors the code
	return c.languageEnv.CheckCompileErrors(c.Filename, c.Language, c.cmdPrefix)
}

// SetupEnvironment initializes the go environment by creating a go module and running go mod tidy. Panics if it fails.
// Go modules are used to manage dependencies in go projects.
func (ge *GoEnvironment) SetupEnvironment(cmdPrefix string, _ []string) {
	// One string
	cmdString := cmdPrefix + " go mod init " + TempModuleName + " && go mod tidy"
	// Split the string into a slice
	cmdSlice := strings.Fields(cmdString) // Fields splits the strings around each instance of one or more consecutive white space characters

	// Make the command
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	// Set its target directory
	cmd.Dir = TempOutputDir
	// Execute the command
	err := cmd.Run()
	if err != nil {
		removeTempFolders(TempOutputDir)
		log.Fatalf("Error initializing go module: %v", err)
	}
}

func (ge *GoEnvironment) CheckCompileErrors(filename string, language Language, cmdPrefix string) (string, error) {

	srcCodeFilename := appendSuffix(filename, language)
	compiledFilename := filename

	cmdString := cmdPrefix + " go build -o " + compiledFilename + " " + srcCodeFilename
	cmdSlice := strings.Fields(cmdString) // Fields splits the string on white space of variable length

	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	cmd.Dir = TempOutputDir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (ge *GoEnvironment) WriteCodeToFile(filename, sourceCode string) error {
	return os.WriteFile(filename, []byte(sourceCode), 0644)
}

func InitCompiler(OS OS, language Language, sourceCode string, filename string, dependencies ...string) ICompiler {
	compiler := &Compiler{}
	compiler.OS = OS
	compiler.Language = language
	compiler.SourceCode = sourceCode
	compiler.Filename = filename
	compiler.Dependencies = dependencies
	compiler.cmdPrefix = getOsPrefix(OS)
	compiler.languageEnv = getLanguageEnv(language)
	return compiler

}

func getOsPrefix(OS OS) string {
	// Set the cmd prefix based on the platform
	switch OS {
	case Windows:
		return "cmd /c "
	case Linux, MacOS:
		return ""
	default:
		panic("Unsupported platform")
	}
}

func getLanguageEnv(language Language) ILanguageEnvironment {
	switch language {
	case Go:
		return &GoEnvironment{}
	case Rust:
		return &RustEnvironment{}
	default:
		panic("Unsupported language")
	}
}

type RustEnvironment struct {
}

func (re *RustEnvironment) RunPipeline(c *Compiler) (string, error) {
	srcCodeFilename := TempOutputDir + appendSuffix(c.Filename, c.Language)
	//compiledFilename := TempOutputDir + c.Filename

	// SetupEnvironment either Go or Rust environment, should be an interface method
	c.languageEnv.SetupEnvironment(c.cmdPrefix, c.Dependencies)

	// Write the source code to a file
	err := c.languageEnv.WriteCodeToFile(srcCodeFilename, c.SourceCode)
	if err != nil {
		log.Fatalf("Error writing source code to file: %v", err)
	}

	// CheckCompileErrors the code
	return c.languageEnv.CheckCompileErrors(c.Filename, c.Language, c.cmdPrefix)
}

// SetupEnvironment initializes the rust environment by creating a cargo project and adding dependencies. Panics if it fails.
func (re *RustEnvironment) SetupEnvironment(cmdPrefix string, dependencies []string) {
	// Initialize the rust cargo project--------------------------------------------------------------------------------
	// Command to initialize a cargo project
	cmdString := cmdPrefix + " cargo init --bin"
	// Split the string into a slice
	cmdSlice := strings.Fields(cmdString)
	// Make the command
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	// Set its target directory
	cmd.Dir = TempOutputDir
	// Execute the command
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error initializing rust project: %v", err)
	}

	// Update rust dependencies in cargo.toml file using cargo add (cargo-edit)-----------------------------------------

	if len(dependencies) == 0 {
		return
	}
	addCommand := cmdPrefix + " cargo add"
	addSlice := strings.Fields(addCommand)
	addSlice = append(addSlice, dependencies...)
	cmd = exec.Command(addSlice[0], addSlice[1:]...)
	cmd.Dir = TempOutputDir
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error adding dependencies: %v", err)
	}
}

// CheckCompileErrors checks the code for errors using cargo check. Returns the output and an error.
// Cargo check does not produce an executable, it only checks the code for errors.
// It also does not need a filename, because it checks the whole cargo project.
func (re *RustEnvironment) CheckCompileErrors(_ string, _ Language, cmdPrefix string) (string, error) {

	cmdString := cmdPrefix + " cargo check"
	cmdSlice := strings.Fields(cmdString)
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	cmd.Dir = TempOutputDir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func (re *RustEnvironment) WriteCodeToFile(_, sourceCode string) error {
	srcCodeFilename := TempOutputDir + "src/" + appendSuffix("main.rs", Rust) // Rust source code file is always named main.rs
	return os.WriteFile(srcCodeFilename, []byte(sourceCode), 0644)
}

/*
Compile compiles the source code and returns the output and an error.
The compiler pipeline
1. Set up OS and Language
2. Set up the temp folders
3. Write the source code to a file
4. SetupEnvironment the code
5. Return the output and error
*/
func (c *Compiler) Compile() (string, error) {
	// Set up temp folders
	setupTempFolders(TempOutputDir)
	defer removeTempFolders(TempOutputDir)

	// CheckCompileErrors the code
	return c.languageEnv.RunPipeline(c)
}

// appendSuffix appends the suffix to the filename if it is not already there depending on the language, panics if the language is not supported
func appendSuffix(filename string, language Language) string {

	suffix := ""
	switch language {
	case Go:
		suffix = ".go"
	case Rust:
		suffix = ".rs"
	default:
		panic("Unsupported language")
	}

	// We check if the filename already has the suffix, if not we append it
	if !strings.HasSuffix(filename, suffix) {
		filename += suffix
	}
	return filename
}

// setupTempFolders creates the temp output directory for compiled files, panics if it fails
func setupTempFolders(tempOutputDir string) {
	// 0777 are the permissions for the directory, everyone can read, write and execute
	err := os.MkdirAll(tempOutputDir, os.ModePerm)
	if err != nil {
		panic("Error creating temp output directory:\n\n" + err.Error())
	}
}

// removeTempFolders removes the temp output directory for compiled files, panics if it fails
func removeTempFolders(tempOutputDir string) {
	err := os.RemoveAll(tempOutputDir)
	if err != nil {
		panic("Error removing temp output directory:\n\n" + err.Error())
	}
}
