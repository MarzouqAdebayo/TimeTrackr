package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop current time-tracking session",
	Long: `Stops the current time-tracking session. 
It will save the data of the current session before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Print("Echo: " + strings.Join(args, " "))
	},
}
