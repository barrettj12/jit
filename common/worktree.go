package common

import (
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"path/filepath"
	"strings"
)

// AddWorktree adds a worktree for the given branch. It assumes the branch
// already exists and is ready to be checked out.
func AddWorktree(repoBasePath, branch string) error {
	worktreePath := worktreePathForBranchName(branch)
	if worktreePath != branch {
		fmt.Printf("WARNING branch name %q contains slashes, worktree path will be %q instead", branch, worktreePath)
	}

	if repoBasePath == "" {
		var err error
		repoBasePath, err = RepoBasePath()
		if err != nil {
			return fmt.Errorf("couldn't get repo base path: %w", err)
		}
	}

	err := git.AddWorktree(git.AddWorktreeArgs{
		Dir:          repoBasePath,
		WorktreePath: worktreePath,
		Branch:       branch,
	})
	if err != nil {
		return fmt.Errorf("couldn't create worktree for branch %q: %w", branch, err)
	}
	return nil
}

func Pull(branch string) error {
	worktreePath := worktreePathForBranchName(branch)
	repoBasePath, err := RepoBasePath()
	if err != nil {
		return fmt.Errorf("couldn't get repo base path: %w", err)
	}

	err = git.Pull(filepath.Join(repoBasePath, worktreePath))
	if err != nil {
		return err
	}
	return nil
}

// If a branch name contains slashes, the corresponding worktree path should
// have them replaced with underscores.
func worktreePathForBranchName(branchName string) string {
	return strings.ReplaceAll(branchName, "/", "_")
}
