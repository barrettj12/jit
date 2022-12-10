package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/barrettj12/jit/common"
)

// For an existing branch "foo", need to run
//
//	git worktree add <repo>/foo
//
// For a new branch "bar" based on "foo", need to run
//
//	git worktree add <repo>/bar foo -b bar
func New(args []string) error {
	branch, err := common.ReqArg(args, 0, "Which branch do you want to create/get?")
	if err != nil {
		return err
	}

	// TODO: if remote branch not found, then try fetch
	//   $ jit new ip-address pengale/main
	//   Preparing worktree (new branch 'ip-address')
	//   fatal: not a valid object name: 'pengale/main'
	//   ERROR: exit status 255

	// TODO: fix this
	//   $ jit new main
	//   Which branch should this be based on? source/main
	//   Preparing worktree (new branch 'main')
	//   fatal: a branch named 'main' already exists
	//   ERROR: exit status 255
	//   $ jit new main main
	//   Preparing worktree (new branch 'main')
	//   fatal: a branch named 'main' already exists
	//   ERROR: exit status 255

	// TODO: should be able to find branch in remote (i.e. no prompt)
	//   $ jit new fix-info-bundle-formatting
	//   Which branch should this be based on? benhoyt/fix-info-bundle-formatting
	//   Preparing worktree (new branch 'fix-info-bundle-formatting')
	//   branch 'fix-info-bundle-formatting' set up to track 'benhoyt/fix-info-bundle-formatting'.
	//   HEAD is now at 675e36bc44 Fix messed-up channels formatting for "juju info" of a bundle

	// cases:
	//   branch exists locally, not checked out
	//    -> create new worktree, check out this branch
	//   branch exists locally, checked out
	//    -> error, branch already exists and checked out
	//   branch exists in a remote
	//    -> create local branch, tracking remote branch
	//   branch doesn't exist
	//    -> create local branch, ask user what branch to base on

	// Check if branch exists locally
	if branchExistsLocally(branch) {
		// Create new worktree using this branch
		path, err := common.WorktreePath(branch)
		if err != nil {
			return err
		}
		return common.Git("worktree", "add", path)
	}

	// Try to find branch in remotes
	branches, err := searchRemotesForBranch(branch)
	if len(branches) > 1 {
		return fmt.Errorf("multiple branches matching %q, please specify remote\n%s\n",
			branch, branches)
	}
	if len(branches) == 1 {
		// Create new local branch tracking remote branch
		remoteRef := branches[0]
		split := strings.Split(remoteRef, "/")
		newBranch := split[len(split)-1]

		path, err := common.WorktreePath(newBranch)
		if err != nil {
			return err
		}
		return common.Git("worktree", "add", path, remoteRef)
	}

	// len(branches) = 0, i.e. branch was not found in remotes.
	// In this case, we create a brand-new branch
	base, err := common.ReqArg(args, 1, "Which branch should this be based on?")
	if err != nil {
		return err
	}

	// Create worktree
	path, err := common.WorktreePath(branch)
	if err != nil {
		return err
	}

	err = common.Git("worktree", "add", path, base, "-b", branch)
	if err != nil {
		return err
	}

	return Edit([]string{branch})
}

func branchExistsLocally(branch string) bool {
	err := common.Git("show-ref", "--quiet", fmt.Sprint("refs/heads/%s", branch))
	return err == nil
}

// If argument is `remote/branch` -> check this exists
// If argument is just `branch` -> search all remotes for this branch
// Return list of all matching `remote/branch`
func searchRemotesForBranch(branch string) ([]string, error) {
	ret := []string{}
	lookBranch := func(string) {
		err := common.Git("show-ref", "--quiet", fmt.Sprint("refs/remotes/%s", branch))
		if err == nil {
			ret = append(ret, branch)
		}
	}

	if strings.Contains(branch, "/") {
		lookBranch(branch)
	} else {
		remotes, err := getRemotes()
		if err != nil {
			return nil, err
		}
		for _, r := range remotes {
			lookBranch(fmt.Sprint("%s/%s", r, branch))
		}
	}

	return ret, nil
}

func getRemotes() ([]string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd := exec.Command("git", "remote")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Read stderr for error info
		errInfo := stderr.String()
		return nil, fmt.Errorf("%s\n%s", errInfo, err)
	}

	remotes := strings.Split(stdout.String(), "\n")
	if remotes[len(remotes)-1] == "" {
		remotes = remotes[:len(remotes)-1]
	}
	return remotes, nil
}
