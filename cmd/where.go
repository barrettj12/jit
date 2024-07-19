package cmd

import (
	"errors"
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/types"
	"github.com/spf13/cobra"
)

var whereDocs = `
To show the push destination for a branch, run

    jit where push <branch>

To show the pull source for a branch, run

    jit where pull <branch>

If you omit <branch>, Jit will default to the current branch.
`[1:]

var whereCmd = &cobra.Command{
	Use:   "where [push/pull] <branch>",
	Short: "Print push/pull destinations for a branch",
	Long:  whereDocs,
	RunE:  Where,
}

func Where(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New(`must provide "push" or "pull" as first argument`)
	}
	action := args[0]

	var branch types.LocalBranch
	if len(args) > 1 {
		branch = types.LocalBranch(args[1])
	} // No arg = current branch should be handled by Git methods

	var target types.RemoteBranch
	var err error
	switch action {
	case "push":
		target, err = git.PushTarget(branch)
	case "pull":
		target, err = git.PullTarget(branch)
	default:
		return fmt.Errorf(`first argument must be "push" or "pull", not %q`, args[0])
	}
	if err != nil {
		return fmt.Errorf("getting %s target: %w", action, err)
	}

	if target == types.NoRemote {
		fmt.Printf("no %s target configured for ", action)
		if branch == "" {
			fmt.Println("current branch")
		} else {
			fmt.Printf("branch %q\n", branch)
		}
	} else {
		fmt.Println(target)
	}
	return nil
}
