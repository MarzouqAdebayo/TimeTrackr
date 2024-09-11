package cmd

import (
	"TimeTrackr/boltDB"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(continueCmd)
}

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "Continue a paused time-tracking session",
	Long: `Continues a paused time-tracking session for the specified task. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := boltdb.ContinuePausedTask(args[0])
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Printf("Time tracking session continued for task %s", args[0])
	},
}
