package git

import (
	"fmt"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"strings"
)

type AddWorktreeArgs struct {
	Dir          path.Dir          // directory to run the command in
	WorktreePath path.Worktree     // path for the new worktree
	Branch       types.LocalBranch // branch to check out in the new worktree (optional)
}

func AddWorktree(opts AddWorktreeArgs) error {
	args := []string{"worktree", "add", opts.WorktreePath.Path()}
	if opts.Branch != "" {
		args = append(args, string(opts.Branch))
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
	Path   path.Worktree
	HEAD   types.LocalBranch
	Branch types.LocalBranch
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
			Path: path.Worktree(worktreePath),
		})
		s++

		headLine := split[s]
		if headLine == "bare" {
			// No more information here
		} else {
			head, _ := strings.CutPrefix(headLine, "HEAD ")
			worktrees[w].HEAD = types.LocalBranch(head)
			s++

			branchLine := split[s]
			ref, _ := strings.CutPrefix(branchLine, "branch ")
			branch, _ := strings.CutPrefix(ref, "refs/heads/")
			worktrees[w].Branch = types.LocalBranch(branch)
		}

		s += 2
		w++
	}

	return worktrees, nil
}
