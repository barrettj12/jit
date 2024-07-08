package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <branch>",
	Short: "Open a branch for editing",
	RunE:  Edit,
}

// Edit opens the given branch in the default editor.
func Edit(cmd *cobra.Command, args []string) error {
	branch, err := common.ReqArg(args, 0, "Which branch would you like to edit?")
	if err != nil {
		return err
	}

	edit, err := common.EditBranch(branch)
	if err == nil {
		err = edit()
	}

	if err != nil {
		return fmt.Errorf("couldn't open branch %q for editing: %w", branch, err)
	}
	return nil
}
