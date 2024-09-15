package cmd

import (
	boltdb "TimeTrackr/boltDB"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// TODO Save installation date into env or some config file
	startDateStr string
	endDateStr   string
)

func init() {
	reportCmd.Flags().StringVar(&startDateStr, "startdate", "", "Start date for the report (e.g., 2024-09-01 or 2024-09-01 15:04:05)")
	reportCmd.Flags().StringVar(&endDateStr, "enddate", "", "End date for the report (defaults to now if not provided)")
	rootCmd.MarkFlagRequired("startdate")
	rootCmd.AddCommand(reportCmd)
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a time tracking report",
	Long:  `Generates a report of tasks tracked between the given start and end dates.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		startTimestamp, endTimestamp, err := boltdb.ParseDateCommand(startDateStr, endDateStr)
		if err != nil {
			cmd.Println(err.Error())
			return
		}

		_report, err := boltdb.GenerateReport(startTimestamp, endTimestamp)
		if err != nil {
			fmt.Println(err.Error())
		}
		var style = lipgloss.NewStyle().
			Bold(true).
			// Background(lipgloss.Color("#FAFAFA")).
			Foreground(lipgloss.Color("#7D56F4")).
			PaddingTop(1).
			PaddingLeft(4)
		cmd.Println(style.Render(_report))
	},
}
