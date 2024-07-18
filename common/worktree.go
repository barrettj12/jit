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
	worktreeName := defaultWorktreeNameForBranchName(branch)
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

	return EditWorktree(filepath.Join(repoBasePath, worktreeName)), nil
}

func Pull(branch string) error {
	worktreePath, err := LookupWorktreeForBranch(branch)
	if err != nil {
		return fmt.Errorf("getting worktree path for branch %q: %w", branch, err)
	}

	upstream, err := git.PullTarget(branch)
	if err != nil {
		return fmt.Errorf("getting upstream for branch %q: %w", branch, err)
	}
	split := strings.SplitN(upstream, "/", 2)

	err = git.Pull(git.PullArgs{
		Remote: split[0],
		Branch: split[1],
		Dir:    worktreePath,
	})
	if err != nil {
		return err
	}
	return nil
}

func LookupWorktreeForBranch(branch string) (string, error) {
	// Get list of worktrees
	worktrees, err := git.Worktrees()
	if err != nil {
		return "", fmt.Errorf("getting worktrees: %w", err)
	}

	// Find worktree corresponding to requested branch
	var worktreePath string
	for _, worktree := range worktrees {
		if worktree.Branch == branch {
			worktreePath = worktree.Path
			break
		}
	}
	if worktreePath == "" {
		return "", fmt.Errorf("no worktree found for branch %q", branch)
	}
	return worktreePath, nil
}

type EditFunc func() error

func EditWorktree(worktreePath string) EditFunc {
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
func defaultWorktreeNameForBranchName(branchName string) string {
	return strings.ReplaceAll(branchName, "/", "_")
}
