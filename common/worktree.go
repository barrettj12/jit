package common

import (
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"strings"
)

// AddWorktree adds a worktree for the given branch. It assumes the branch
// already exists and is ready to be checked out. It returns a function which
// can be called to open the new worktree for editing.
func AddWorktree(repoBasePath path.Repo, branch types.LocalBranch) (EditFunc, error) {
	worktreeName := defaultWorktreeNameForBranchName(branch)
	if worktreeName != string(branch) {
		fmt.Printf("WARNING branch name %q contains slashes, worktree path will be %q instead\n", branch, worktreeName)
	}

	if repoBasePath == "" {
		var err error
		repoBasePath, err = RepoBasePath()
		if err != nil {
			return nil, fmt.Errorf("couldn't get repo base path: %w", err)
		}
	}

	worktreePath := path.WorktreePath(repoBasePath, worktreeName)
	err := git.AddWorktree(git.AddWorktreeArgs{
		Dir:          repoBasePath,
		WorktreePath: worktreePath,
		Branch:       branch,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't create worktree for branch %q: %w", branch, err)
	}
	fmt.Printf("successfully added worktree %q\n", worktreeName)

	return EditWorktree(worktreePath), nil
}

func Pull(branch types.LocalBranch) error {
	worktreePath, err := LookupWorktreeForBranch(branch)
	if err != nil {
		return fmt.Errorf("getting worktree path for branch %q: %w", branch, err)
	}

	unknownRemote := false
	upstream, err := git.PullTarget(branch)
	if git.IsAmbiguousArgErr(err) {
		// Remote tracking branch doesn't exist. We'll try to pull anyway with
		// no args, but if this fails, we should tell the user to fetch the
		// remote-tracking branch.
		unknownRemote = true
	}
	if err != nil {
		fmt.Printf("WARNING: unable to get upstream for branch %q\n", branch)
		// Try to do a 'pull' anyway with no args
		upstream = types.NoRemote
	}

	err = git.Pull(git.PullArgs{
		LocalBranch:  branch,
		RemoteBranch: upstream,
		Dir:          worktreePath,
	})
	if err != nil && unknownRemote {
		return fmt.Errorf(`
unknown remote for branch %q. Please run
    jit fetch <remote> %s
to fetch the remote-tracking branch.`[1:], branch, branch)
	}
	return err
}

func LookupWorktreeForBranch(branch types.LocalBranch) (path.Worktree, error) {
	// Get list of worktrees
	worktrees, err := git.Worktrees()
	if err != nil {
		return "", fmt.Errorf("getting worktrees: %w", err)
	}

	// Find worktree corresponding to requested branch
	var worktreePath path.Worktree
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

func EditWorktree(path path.Worktree) EditFunc {
	editor := defaultEditor()
	return func() error {
		res := Exec(ExecArgs{
			Cmd:        editor,
			Args:       []string{path.Path()},
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
func defaultWorktreeNameForBranchName(branchName types.LocalBranch) string {
	return strings.ReplaceAll(string(branchName), "/", "_")
}
