package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestPauseCmd(t *testing.T) {
	t.Run("pause - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		want := "Echo: " + args[0]

		testCmd := &cobra.Command{
			Use:   pauseCmd.Use,
			Short: pauseCmd.Short,
			Run:   pauseCmd.Run,
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
