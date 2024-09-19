package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestWebServerCmd(t *testing.T) {
	t.Run("start - Should print appropriate command when passed", func(t *testing.T) {
		testCmd := &cobra.Command{
			Use:   webServerCmd.Use,
			Short: webServerCmd.Short,
			Run:   webServerCmd.Run,
		}

		err := testCmd.Execute()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
