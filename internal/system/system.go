package system

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Run a binary if it exists, stream its stdout/stderr, and return the Cmd.
func Run(name string, args ...string) (*exec.Cmd, error) {
	// Verify command exists in PATH
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("command %q not found in PATH", name)
	}

	cmd := exec.Command(path, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("%s stdout: %w", name, err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("%s stderr: %w", name, err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start %s: %w", name, err)
	}

	// Stream stdout/stderr in background goroutines
	// TODO: prefix with container name
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	return cmd, nil
}
