package displayindicator

import (
	"fmt"
	"time"
)

// Function to display a waiting/loading indicator
func DisplayLoadingIndicator(done chan bool) {
	indicator := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r") // Clear the waiting indicator when done
			return
		default:
			fmt.Printf("\r%s Generating...", indicator[i%len(indicator)])
			i++
			time.Sleep(200 * time.Millisecond)
		}
	}
}
