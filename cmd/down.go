package cmd

import (
	"fmt"
	"log"

	"github.com/jeduardo/container-compose/internal/container"
	"github.com/jeduardo/container-compose/pkg/compose"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers",
	Run:   down,
}

func down(cmd *cobra.Command, args []string) {
	composeFile, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Fatalf("no compose file informed: %s", err)
	}

	config := compose.Parse(composeFile)
	for name, _ := range config.Services {
		serviceName := fmt.Sprintf("%s_%s_%d", "compose", name, 1)

		fmt.Printf("Stopping %s...", serviceName)
		err := container.Stop(serviceName)
		if err != nil {
			fmt.Println("stop failed:", err)
		}
		if err := container.Remove(serviceName); err != nil {
			fmt.Println("remove failed:", err)
		}
		fmt.Println("done!")
	}
}

func init() {
	rootCmd.AddCommand(downCmd)
}
