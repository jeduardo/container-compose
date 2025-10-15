package cmd

import (
	"fmt"
	"log"

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

	config := compose.Parse(composeFile)
	for name, service := range config.Services {
		serviceName := fmt.Sprintf("%s_%s_%d", "compose", name, 1)
		fmt.Printf("Creating %s...\n", serviceName)
		container.Run(serviceName, service.Image)
	}
}

func init() {
	rootCmd.AddCommand(upCmd)
}
