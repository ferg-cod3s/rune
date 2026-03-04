package main

import (
	"os/exec"
	"testing"
)

func TestBinaryBuild(t *testing.T) {
	// Verify the binary can be built successfully
	cmd := exec.Command("go", "build", "-o", "/dev/null", ".")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("binary build failed: %v\n%s", err, output)
	}
}
