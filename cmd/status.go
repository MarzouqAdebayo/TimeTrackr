package cmd

import (
	"TimeTrackr/boltDB"
	"TimeTrackr/ui"

	"github.com/spf13/cobra"
)

func init() {
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
		filters := boltdb.FilterObject{
			Status: boltdb.TaskStatus(boltdb.ONGOING),
		}
		// singleResult, err = boltdb.GetTask(idVar)
		multipleResult, err := boltdb.Status(&filters)
		if err != nil {
			cmd.Println(err.Error())
			return
		}
		if len(multipleResult) == 1 {
			// TODO Print single task ui here
			cmd.Println(ui.PrintTaskList(multipleResult))
		} else {
			cmd.Println(ui.PrintTaskList(multipleResult))
		}
	},
}
