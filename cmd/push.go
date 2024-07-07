package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push up the current branch",
	RunE:  Push,
}

func Push(cmd *cobra.Command, args []string) error {
	pushArgs := git.PushArgs{}

	// Check if push destination is set up
	pushTarget, err := git.PushTarget("")
	if err != nil {
		return fmt.Errorf("calculating push destination: %w", err)
	}
	if pushTarget == "" {
		// No upstream is set, so we need to set it in the push.
		remote := common.GitHubUser()
		if remote == "" {
			return fmt.Errorf(`can't determine push destination for current branch.
Please set the GH_USER env var or set the upstream manually.`)
		}
		pushArgs.Remote = remote

		currentBranch, err := git.CurrentBranch()
		if err != nil {
			return fmt.Errorf("can't get current branch: %w", err)
		}
		pushArgs.Branch = currentBranch
		pushArgs.SetUpstream = true
	}

	err = git.Push(pushArgs)
	if err != nil {
		return fmt.Errorf("couldn't push: %w", err)
	}
	return nil
}