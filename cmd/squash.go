package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/gh"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/spf13/cobra"
)

var squashCmd = &cobra.Command{
	Use:   "squash",
	Short: "Squash all commits on a branch",
	RunE:  Squash,
}

func Squash(cmd *cobra.Command, args []string) error {
	prInfo, err := gh.GetPRInfo("")
	if err != nil {
		return fmt.Errorf("getting pull request for current branch: %w", err)
	}

	// Find base branch to squash against
	base := types.LocalBranch(prInfo.BaseBranch)
	err = common.Pull(base)
	if err != nil {
		fmt.Printf("WARNING: couldn't pull branch %q: %s\n", base, err)
	}

	currentBranch, err := git.CurrentBranch(path.CurrentDir)
	if err != nil {
		return fmt.Errorf("couldn't get current branch: %w", err)
	}

	// Create new temp branch based on base branch
	squashBranch := squashBranchName(currentBranch)
	err = git.CreateBranch(squashBranch, base)
	if err != nil {
		return fmt.Errorf("couldn't create new branch: %w", err)
	}
	// Cleanup branch later
	defer func() {
		err = git.DeleteBranch(squashBranch, false)
		if err != nil {
			fmt.Printf("WARNING: couldn't delete branch %q: %s\n", squashBranch, err)
		}
	}()

	// Switch to new temp branch
	err = git.Switch(squashBranch)
	if err != nil {
		return fmt.Errorf("switching to branch %q: %w", squashBranch, err)
	}

	// git merge --squash <old-branch>
	err = git.Merge(git.MergeArgs{
		Branch: currentBranch,
		Squash: true,
	})
	if err != nil {
		return fmt.Errorf("merging branch %q into %q: %w", currentBranch, squashBranch, err)
	}

	// Commit all changes
	// TODO: use the -c / --reedit-message option to edit existing message
	err = git.Commit(git.CommitArgs{})
	if err != nil {
		return fmt.Errorf("committing changes to branch %q: %w", squashBranch, err)
	}

	// Reset the old branch to the new branch's HEAD
	err = git.Switch(currentBranch)
	if err != nil {
		return fmt.Errorf("switching to branch %q: %w", currentBranch, err)
	}
	err = git.Reset(git.ResetArgs{
		Branch: squashBranch,
		Mode:   git.HardReset,
	})
	if err != nil {
		return fmt.Errorf("resetting branch %q to %q: %w", currentBranch, squashBranch, err)
	}

	return nil
}

// squashBranchName returns the name to use for the temporary branch used for
// squashing.
func squashBranchName(currentBranch types.LocalBranch) types.LocalBranch {
	return "jit/squash/" + currentBranch
}
