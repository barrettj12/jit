package cmd

import (
	"github.com/barrettj12/jit/common"
	"github.com/spf13/cobra"
	"path/filepath"
	"runtime"
)

var whereCmd = &cobra.Command{
	Use:   "where [push/pull] <branch>",
	Short: "Print push/pull destinations for a branch",
	RunE:  Where,
}

func Where(cmd *cobra.Command, args []string) error {
	// Get source directory
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine source directory")
	}
	srcDir := filepath.Dir(file)

	return common.Execute(filepath.Join(srcDir, "cmd/git-where"), args...)
}
