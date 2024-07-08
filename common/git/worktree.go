package git

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
	return err
}
