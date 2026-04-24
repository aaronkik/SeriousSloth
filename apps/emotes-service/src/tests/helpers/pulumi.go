package helpers

import (
	"os/exec"
	"strings"
	"testing"
)

func GetPulumiExport(t *testing.T, exportName string) string {
	stackName, err := exec.Command("../../../scripts/stack-name.sh").Output()
	if err != nil {
		t.Fatalf("Failed to get stack name: %v", err)
	}

	stack := strings.TrimSpace(string(stackName))

	stackExportValue, err := exec.Command(
		"pulumi",
		"stack",
		"output",
		exportName,
		"-s",
		stack,
	).CombinedOutput()

	if err != nil {
		t.Fatalf("Failed to get stack export %s: %v: %s", exportName, err, stackExportValue)
	}

	return strings.TrimSpace(string(stackExportValue))
}
