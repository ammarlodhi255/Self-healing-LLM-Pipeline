package utils

import (
	p "llama/modules/compiler_v2/platform"
	"os/exec"
	"runtime"
)

// MakeCommand creates a command based on the runtime platform
func MakeCommand(cmd string) *exec.Cmd {
	platform := runtime.GOOS
	switch platform {
	case p.Windows:
		return exec.Command("cmd", "/c", cmd)
	case p.Linux, p.MacOS:
		return exec.Command("bash", "-c", cmd)
	default:
		panic("Unsupported platform")
	}
}
