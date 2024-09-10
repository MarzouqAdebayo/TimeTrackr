package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Source string

var rootCmd = &cobra.Command{
	Use:   "trackr",
	Short: "trackr is a time tracking tool",
	Long: `trackr is a CLI time tracking tool to 
  help manage the time you spend on tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Print("Welcome to TimeTrackr, use trackr --help to get started")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
