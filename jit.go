package main

import (
	"errors"
	"github.com/barrettj12/jit/cmd"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func main() {
	// Set up global flags
	baseCmd.PersistentFlags().StringP("repo", "R", "", "repo to execute commands in")

	cmd.AddSubcommands(baseCmd)

	err := baseCmd.Execute()
	if err != nil {
		// Try to extract exit code from error
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		// Otherwise, just exit with a code of 1.
		os.Exit(1)
	}
}

// baseCmd represents the base command when called without any subcommands
var baseCmd = &cobra.Command{
	Use:          "jit",
	SilenceUsage: true,
}
