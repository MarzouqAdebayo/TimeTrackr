package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestStatusCmd(t *testing.T) {
	t.Run("status - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		testCmd := &cobra.Command{
			Use:   statusCmd.Use,
			Short: statusCmd.Short,
			Run:   statusCmd.Run,
		}
		testCmd.SetOut(&buf)
		testCmd.SetArgs(args)

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
