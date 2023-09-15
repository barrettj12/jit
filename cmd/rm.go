package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/barrettj12/jit/common"
)

func Remove(args []string) error {
	// TODO: add --force option
	branch, err := common.ReqArg(args, 0, "Which branch would you like to remove?")
	if err != nil {
		return err
	}

	split := strings.SplitN(branch, ":", 2)
	if len(split) >= 2 {
		branch = split[1]
	}

	// Delete remote tracking branch
	// git push -d <remote_name> <branchname>
	// TODO: not working when upstream is in origin/source repo

	// TODO: this should work if the remote branch doesn't exist (no-op)
	//Delete remote tracking branch barrettj12/rm-webster? [y/n]: y
	//error: unable to delete 'rm-webster': remote ref does not exist
	//error: failed to push some refs to 'https://github.com/barrettj12/interview-questions'
	//ERROR: exit status 1

	// TODO: need to be able to handle branches with "/" in the name
	//   $ jit rm imerge/3.3
	//   Delete remote tracking branch barrettj12/imerge? [y/n]
	remote, remoteBranch, err := common.PushLoc(branch)
	switch err {
	case common.ErrUpstreamNotFound:
		// no-op
		fmt.Printf("no remote tracking branch found for branch %q\n", branch)

	case nil:
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

	default: // non-nil error
		return err
	}

	// Delete worktree
	wktreePath, err := common.WorktreePath(branch)
	if err != nil {
		return err
	}
	// TODO: worktrees can be removed but the worktree still exists in Git.
	// Check if it exists in `git worktrees list`, rather than doing a stat.
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
				// Usually the error is "worktree contains modified or untracked files"
				// so print these files for the user to see.

				// TODO: the err could be
				//   fatal: <path> is not a working tree
				// in which case we need to just skip to the branch removal step

				untrackedFiles, _ := common.ExecGit(wktreePath, "status", "--porcelain", "--ignore-submodules=none")
				fmt.Println(untrackedFiles)

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

	// Fetch merged branch, so that we don't get the error message
	//     error: The branch ____ is not fully merged.
	// TODO: get the remote/branch from gh pr view
	_ = common.Fetch("", "")

	// Delete local branch
	// git branch -d <branchname>
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
