package cmd

import (
	"TimeTrackr/boltDB"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("category", "c", "miscellaneous", "group this task into a category.")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Initiate a new time-tracking session",
	Long: `Starts a new time-tracking session for the specified task. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		categoryFlag, _ := cmd.Flags().GetString("category")
		err := boltdb.StartTask(args[0], categoryFlag)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Printf("A new time tracking session started for task %s", args[0])
	},
}
