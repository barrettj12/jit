package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <branch>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a branch from remote & local",
	RunE:    Remove,
}

func Remove(cmd *cobra.Command, args []string) error {
	var remoteBranch types.GitHubBranch
	var worktree path.Worktree
	var localBranch types.LocalBranch

	deleteRemoteBranch(remoteBranch)
	deleteWorktree(worktree)
	deleteLocalBranch(localBranch)
}

func deleteRemoteBranch(branch types.GitHubBranch) error {
	remote, err := common.GitRemoteFromURL(branch.RepoURL)
	if err != nil {
		fmt.Printf("WARNING: couldn't find remote matching %q: %v\n", branch.RepoURL.Owner(), err)
		fmt.Printf("assuming remote name is %q\n", branch.RepoURL.Owner())
		remote = types.RemoteName(branch.RepoURL.Owner())
	}

	// TODO: lookup local branch that's tracking the given remote
	localBranch := types.LocalBranch(branch.Branch)

	err = git.Push(git.PushArgs{
		Branch: localBranch,
		Remote: remote,
		Delete: true,
	})
	if err != nil {
		return fmt.Errorf("deleting remote branch %q: %w", branch, err)
	}
	return nil
}

func deleteWorktree(worktree path.Worktree) {

}

func deleteLocalBranch(branch types.LocalBranch) {

}
