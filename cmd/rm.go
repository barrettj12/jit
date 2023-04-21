package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/barrettj12/jit/common"
)

func Remove(args []string) error {
	// TODO: add --force option
	branch, err := common.ReqArg(args, 0, "Which branch would you like to remove?")
	if err != nil {
		return err
	}

	// Delete remote tracking branch
	// git push -d <remote_name> <branchname>
	// TODO: not working when upstream is in origin/source repo

	// TODO: this should work if the remote branch doesn't exist (no-op)
	//Delete remote tracking branch barrettj12/rm-webster? [y/n]: y
	//error: unable to delete 'rm-webster': remote ref does not exist
	//error: failed to push some refs to 'https://github.com/barrettj12/interview-questions'
	//ERROR: exit status 1
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
	_, err = os.Stat(wktreePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("no worktree found at %s\n", wktreePath)
	} else if err != nil {
		return err
	} else {
		ok, err := confirm(fmt.Sprintf("Delete worktree at %s", wktreePath))
		if err != nil {
			return err
		}
		if ok {
			err = common.Git("worktree", "remove", wktreePath)
			if err != nil {
				force, err := confirm("Worktree deletion failed, try again with --force")
				if err != nil {
					return err
				}
				if force {
					err = common.Git("worktree", "remove", wktreePath, "--force")
					if err != nil {
						return err
					}
				}
			}
			err = os.RemoveAll(wktreePath)
			if err != nil {
				return err
			}
		}
	}

	// Delete local branch
	// git branch -d <branchname>
	// TODO: git fetch so that we don't get the error message
	//     error: The branch ____ is not fully merged.
	ok, err := confirm(fmt.Sprintf("Delete local branch %q", branch))
	if err != nil {
		return err
	}
	if ok {
		err = common.Git("branch", "-d", branch)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			force, err := confirm("Branch deletion failed, try again with force")
			if err != nil {
				return err
			}
			if force {
				err = common.Git("branch", "-D", branch)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// TODO: move this to common
func confirm(prompt string) (bool, error) {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s? [y/n]: ", prompt)
	sc.Scan()
	if err := sc.Err(); err != nil {
		return false, fmt.Errorf("error reading input: %w", err)
	}
	return sc.Text() == "y", nil
}
