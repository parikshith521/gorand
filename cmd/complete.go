package cmd

import (
	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [task to complete]",
	Short: "Completes a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteRecordFromCSVFile("data.csv", args[0], true)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
