package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jeduardo/container-compose/internal/container"
	"github.com/jeduardo/container-compose/pkg/compose"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers",
	Run:   up,
}

func up(cmd *cobra.Command, args []string) {
	composeFile, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Fatalf("no compose file informed: %s", err)
	}

	// Capture Ctrl-C or SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// launch services
	config := compose.Parse(composeFile)
	serviceCmds := make(map[string]*exec.Cmd)
	for name, service := range config.Services {
		serviceName := fmt.Sprintf("%s_%s_%d", "compose", name, 1)
		fmt.Printf("Creating %s...\n", serviceName)
		serviceCmds[serviceName] = container.Run(serviceName, service.Image)
	}

	go func() {
		sig := <-sigChan
		for _, serviceCmd := range serviceCmds {
			_ = serviceCmd.Process.Signal(sig)
		}
	}()

	// Wait for all containers to exit or terminate with Ctrl-C
	var wg sync.WaitGroup
	wg.Add(len(serviceCmds))
	for name, serviceCmd := range serviceCmds {
		go func() {
			defer wg.Done()
			_ = serviceCmd.Wait()
			container.Remove(name)
		}()
	}
	wg.Wait()
}

func init() {
	rootCmd.AddCommand(upCmd)
}
