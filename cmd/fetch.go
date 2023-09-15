package cmd

import (
	"fmt"

	"github.com/barrettj12/jit/common"
)

// Accepts the following formats:
//
//	fetch <remote> <branch>
//	fetch <remote>:<branch>
//	fetch <remote>/<branch>
func Fetch(args []string) error {
	git := newGitProvider()

	remote, err := common.ReqArg(args, 0, "What branch would you like to fetch?")
	if err != nil {
		return err
	}

	if len(args) < 2 {
		_, resolved := git.ResolveBranch(remote)
		if resolved {
			return nil
		}
	}

	// Assume remote was just remote name - get branch name
	branch, err := common.ReqArg(args, 0, fmt.Sprintf(
		"What branch would you like to fetch from remote %q?", remote))
	if err != nil {
		return err
	}

	// Add remote if necessary
	_, err = git.GetRemote(remote)
	if err != nil {
		err = git.AddRemote(remote, "")
		if err != nil {
			return fmt.Errorf("could not add remote %q: %w", remote, err)
		}
	}

	// Fetch branch if necessary
	err = git.Fetch(remote, branch)
	if err != nil {
		return fmt.Errorf("could not fetch remote branch %q: %w", branch, err)
	}

	return nil
}
