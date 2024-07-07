package cmd

import (
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

	// TODO: this needs to use the methods in common/worktree.go
	path, err := common.WorktreePath(branch)
	if err != nil {
		return err
	}

	editor := defaultEditor()

	res := common.Exec(common.ExecArgs{
		Cmd:        editor,
		Args:       []string{path},
		Background: true,
	})
	return res.RunError
}

func defaultEditor() string {
	// TODO: allow specifying default editor on a per-repo basis
	return "goland"
}
