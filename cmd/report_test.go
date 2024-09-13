package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestReportCmd(t *testing.T) {
	t.Run("report - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		testCmd := &cobra.Command{
			Use:   reportCmd.Use,
			Short: reportCmd.Short,
			Run:   reportCmd.Run,
		}
		testCmd.SetOut(&buf)
		testCmd.SetArgs(args)

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
