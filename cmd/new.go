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

func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <branch> [base]",
		Short: "Create a new branch",
		Long:  newDocs,
		RunE:  New,
	}

	// Set flags
	cmd.Flags().Bool("no-edit", false, "don't open new branch for editing")

	return cmd
}

func New(cmd *cobra.Command, args []string) error {
	var branchName string
	var edit func() error
	var err error

	if len(args) > 1 {
		branchName = args[0]
		edit, err = newWorktreeNewBranchWithBase(branchName, args[1])
	} else if strings.Contains(args[0], ":") {
		split := strings.SplitN(args[0], ":", 2)
		branchName = split[1]
		edit, err = newWorktreeBasedOnRemoteBranch(split[0], branchName)
	} else {
		branchName = args[0]
		edit, err = newWorktreeBasedOnExistingBranch(branchName)
	}

	if err != nil {
		return fmt.Errorf("failed to create new branch: %w", err)
	}

	noEdit, err := cmd.Flags().GetBool("no-edit")
	if err != nil {
		fmt.Printf("WARNING could not get value of --no-edit flag, will open new branch for editing anyway\n")
		noEdit = false
	}
	if !noEdit {
		err = edit()
		if err != nil {
			fmt.Printf("WARNING could not open branch %q for editing: %v\n", branchName, err)
		}
	}

	return nil
}

func newWorktreeBasedOnExistingBranch(branch string) (common.EditFunc, error) {
	fmt.Printf("creating new worktree based on existing local branch %q\n", branch)
	return common.AddWorktree("", branch)
}

func newWorktreeNewBranchWithBase(newBranch, base string) (common.EditFunc, error) {
	fmt.Printf("Pulling branch %q...\n", base)
	err := common.Pull(base)
	if err != nil {
		fmt.Printf("WARNING failed to pull branch %q: %v\n", base, err)
	}

	fmt.Printf("creating new branch %q based on %q\n", newBranch, base)
	err = git.CreateBranch(newBranch, base)
	if err != nil {
		return nil, fmt.Errorf("couldn't create branch %q: %w", newBranch, err)
	}

	return common.AddWorktree("", newBranch)
}

func newWorktreeBasedOnRemoteBranch(remote, branch string) (common.EditFunc, error) {
	repoBasePath, err := common.RepoBasePath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get repo base path: %w", err)
	}

	fmt.Printf("creating new worktree based on remote branch %s:%s\n", remote, branch)

	// Add remote if it doesn't exist
	remoteExists, err := git.RemoteExists(repoBasePath, remote)
	if err != nil {
		return nil, fmt.Errorf("couldn't calculate if remote %q exists: %w", remote, err)
	}
	if !remoteExists {
		err = addRemote(remote, "")
		if err != nil {
			return nil, fmt.Errorf("couldn't add remote %q: %w", remote, err)
		}
	}

	// Fetch the remote branch
	err = git.Fetch("", remote, branch)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch remote branch %s:%s: %w", remote, branch, err)
	}

	// Create new branch
	// This step also sets the upstream, since we are basing it on a remote branch.
	err = git.CreateBranch(branch, fmt.Sprintf("%s/%s", remote, branch))
	if err != nil {
		return nil, fmt.Errorf("couldn't create branch %q: %w", branch, err)
	}

	return common.AddWorktree("", branch)
}
