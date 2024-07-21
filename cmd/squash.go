package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common/gh"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

	// Find base commit to rebase against
	base := types.LocalBranch(prInfo.BaseBranch)
	mergeBase, err := git.MergeBase(types.HEAD, base)
	if err != nil {
		return fmt.Errorf("finding merge base: %w", err)
	}

	fmt.Printf("squashing against %q commit %s...\n", base, mergeBase[:10])
	err = git.Rebase(git.RebaseArgs{
		Base:        mergeBase,
		Interactive: true,
		Env:         []string{"GIT_SEQUENCE_EDITOR=jit squash-editor"},
	})
	if err != nil {
		return fmt.Errorf("rebasing: %w", err)
	}

	fmt.Println("successfully squashed")
	return nil
}

var squashEditorCmd = &cobra.Command{
	Use:    "squash-editor",
	Hidden: true,
	RunE:   SquashEditor,
}

func SquashEditor(cmd *cobra.Command, args []string) error {
	filename := args[0]

	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file %q: %w", filename, err)
	}
	lines := strings.Split(string(fileContents), "\n")

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("opening file %q: %w", filename, err)
	}

	// First line (commit) unchanged
	_, err = file.WriteString(lines[0] + "\n")
	if err != nil {
		return fmt.Errorf("writing to file %q: %w", filename, err)
	}

	for _, line := range lines[1:] {
		if strings.HasPrefix(line, "pick ") {
			cutLine, _ := strings.CutPrefix(line, "pick ")
			line = "f " + cutLine
		}

		_, err = file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("writing to file %q: %w", filename, err)
		}
	}
	return nil
}
