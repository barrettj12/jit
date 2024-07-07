package cmd

import (
	"github.com/barrettj12/jit/common"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log <args>",
	Short: "Print Git log in short format",
	RunE:  Log,
}

func Log(cmd *cobra.Command, args []string) error {
	gitArgs := append([]string{"--oneline"}, args...)
	return common.Git("log", gitArgs...)
}
