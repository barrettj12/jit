package cmd

import (
	"os"

	"github.com/barrettj12/jit/common"
)

func Pull(args []string) error {
	branch, err := common.ReqArg(args, 0, "Which branch would you like to pull?")
	if err != nil {
		return err
	}

	var pullArgs []string
	if len(args) >= 1 {
		pullArgs = args[1:]
	}

	return pull(branch, pullArgs...)
}

func pull(branch string, pullArgs ...string) error {
	path, err := common.WorktreePath(branch)
	if err != nil {
		return err
	}

	res := common.Exec(common.ExecArgs{
		Cmd:    "git",
		Args:   append([]string{"pull"}, pullArgs...),
		Dir:    path,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	return res.RunError
}
