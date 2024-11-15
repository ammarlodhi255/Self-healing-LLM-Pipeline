package extraction

import (
	"fmt"
	// "strings"
	"regexp"
)

// var GoPrompt = "The code should be in the Go programming language. There should also be 3 robust test cases within the same file, these test cases should use 'testing' module. There should also be a main function inside of which the execution of the implemented function takes place. Please always provide the source code and no further explanation, The format should be ```go <yourcode + testcases> ```"

var GoPrompt = "The code should be in the Go programming language. You should generate two files. One file should contain the main code (main function plus a function that performs the required job) and another file should contain 3 robust test cases. Please note that only use main package in both files. Please always provide the source codes and no further explanation, The format should be ```go main code``` and ```go testcode```"

var RustPrompt = "The code should be in the Rust programming language. There should also be 3 robust test cases within the same code. There should also be a main function inside of which all the execution takes place. Please only provide the source code and no further explanation, The format should be ```rust <yourcode + testcases> ```"

// func Extract(output string) string {
// 	parts := strings.Split(output, "```")
// 	var extracted = ""
// 	if strings.Contains(parts[1], "rust") {
// 		extracted = strings.TrimLeft(parts[1], "rust")
// 	} else {
// 		extracted = strings.TrimLeft(parts[1], "go")
// 	}
// 	return extracted
// }

// Extract extracts the code snippet between ``` and removes the language identifier.

func Extract(input string) (string, string, error) {
	// Define regex patterns to match each code block starting with "package main"
	mainPattern := regexp.MustCompile("(?s)```.*?\\n(.*?package main.*?\\n.*?)```")
	testPattern := regexp.MustCompile("(?s)```.*?\\n(.*?package main.*?\\n.*?)```.*?```.*?\\n(.*?package main.*?\\n.*?)```")

	// Match the main code block (first match)
	mainMatch := mainPattern.FindStringSubmatch(input)
	if len(mainMatch) < 2 {
		return "", "", fmt.Errorf("main code block not found")
	}
	mainCode := mainMatch[1]

	// Match the test code block (second match)
	testMatch := testPattern.FindStringSubmatch(input)
	if len(testMatch) < 3 {
		return "", "", fmt.Errorf("test code block not found")
	}
	testCode := testMatch[2]

	return mainCode, testCode, nil
}
