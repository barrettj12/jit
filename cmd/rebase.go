package cmd

import (
	"github.com/barrettj12/jit/common"
	"github.com/spf13/cobra"
	"path/filepath"
	"runtime"
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase <branch>",
	Short: "Rebase a branch",
	RunE:  Rebase,
}

func Rebase(cmd *cobra.Command, args []string) error {
	// Get source directory
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine source directory")
	}
	srcDir := filepath.Dir(file)

	return common.Execute(filepath.Join(srcDir, "cmd/git-rebase"), args...)
}
