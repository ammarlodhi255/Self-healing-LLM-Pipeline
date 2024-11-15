package promptlist

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func PromptList() []string {
	promptText, err := os.Open("data/promptList.txt")
	if err != nil {
		fmt.Println("Error opening file", err)
	}

	text, err := io.ReadAll(promptText)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	prompt := string(text)

	promptList := strings.Split(prompt, ",")

	return promptList
}
