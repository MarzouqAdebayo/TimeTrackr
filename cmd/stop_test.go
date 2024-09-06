package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestStopCmd(t *testing.T) {
	t.Run("stop - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{"task1"}

		want := "Echo: " + args[0]

		testCmd := &cobra.Command{
			Use:   stopCmd.Use,
			Short: stopCmd.Short,
			Run:   stopCmd.Run,
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
