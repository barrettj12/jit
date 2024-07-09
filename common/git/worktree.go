package git

import (
	"fmt"
	"strings"
)

type AddWorktreeArgs struct {
	Dir          string // directory to run the command in
	WorktreePath string // path for the new worktree
	Branch       string // branch to check out in the new worktree (optional)
}

func AddWorktree(opts AddWorktreeArgs) error {
	args := []string{"worktree", "add", opts.WorktreePath}
	if opts.Branch != "" {
		args = append(args, opts.Branch)
	}
	_, err := internalExec(internalExecArgs{
		args: args,
		dir:  opts.Dir,
	})
	if IsNoSuchBranchErr(err) {
		return fmt.Errorf("branch %q doesn't exist", opts.Branch)
	}
	return err
}

type WorktreeInfo struct {
	Path   string
	HEAD   string
	Branch string
}

// Returns information on all worktrees.
func Worktrees() ([]WorktreeInfo, error) {
	out, err := internalExec(internalExecArgs{
		args: []string{"worktree", "list", "--porcelain", "-z"},
	})
	if err != nil {
		return nil, err
	}

	// Parse output into worktree info list
	split := strings.Split(out, "\x00")
	var worktrees []WorktreeInfo
	var s, w int
	for s < len(split)-1 {
		worktreeLine := split[s]
		worktreePath, _ := strings.CutPrefix(worktreeLine, "worktree ")
		worktrees = append(worktrees, WorktreeInfo{
			Path: worktreePath,
		})
		s++

		headLine := split[s]
		if headLine == "bare" {
			// No more information here
		} else {
			head, _ := strings.CutPrefix(headLine, "HEAD ")
			worktrees[w].HEAD = head
			s++

			branchLine := split[s]
			ref, _ := strings.CutPrefix(branchLine, "branch ")
			branch, _ := strings.CutPrefix(ref, "refs/heads/")
			worktrees[w].Branch = branch
		}

		s += 2
		w++
	}

	return worktrees, nil
}
