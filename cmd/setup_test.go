package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestSetupCmd(t *testing.T) {
	t.Run("setup - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		testCmd := &cobra.Command{
			Use:   setupCmd.Use,
			Short: setupCmd.Short,
			Run:   setupCmd.Run,
		}
		testCmd.SetOut(&buf)
		testCmd.SetArgs(args)

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
