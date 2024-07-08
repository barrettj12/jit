package common

import (
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"path/filepath"
	"strings"
)

// AddWorktree adds a worktree for the given branch. It assumes the branch
// already exists and is ready to be checked out. It returns a function which
// can be called to open the new worktree for editing.
func AddWorktree(repoBasePath, branch string) (EditFunc, error) {
	worktreeName := worktreeNameForBranchName(branch)
	if worktreeName != branch {
		fmt.Printf("WARNING branch name %q contains slashes, worktree path will be %q instead\n", branch, worktreeName)
	}

	if repoBasePath == "" {
		var err error
		repoBasePath, err = RepoBasePath()
		if err != nil {
			return nil, fmt.Errorf("couldn't get repo base path: %w", err)
		}
	}

	err := git.AddWorktree(git.AddWorktreeArgs{
		Dir:          repoBasePath,
		WorktreePath: worktreeName,
		Branch:       branch,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create worktree for branch %q: %w", branch, err)
	}
	fmt.Printf("successfully added worktree %q\n", worktreeName)

	return editWorktree(filepath.Join(repoBasePath, worktreeName)), nil
}

func Pull(branch string) error {
	worktreeName := worktreeNameForBranchName(branch)
	repoBasePath, err := RepoBasePath()
	if err != nil {
		return fmt.Errorf("couldn't get repo base path: %w", err)
	}

	err = git.Pull(filepath.Join(repoBasePath, worktreeName))
	if err != nil {
		return err
	}
	return nil
}

type EditFunc func() error

func EditBranch(branch string) (EditFunc, error) {
	worktreeName := worktreeNameForBranchName(branch)
	repoBasePath, err := RepoBasePath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get repo base path: %w", err)
	}
	return editWorktree(filepath.Join(repoBasePath, worktreeName)), nil
}

func editWorktree(worktreePath string) EditFunc {
	editor := defaultEditor()
	return func() error {
		res := Exec(ExecArgs{
			Cmd:        editor,
			Args:       []string{worktreePath},
			Background: true,
		})
		return res.RunError
	}
}

// TODO: allow configuring this value on a per-repo / global basis
func defaultEditor() string {
	// TODO: allow specifying default editor on a per-repo basis
	return "goland"
}

// If a branch name contains slashes, the corresponding worktree path should
// have them replaced with underscores.
func worktreeNameForBranchName(branchName string) string {
	return strings.ReplaceAll(branchName, "/", "_")
}
