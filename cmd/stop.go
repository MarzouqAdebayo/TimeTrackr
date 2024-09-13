package cmd

import (
	b "TimeTrackr/boltDB"

	"github.com/spf13/cobra"
)

func init() {
	stopCmd.Flags().IntVar(&idVar, "id", 0, "")
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop current time-tracking session",
	Long: `Stops the current time-tracking session. 
It will save the data of the current session before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var status string
		var err error
		if idVar <= 0 {
			status, err = b.StopCurrentTask(nil)
		} else {
			status, err = b.StopCurrentTask(&idVar)
		}
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Println("Current time tracking session stopped")
		cmd.Println(status)
	},
}
