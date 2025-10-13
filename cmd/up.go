package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Create and start containers",
	Run:   up,
}

func up(cmd *cobra.Command, args []string) {
	fmt.Println("up called")
}

func init() {
	rootCmd.AddCommand(upCmd)
}
