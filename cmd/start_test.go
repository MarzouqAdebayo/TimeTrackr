package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestStartCmd(t *testing.T) {
	t.Run("start - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		want := "Task " + args[0] + " started"

		testCmd := &cobra.Command{
			Use:   startCmd.Use,
			Short: startCmd.Short,
			Run:   startCmd.Run,
		}
		testCmd.SetOut(&buf)
		testCmd.SetArgs(args)

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		got := buf.String()

		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
}
