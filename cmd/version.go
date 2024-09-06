package cmd

import (
	"github.com/spf13/cobra"
)

const versionResponse = "TimeTrackr v0.0.1 -- HEAD\n"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of trackr",
	Long:  `All software have versions. This is trackr's`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Print(versionResponse)
	},
}
