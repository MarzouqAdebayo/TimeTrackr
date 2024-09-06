package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pauseCmd)
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause current time-tracking session",
	Long: `Pauses the current time-tracking session. 

Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("Echo: %s", strings.Join(args, " "))
	},
}
