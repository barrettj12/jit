package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
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
	var branchName types.LocalBranch
	var edit func() error
	var err error

	if len(args) > 1 {
		// New branch based on existing
		newBranch := types.LocalBranch(args[0])
		base := types.LocalBranch(args[1])
		branchName, edit, err = newWorktreeNewBranchWithBase(newBranch, base)
	} else if strings.Contains(args[0], ":") {
		branch := types.ParseGitHubBranch(args[0])
		branchName, edit, err = newWorktreeBasedOnGitHubBranch(branch)
	} else {
		branch := types.LocalBranch(args[0])
		branchName, edit, err = newWorktreeBasedOnExistingBranch(branch)
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

func newWorktreeBasedOnExistingBranch(branch types.LocalBranch) (types.LocalBranch, common.EditFunc, error) {
	fmt.Printf("creating new worktree based on existing local branch %q\n", branch)
	edit, err := common.AddWorktree("", branch)
	return branch, edit, err
}

func newWorktreeNewBranchWithBase(newBranch, base types.LocalBranch) (types.LocalBranch, common.EditFunc, error) {
	fmt.Printf("Pulling branch %q...\n", base)
	err := common.Pull(base)
	if err != nil {
		fmt.Printf("WARNING failed to pull branch %q: %v\n", base, err)
	}

	fmt.Printf("creating new branch %q based on %q\n", newBranch, base)
	err = git.CreateBranch(newBranch, base)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't create branch %q: %w", newBranch, err)
	}

	edit, err := common.AddWorktree("", newBranch)
	return newBranch, edit, err
}

func newWorktreeBasedOnGitHubBranch(gitHubBranch types.GitHubBranch) (types.LocalBranch, common.EditFunc, error) {
	repoBasePath, err := common.RepoBasePath()
	if err != nil {
		return "", nil, fmt.Errorf("couldn't get repo base path: %w", err)
	}

	fmt.Printf("creating new worktree based on remote branch %q\n", gitHubBranch)

	// Add remote if it doesn't exist
	// TODO: possible that remote name != username, we should check the list of Git remotes
	remote := types.RemoteName(gitHubBranch.RepoURL.Owner())
	remoteExists, err := git.RemoteExists(repoBasePath, remote)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't calculate if remote %q exists: %w", remote, err)
	}
	if !remoteExists {
		err = addRemote(remote, url.Nil)
		if err != nil {
			return "", nil, fmt.Errorf("couldn't add remote %q: %w", remote, err)
		}
	}

	// Fetch the remote branch
	remoteBranch := types.RemoteBranch{
		Remote: remote,
		Branch: gitHubBranch.Branch,
	}
	err = git.Fetch(path.CurrentDir, remoteBranch)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't fetch remote branch %q: %w", remoteBranch, err)
	}

	// Create new branch
	// This step also sets the upstream, since we are basing it on a remote branch.
	localBranch := types.LocalBranch(remoteBranch.Branch)
	err = git.CreateBranch(localBranch, remoteBranch.AsLocalBranch())
	if err != nil {
		return "", nil, fmt.Errorf("couldn't create branch %q: %w", localBranch, err)
	}

	edit, err := common.AddWorktree("", localBranch)
	return localBranch, edit, err
}
