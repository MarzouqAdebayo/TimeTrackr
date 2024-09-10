package cmd

import (
	"TimeTrackr/boltDB"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("status", "s", "ongoing", "gets task status")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get a summary into your current time tracking session",
	Long: `Get a summary into your current time tracking session. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		status, _ := cmd.Flags().GetString("status")
		result, err := boltdb.Status(status)
		if err != nil {
			cmd.Println(err.Error())
			return
		}
		cmd.Println(result)
	},
}
