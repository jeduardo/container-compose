package cmd

import (
	"fmt"
	"log"

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
	fmt.Println(config)
}

func init() {
	rootCmd.AddCommand(upCmd)
}
