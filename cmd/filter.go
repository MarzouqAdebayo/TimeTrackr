package cmd

import (
	"TimeTrackr/boltDB"
	"TimeTrackr/ui"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().IntVar(&idVar, "id", 0, "specific task ID. Usage `trackr status --id=1`")
	filterCmd.Flags().BoolP("all", "a", false, "gets all tasks. This flag overrides all other flags")
	filterCmd.Flags().StringP("status", "s", "", "filter by status")
	filterCmd.Flags().StringP("category", "c", "", "filter by category")
	filterCmd.Flags().StringP("name", "n", "", "filter by name. Usage `trackr status -n=one`")
	filterCmd.Flags().String("startdate", "", "filter by start date (e.g., 2024-09-01 or 2024-09-01 15:04:05)")
	filterCmd.Flags().String("enddate", "", "filter by end date (e.g., 2024-09-01 or 2024-09-01 15:04:05)")
	filterCmd.Flags().String("minDuration", "", "filter by min duration (e.g., 2h3m30s or 2H3M30S)")
	filterCmd.Flags().String("maxDuration", "", "filter by max duration (e.g., 2h3m30s or 2H3M30S)")
}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Get a list of previous time tracking sessions based on filter flags",
	Long: `Get a list of previous time tracking sessions based on filter flags. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// var singleResult string
		var multipleResult [][]string

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		status, err := cmd.Flags().GetString("status")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		category, err := cmd.Flags().GetString("category")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		startDate, err := cmd.Flags().GetString("startdate")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		endDate, err := cmd.Flags().GetString("enddate")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		startTimestamp, endTimestamp, err := boltdb.ParseDateCommand(startDate, endDate)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		minDuration, err := cmd.Flags().GetString("minDuration")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		maxDuration, err := cmd.Flags().GetString("maxDuration")
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		minDurationTm, err := boltdb.DurationParser(minDuration)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		maxDurationTm, err := boltdb.DurationParser(maxDuration)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		filters := boltdb.FilterObject{
			Name:        name,
			Status:      boltdb.TaskStatus(status),
			Category:    category,
			StartDate:   startTimestamp,
			EndDate:     endTimestamp,
			MinDuration: minDurationTm,
			MaxDuration: maxDurationTm,
		}
		if all {
			filters = boltdb.FilterObject{}
		}

		if idVar <= 0 {
			multipleResult, err = boltdb.Filter(&filters)
		} else {
			// singleResult, err = boltdb.GetTask(idVar)
			multipleResult, err = boltdb.Filter(&filters)
		}
		if err != nil {
			cmd.Println(err.Error())
			return
		}
		if len(multipleResult) == 1 && idVar > 0 {
			// TODO Print single task ui here
			cmd.Println(ui.PrintTaskList(multipleResult))
		} else {
			cmd.Println(ui.PrintTaskList(multipleResult))
		}
	},
}
