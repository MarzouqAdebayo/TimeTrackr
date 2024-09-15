package cmd

import (
	boltdb "TimeTrackr/boltDB"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func init() {
	startCmd.Flags().StringP("category", "c", boltdb.DEFAULT_CATEGORY, "group this task into a category. will use "+boltdb.DEFAULT_CATEGORY+" by default")
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
		categoryFlag, _ := cmd.Flags().GetString("category")
		if categoryFlag == boltdb.DEFAULT_CATEGORY {
			cmd.Printf("category flag not passed, defaulting to '%s'\n", boltdb.DEFAULT_CATEGORY)
		}
		// TODO Get default start date from config file
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("5"))

		err := boltdb.StartTask(args[0], categoryFlag)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Printf(style.Render(fmt.Sprintf("A new time tracking session started for task %s (%s)\n", args[0], categoryFlag)))
	},
}
