package go_compiler_v2

import (
	"llama/modules/compiler_v2/consts"
	"llama/modules/compiler_v2/utils"
	"os"
)

const fileName = "main.go"
const testFileName = "main_test.go"

type GoCompiler struct{}

// NewGoCompiler creates a new GoCompiler
func NewGoCompiler() *GoCompiler {
	return &GoCompiler{}
}

// CheckCompileErrors takes Go source code and checks for compile errors.
//
// The dependencies are handled automatically by go mod and go tidy.
//
// NOTE: Make sure you have an up-to-date Go installed on the system
//
// Returns the output of the compilation and an error if any
func (gb *GoCompiler) CheckCompileErrors(srcCode string, testCode string) ([]byte, error) {
	// Set up temporary folders
	utils.SetupTempFolders(consts.TempOutputDir)
	defer utils.RemoveTempFolders(consts.TempOutputDir)

	// Regular expression to identify and extract test functions
	// re := regexp.MustCompile(`(?m)^func\s+(Test\w+)\s*\(t\s+\*testing\.T\)\s*{[\s\S]*?^}`)

	// Extract test functions from srcCode
	// testFunctions := re.FindAllString(string(srcCode), -1)

	// Remove the test functions from the main source code
	// nonTestContent := re.ReplaceAllString(string(srcCode), "")

	// Remove the "testing" import if it exists in the non-test code
	// nonTestContent = strings.Replace(nonTestContent, `import "testing"`, "", 1)

		
	// Write the cleaned main code to the primary file
	mainFilePath := consts.TempOutputDir + fileName
	err := os.WriteFile(mainFilePath, []byte(srcCode), 0644)
	if err != nil {
		return nil, err
	}


	testFilePath := consts.TempOutputDir + testFileName
	err2 := os.WriteFile(testFilePath, []byte(testCode), 0644)
	if err2 != nil {
		return nil, err
	}

	// Construct the content for the separate _test.go file with the necessary testing import
	// testFileContent := "package main\n\nimport \"testing\"\n\n"
	// for _, match := range testFunctions {
	// 	testFileContent += match + "\n\n"
	// }

	// Write the test code to the _test.go file
	// testFilePath := consts.TempOutputDir + testFileName
	// err2 := os.WriteFile(testFilePath, []byte(testFileContent), 0644)
	// if err2 != nil {
	// 	return nil, err2
	// }

	// Print contents of the main file
	// fmt.Println("Main file content:")
	// mainFileData, err := os.ReadFile(mainFilePath)
	// if err == nil {
	// 	fmt.Println(string(mainFileData))
	// } else {
	// 	fmt.Println("Error reading main file:", err)
	// }

	// // Print contents of the test file
	// fmt.Println("\nTest file content:")
	// testFileData, err := os.ReadFile(testFilePath)
	// if err == nil {
	// 	fmt.Println(string(testFileData))
	// } else {
	// 	fmt.Println("Error reading test file:", err)
	// }

    // Initialize Go module and tidy dependencies
    cmdString := "go mod init tempOutput && go mod tidy"

    // Run the main code file to capture its output
    // cmdString += " && go run main.go"

    // Compile and run tests in the separate test file
    cmdString += " && go build -o main " + fileName
    cmdString += " && go test -v"

	// Execute the command string
	cmd := utils.MakeCommand(cmdString)
	cmd.Dir = consts.TempOutputDir
	return cmd.CombinedOutput()
}


