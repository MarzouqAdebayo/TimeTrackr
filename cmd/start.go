package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Initiate a new time-tracking session",
	Long: `Starts a new time-tracking session for the specified task. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create function for all these. They return a string to be printed to the user
		// get ongoing task from db and tell user that a task is currently running
		// if no ongoing task, create a task with intialized values, duration is 0 and endtime is 0
		// If not error, return response to user
		cmd.Printf("Task %s started", args[0])
	},
}
