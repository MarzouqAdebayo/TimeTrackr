package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCmd(t *testing.T) {
	t.Run("version - Should print appropriate command when passed", func(t *testing.T) {
		buf := bytes.Buffer{}
		args := []string{}

		testCmd := &cobra.Command{
			Use:   versionCmd.Use,
			Short: versionCmd.Short,
			Run:   versionCmd.Run,
		}
		testCmd.SetOut(&buf)
		testCmd.SetArgs(args)

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		got := buf.String()

		if got != versionResponse {
			t.Errorf("want: %s, got: %s", versionResponse, got)
		}
	})
}
