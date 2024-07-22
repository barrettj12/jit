package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
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
	if pushTarget == types.NoRemote {
		// No upstream is set, so we need to set it in the push.
		remote, err := common.DefaultRemote()
		if err != nil {
			return fmt.Errorf(`can't determine push destination for current branch.
Please set the GH_USER env var or set the upstream manually.`)
		}

		// Add remote if it doesn't already exist
		err = addRemote(remote, url.Nil)
		if err != nil && !git.IsRemoteAlreadyExistsErr(err) {
			return fmt.Errorf("couldn't add remote %q: %w", remote, err)
		}
		pushArgs.Remote = remote

		currentBranch, err := git.CurrentBranch(path.CurrentDir)
		if err != nil {
			return fmt.Errorf("can't get current branch: %w", err)
		}
		pushArgs.Branch = currentBranch
		pushArgs.SetUpstream = true
		pushTarget = types.RemoteBranch{
			Remote: remote,
			Branch: string(currentBranch),
		}
	}

	fmt.Printf("pushing to %q...\n", pushTarget)
	err = git.Push(pushArgs)
	if err != nil {
		return fmt.Errorf("couldn't push: %w", err)
	}
	fmt.Printf("successfully pushed current branch to %q\n", pushTarget)
	return nil
}
