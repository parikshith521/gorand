package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gorand",
	Short: "A CLI tool to manage tasks",
	Long:  `This application has four commands: list, add, complete, and delete.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
