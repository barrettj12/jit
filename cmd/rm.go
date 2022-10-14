package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/barrettj12/jit/common"
)

func Remove(args []string) error {
	branch, err := common.ReqArg(args, 0, "Which branch would you like to remove?")
	if err != nil {
		return err
	}

	// Delete remote tracking branch
	// git push -d <remote_name> <branchname>
	remote, remoteBranch, err := common.PushLoc(branch)
	if err == common.ErrUpstreamNotFound {
		// no-op
		fmt.Printf("no remote tracking branch found for branch %q\n", branch)
	} else {
		if err != nil {
			return err
		}
		ok, err := confirm(fmt.Sprintf("Delete remote tracking branch %s/%s", remote, remoteBranch))
		if err != nil {
			return err
		}
		if ok {
			err = common.Git("push", "-d", remote, remoteBranch)
			if err != nil {
				return err
			}
		}
	}

	// Delete worktree
	wktreePath, err := common.WorktreePath(branch)
	if err != nil {
		return err
	}
	ok, err := confirm(fmt.Sprintf("Delete worktree at %s", wktreePath))
	if err != nil {
		return err
	}
	if ok {
		err = common.Git("worktree", "remove", wktreePath)
		if err != nil {
			return err
		}
		err = common.Execute("rm", "-r", wktreePath)
		if err != nil {
			return err
		}
	}

	// Delete local branch
	// git branch -d <branchname>
	ok, err = confirm(fmt.Sprintf("Delete local branch %q", branch))
	if err != nil {
		return err
	}
	if ok {
		err = common.Git("branch", "-d", branch)
		if err != nil {
			return err
		}
	}

	return nil
}

func confirm(prompt string) (bool, error) {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s? [y/n]: ", prompt)
	sc.Scan()
	if err := sc.Err(); err != nil {
		return false, fmt.Errorf("error reading input: %w", err)
	}
	return sc.Text() == "y", nil
}
