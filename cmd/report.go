package cmd

import (
	// "TimeTrackr/boltDB"
	boltdb "TimeTrackr/boltDB"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	startDateStr string
	endDateStr   string
	full         bool
	analysis     bool
	statistics   bool
	tdd          bool
	tdh          bool
	misc         bool
)

func init() {
	reportCmd.Flags().StringVar(&startDateStr, "startdate", "", "Start date for the report (e.g., 2024-09-01 or 2024-09-01 15:04:05)")
	reportCmd.Flags().StringVar(&endDateStr, "enddate", "", "End date for the report (defaults to now if not provided)")
	reportCmd.Flags().BoolVarP(&full, "full", "f", false, "Show full analysis")
	reportCmd.Flags().BoolVarP(&analysis, "analysis", "a", false, "Show category analysis")
	reportCmd.Flags().BoolVarP(&statistics, "statistics", "s", false, "Show completion statistics")
	reportCmd.Flags().BoolVarP(&tdd, "daydist", "y", false, "Show time distribution by day")
	reportCmd.Flags().BoolVarP(&tdh, "hourdist", "r", false, "Show time distribution by hour")
	reportCmd.Flags().BoolVarP(&misc, "misc", "m", false, "Show miscellaneous insights")

	rootCmd.MarkFlagRequired("startdate")
	rootCmd.AddCommand(reportCmd)
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a time tracking report",
	Long:  `Generates a report of tasks tracked between the given start and end dates.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		loc := time.Local
		startDate, err := parseDateWithTimezone(startDateStr, loc)
		if err != nil {
			cmd.Println("Invalid start date format. Please use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
			return
		}

		var endDate time.Time
		if endDateStr == "" {
			endDate = time.Now().In(loc)
		} else {
			endDate, err = parseDateWithTimezone(endDateStr, loc)
			if err != nil {
				cmd.Println("Invalid end date format. Please use YYYY-MM-DD or YYYY-MM-DD HH:MM:SS")
				return
			}
		}

		startTimestamp := startDate.Unix()
		endTimestamp := endDate.Unix()

		if startTimestamp >= endTimestamp {
			cmd.Println("Error: Start date must be before end date")
			return
		}

		cmd.Printf("startTimestamp: %d, endTimestamp: %d\n", startTimestamp, endTimestamp)
		_report, err := boltdb.GenerateReport(startTimestamp, endTimestamp)
		if err != nil {
			fmt.Println(err.Error())
		}
		cmd.Println(_report)
	},
}

func parseDateWithTimezone(dateStr string, loc *time.Location) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
	}
	var parsedDate time.Time
	var err error
	for _, layout := range layouts {
		parsedDate, err = time.ParseInLocation(layout, dateStr, loc)
		if err == nil {
			return parsedDate, nil
		}
	}
	return time.Time{}, fmt.Errorf("Invalid date format")
}
