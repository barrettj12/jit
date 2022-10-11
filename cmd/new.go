package cmd

import (
	"path/filepath"

	"github.com/barrettj12/jit/common"
)

// For an existing branch "foo", need to run
//
//	git worktree add <repo>/foo
//
// For a new branch "bar" based on "foo", need to run
//
//	git worktree add <repo>/bar foo -b bar
func New(args []string) error {
	newB, err := common.ReqArg(args, 0, "Enter a name for the new branch:")
	if err != nil {
		return err
	}
	base, err := common.ReqArg(args, 1, "Which branch should this be based on?")
	if err != nil {
		return err
	}

	// Get path to new worktree
	gitDir, err := common.RepoBasePath()
	if err != nil {
		return err
	}
	path := filepath.Join(gitDir, newB)

	// Create worktree
	return common.Git("worktree", []string{"add", path, base, "-b", newB})
}
