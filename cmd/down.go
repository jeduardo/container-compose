package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers",
	Run:   down,
}

func down(cmd *cobra.Command, args []string) {
	fmt.Println("down called")
}

func init() {
	rootCmd.AddCommand(downCmd)
}
