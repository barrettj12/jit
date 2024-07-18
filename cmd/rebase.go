package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/gh"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "Rebase against the latest version of the target branch",
	RunE:  Rebase,
}

func Rebase(cmd *cobra.Command, args []string) error {
	prInfo, err := gh.GetPRInfo("")
	if err != nil {
		return fmt.Errorf("getting pull request for current branch: %w", err)
	}

	// Pull base branch
	base := prInfo.BaseBranch
	fmt.Printf("Pulling branch %q...\n", base)
	err = common.Pull(base)
	if err != nil {
		return fmt.Errorf("pulling branch %q: %w", base, err)
	}

	fmt.Printf("Rebasing against latest version of %q...\n", base)
	err = git.Rebase(git.RebaseArgs{
		Base:        base,
		Interactive: false,
	})
	if err != nil {
		return fmt.Errorf("rebasing: %w", err)
	}

	fmt.Println("Successfully rebased")
	return nil
}
