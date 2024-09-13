package cmd

import (
	boltdb "TimeTrackr/boltDB"

	"github.com/spf13/cobra"
)

func init() {
	pauseCmd.Flags().IntVar(&idVar, "id", 0, "")
	rootCmd.AddCommand(pauseCmd)
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause current time-tracking session",
	Long: `Pauses the current time-tracking session. 

Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var msg string
		var err error
		if idVar <= 0 {
			msg, err = boltdb.PauseCurrentTask(nil)
		} else {
			msg, err = boltdb.PauseCurrentTask(&idVar)
		}
		if err != nil {
			cmd.Println(err.Error())
			return
		}
		cmd.Println(msg)
	},
}
