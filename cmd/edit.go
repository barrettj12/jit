package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/config"
	"github.com/barrettj12/jit/common/types"
	"github.com/spf13/cobra"
	"strings"
)

var editCmd = &cobra.Command{
	Use:   "edit <branch>",
	Short: "Open a branch for editing",
	RunE:  Edit,
}

// Edit opens the given branch in the default editor.
func Edit(cmd *cobra.Command, args []string) error {
	editor, err := config.Editor()
	if err != nil {
		fmt.Printf("error getting default editor: %v\n", err)
	} else {
		fmt.Printf("default editor is %q\n", editor)
	}
	return nil

	branch, err := common.ReqArg(args, 0, "Which branch would you like to edit?")
	if err != nil {
		return err
	}

	// TODO: handle the different cases for branch
	//   types.LocalBranch, types.GitHubBranch, (types.RemoteBranch?)

	// Strip remote from branch name
	if strings.Contains(branch, ":") {
		split := strings.SplitN(branch, ":", 2)
		branch = split[1]
	}

	worktreePath, err := common.LookupWorktreeForBranch(types.LocalBranch(branch))
	if err != nil {
		return err
	}

	err = common.EditWorktree(worktreePath)()
	if err != nil {
		return fmt.Errorf("couldn't open branch %q for editing: %w", branch, err)
	}
	return nil
}
