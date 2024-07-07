package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
	"strings"
)

var newDocs = `
Create a new worktree or branch. This command operates in three modes.

    jit new <branch>

checks out an existing branch named <branch>, in a new worktree, which will
also be named <branch>.

    jit new <branch> <base>

creates a new branch <branch> based on <base>, and checks out <branch> in a new
worktree named <branch>.

    jit new <remote>:<branch>

checks out <branch> from the remote <remote> in a new worktree called <branch>.
If the remote does not already exist, it will automatically be created.
`[1:]

var newCmd = &cobra.Command{
	Use:   "new <branch> [base]",
	Short: "Create a new branch",
	Long:  newDocs,
	RunE:  New,
}

func New(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return newWorktreeNewBranchWithBase(args[0], args[1])
	}
	if strings.Contains(args[0], ":") {
		split := strings.SplitN(args[0], ":", 2)
		return newWorktreeBasedOnRemoteBranch(split[0], split[1])
	}
	return newWorktreeBasedOnExistingBranch(args[0])
}

func newWorktreeBasedOnExistingBranch(branch string) error {
	repoBasePath, err := common.RepoBasePath()
	if err != nil {
		return fmt.Errorf("couldn't get repo base path: %w", err)
	}

	fmt.Printf("creating new worktree based on existing local branch %q\n", branch)

	worktreePath := worktreePathForBranchName(branch)
	if worktreePath != branch {
		fmt.Printf("WARNING branch name %q contains slashes, worktree path will be %q instead", branch, worktreePath)
	}

	err = git.AddWorktree(git.AddWorktreeArgs{
		Dir:          repoBasePath,
		WorktreePath: worktreePath,
		Branch:       branch,
	})
	if err != nil {
		return fmt.Errorf("couldn't create worktree for branch %q: %w", branch, err)
	}

	// We can't set the upstream branch yet, as it doesn't yet exist. We have
	// to wait until push time to do this.
	return nil
}

func newWorktreeNewBranchWithBase(newBranch, base string) error {
	/// TODO: pull base branch

	fmt.Printf("creating new branch %q based on %q\n", newBranch, base)
	err := git.CreateBranch(newBranch, base)
	if err != nil {
		return fmt.Errorf("couldn't create branch %q: %w", newBranch, err)
	}

	repoBasePath, err := common.RepoBasePath()
	if err != nil {
		return fmt.Errorf("couldn't get repo base path: %w", err)
	}

	worktreePath := worktreePathForBranchName(newBranch)
	if worktreePath != newBranch {
		fmt.Printf("WARNING branch name %q contains slashes, worktree path will be %q instead", newBranch, worktreePath)
	}

	err = git.AddWorktree(git.AddWorktreeArgs{
		Dir:          repoBasePath,
		WorktreePath: worktreePath,
		Branch:       newBranch,
	})
	if err != nil {
		return fmt.Errorf("couldn't create worktree for branch %q: %w", newBranch, err)
	}

	// We can't set the upstream branch yet, as it doesn't yet exist. We have
	// to wait until push time to do this.
	return nil
}

func newWorktreeBasedOnRemoteBranch(remote, branch string) error {
	repoBasePath, err := common.RepoBasePath()
	if err != nil {
		return fmt.Errorf("couldn't get repo base path: %w", err)
	}

	fmt.Printf("creating new worktree based on remote branch %s:%s\n", remote, branch)

	// Add remote if it doesn't exist
	remoteExists, err := git.RemoteExists(repoBasePath, remote)
	if err != nil {
		return fmt.Errorf("couldn't calculate if remote %q exists: %w", remote, err)
	}
	if !remoteExists {
		err = addRemote(remote, "")
		if err != nil {
			return fmt.Errorf("couldn't add remote %q: %w", remote, err)
		}
	}

	// Fetch the remote branch
	err = git.Fetch(remote, branch)
	if err != nil {
		return fmt.Errorf("couldn't fetch remote branch %s:%s: %w", remote, branch, err)
	}

	// Create new branch
	// This step also sets the upstream, since we are basing it on a remote branch.
	err = git.CreateBranch(branch, fmt.Sprintf("%s/%s", remote, branch))
	if err != nil {
		return fmt.Errorf("couldn't create branch %q: %w", branch, err)
	}

	worktreePath := worktreePathForBranchName(branch)
	if worktreePath != branch {
		fmt.Printf("WARNING branch name %q contains slashes, worktree path will be %q instead", branch, worktreePath)
	}

	err = git.AddWorktree(git.AddWorktreeArgs{
		Dir:          repoBasePath,
		WorktreePath: worktreePath,
		Branch:       branch,
	})
	if err != nil {
		return fmt.Errorf("couldn't create worktree for branch %q: %w", branch, err)
	}
	return nil
}

// If a branch name contains slashes, the corresponding worktree path should
// have them replaced with underscores.
func worktreePathForBranchName(branchName string) string {
	return strings.ReplaceAll(branchName, "/", "_")
}
