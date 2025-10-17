package container

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/jeduardo/container-compose/internal/system"
)

// Run runs a named container and streams its logs until exit.
func Run(containerName, image string) *exec.Cmd {
	args := []string{"run", "--name", containerName, image}
	cmd, err := system.Run("container", args...)
	if err != nil {
		log.Fatalln(err)
	}
	return cmd
}

// Stops a running container by name and waits for it to exit.
// It uses a context with a timeout to avoid hanging forever.
func Stop(containerName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "container", "stop", containerName)
	if err := cmd.Run(); err != nil {
		// If the context expired, surface that clearly
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout stopping container %q", containerName)
		}
		return fmt.Errorf("failed to stop container %q: %w", containerName, err)
	}
	return nil
}

// Removes a stopped container by name.
// It also uses a short timeout to prevent hangs.
func Remove(containerName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "container", "rm", containerName)
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout removing container %q", containerName)
		}
		return fmt.Errorf("failed to remove container %q: %w", containerName, err)
	}
	return nil
}
