package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestContinueCmd(t *testing.T) {
	t.Run("continue - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		testCmd := &cobra.Command{
			Use:   continueCmd.Use,
			Short: continueCmd.Short,
			Run:   continueCmd.Run,
		}
		testCmd.SetOut(&buf)
		testCmd.SetArgs(args)

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
