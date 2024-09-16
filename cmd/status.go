package cmd

import (
	"TimeTrackr/boltDB"
	"TimeTrackr/ui"

	"github.com/spf13/cobra"
)

func init() {
	statusCmd.Flags().IntVar(&idVar, "id", 0, "View a specific tracking session by ID")
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get a summary into your current time tracking session",
	Long: `Get a summary into your current time tracking session. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var result string
		var err error
		if idVar <= 0 {
			result, err = boltdb.Status(nil)
		} else {
			result, err = boltdb.Status(&idVar)
		}
		if err != nil {
			cmd.Println(err.Error())
			return
		}
		cmd.Println(ui.PrintTask(result))
	},
}
