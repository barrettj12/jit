package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
)

var pullDocs = `
Pull a given local branch using:

    jit pull <branch>

Inside a worktree, you can pull the current branch by simply running:

    jit pull
`[1:]

var pullCmd = &cobra.Command{
	Use:   "pull <branch>",
	Short: "Pull a remote branch",
	Long:  pullDocs,
	RunE:  Pull,
}

func Pull(cmd *cobra.Command, args []string) error {
	// If a branch was specified, just attempt to pull that branch.
	if len(args) > 0 {
		return pullBranch(args[0])
	}

	// No branch specified. First, let's try to pull the current branch.
	branch, err := git.CurrentBranch("")
	if err == nil {
		return pullBranch(branch)
	}

	// If there is no current branch (e.g. we are not inside a worktree), ask
	// the user which branch they'd like to pull.
	branch, err = common.Prompt("Which branch would you like to pull?")
	if err != nil {
		return err
	}
	return pullBranch(branch)
}

func pullBranch(branch string) error {
	fmt.Printf("pulling branch %q...\n", branch)
	err := common.Pull(branch)
	if err != nil {
		return err
	}

	fmt.Printf("successfully pulled branch %q\n", branch)
	return nil
}
