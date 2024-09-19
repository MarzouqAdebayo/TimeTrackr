package cmd

import (
	"TimeTrackr/web"

	"github.com/spf13/cobra"
)

func init() {
	webServerCmd.Flags().IntP("port", "p", 0, "run server on custom port")
	rootCmd.AddCommand(webServerCmd)
}

var webServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a web server to interact with TimeTrackr",
	Long: `Start a web server to interact with TimeTrackr. 
If a timer is currently running, it will automatically stop that timer and save its data before starting a new session for the provided task name. 
Use this command to accurately track the time spent on each activity throughout your day.`,
	// Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		portFlag, _ := cmd.Flags().GetInt("port")
		if portFlag == 0 {
			portFlag = 8080
		}
		web.RunServer(&portFlag)
	},
}
