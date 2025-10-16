package container

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

// runCommand starts `container <args...>` and returns its stdin, stdout, stderr, and the Cmd.
// Caller is responsible for calling cmd.Wait().
func runCommand(args ...string) (stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser, cmd *exec.Cmd, err error) {
	// 1) Ensure the binary can be found
	if _, lookErr := exec.LookPath("container"); lookErr != nil {
		return nil, nil, nil, nil, fmt.Errorf("container binary not found in PATH: %w", lookErr)
	}

	// 2) Build the command and wire pipes BEFORE Start
	cmd = exec.Command("container", args...)

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err = cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("stderr pipe: %w", err)
	}

	// 3) Start the process
	if err = cmd.Start(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("start: %w", err)
	}

	return stdin, stdout, stderr, cmd, nil
}

// Run runs a named container and streams its logs until exit.
func Run(containerName, image string) error {
	args := []string{"run", "--name", containerName, image}
	stdin, stdout, stderr, cmd, err := runCommand(args...)
	if err != nil {
		return err
	}
	// Closing stdin as the child doesn't need input.
	// NOTE: change when supporting -i flag.
	_ = stdin.Close()

	// Stream stdout
	outDone := make(chan error, 1)
	go func() {
		sc := bufio.NewScanner(stdout)
		// Optional: grow max token size beyond 64 KiB
		buf := make([]byte, 0, 128*1024)
		sc.Buffer(buf, 10*1024*1024)
		for sc.Scan() {
			fmt.Printf("%s | stdout | %s\n", containerName, sc.Text())
		}
		outDone <- sc.Err()
	}()

	// Stream stderr
	errDone := make(chan error, 1)
	go func() {
		sc := bufio.NewScanner(stderr)
		sc.Buffer(make([]byte, 0, 128*1024), 10*1024*1024)
		for sc.Scan() {
			// send to stderr so users/tools can differentiate
			fmt.Fprintf(os.Stderr, "%s | stderr | %s\n", containerName, sc.Text())
		}
		errDone <- sc.Err()
	}()

	waitErr := cmd.Wait()

	// Give scanners a moment to drain after process exit (best-effort).
	select {
	case e := <-outDone:
		if e != nil {
			if !errors.Is(e, io.EOF) {
				fmt.Fprintf(os.Stderr, "stdout scan error: %v\n", e)
			}
		}
	case <-time.After(200 * time.Millisecond):
	}
	select {
	case e := <-errDone:
		if e != nil {
			if !errors.Is(e, io.EOF) {
				fmt.Fprintf(os.Stderr, "stderr scan error: %v\n", e)
			}
		}
	case <-time.After(200 * time.Millisecond):
	}

	return waitErr
}
