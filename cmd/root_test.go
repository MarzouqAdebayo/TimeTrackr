package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd(t *testing.T) {
	t.Run("root - Should print appropriate command when passed", func(t *testing.T) {
		want := "Welcome to TimeTrackr, use trackr --help to get started"
		buf := bytes.Buffer{}
		rootCmd.SetOut(&buf)
		rootCmd.SetArgs([]string{})

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		got := buf.String()

		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
}
