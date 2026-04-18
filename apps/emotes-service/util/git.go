package util

import (
	"os/exec"
	"strings"
)

func GitCli(args ...string) string {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "Unknown"
	}
	if s := strings.TrimSpace(string(out)); s != "" {
		return s
	}
	return "Unknown"
}
