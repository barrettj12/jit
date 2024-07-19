package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
	"github.com/spf13/cobra"
	"strconv"
)

var cloneDocs = `
Clone a repo from GitHub. The following formats are all equivalent, and will
clone the GitHub repository located at https://github.com/<user>/<repo>:

    jit clone <user> <repo>
    jit clone <user>/<repo>
    jit clone https://github.com/<user>/<repo>

The clone will be placed in $JIT_DIR/<user>/<repo>. It will be set up as a bare
repository, with an initial worktree creating tracking the default branch.

Use the --fork flag to specify whether you would like to create a fork of the
repo (requires 'gh', the GitHub CLI, to be installed). If not specified, you
will be prompted after cloning.
`[1:]

func newCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <user>/<repo>",
		Short: "Clone a repo from GitHub",
		Long:  cloneDocs,
		RunE:  Clone,
	}

	// Set flags
	cmd.Flags().String("fork", "", "whether to create a fork")
	cmd.Flags().Bool("no-edit", false, "don't open new repo for editing")

	return cmd
}

// Clone clones the provided repo, using the workflow described in
// https://morgan.cugerone.com/blog/how-to-use-git-worktree-and-in-a-clean-way/
func Clone(cmd *cobra.Command, args []string) error {
	githubRepo := url.GitHubURL(args...)
	user := githubRepo.Owner()
	if user == "" {
		return fmt.Errorf("must specify a user to clone repo from")
	}
	repo := githubRepo.RepoName()
	if repo == "" {
		return fmt.Errorf("must specify a repo to clone")
	}

	// Use JIT_DIR to find clone path
	cloneDir, err := common.DefaultRepoBasePath(user, repo)
	if err != nil {
		return fmt.Errorf("getting clone path: %w", err)
	}

	// Clone the repo
	remote := types.RemoteName(user)
	err = git.Clone(git.CloneArgs{
		Repo:       githubRepo,
		CloneDir:   path.GitFolderPath(cloneDir),
		Bare:       true,
		OriginName: remote,
	})
	if err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	// Success - print message to user
	fmt.Printf(`
Successfully cloned repo %s/%s into %v
Create new branches using
    jit new <branch> [<remote>/]<base>
`[1:], user, repo, cloneDir)

	// Set up correct fetch config for remote
	err = git.SetConfig(cloneDir,
		fmt.Sprintf("remote.%s.fetch", remote),
		fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", remote),
	)
	if err != nil {
		return fmt.Errorf("failed to set fetch config for remote %q: %w", user, err)
	}

	// Fork repo and add as remote
	forkFlagVal, err := cmd.Flags().GetString("fork")
	if err != nil {
		return fmt.Errorf("internal error: couldn't get value of --fork flag: %w", err)
	}

	var shouldFork bool
	if forkFlagVal == "" {
		// The user did not specify when typing the command whether we should
		// fork the repo or not. Ask them.
		shouldFork, err = confirm("Create a fork")
		if err != nil {
			return err
		}
	} else {
		shouldFork, err = strconv.ParseBool(forkFlagVal)
		if err != nil {
			return fmt.Errorf("couldn't parse value %q of --fork flag: %w", forkFlagVal, err)
		}
	}

	if shouldFork {
		err = fork(cloneDir, user, repo)
		if err != nil {
			return err
		}
	}

	// Create new worktree tracking HEAD of source remote
	currentBranch, err := git.CurrentBranch(cloneDir)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	edit, err := common.AddWorktree(cloneDir, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to create initial worktree: %w", err)
	}

	// Set upstream for default branch
	remoteBranchName := string(currentBranch)
	remoteBranch := types.RemoteBranch{
		Remote: remote,
		Branch: remoteBranchName,
	}
	err = git.Fetch(cloneDir, remoteBranch)
	if err != nil {
		fmt.Printf("WARNING could not fetch remote branch %s/%s: %v\n", remote, currentBranch, err)
	}
	err = git.SetUpstream(git.SetUpstreamArgs{
		LocalBranch:  currentBranch,
		RemoteBranch: remoteBranch,
		Dir:          cloneDir,
	})
	if err != nil {
		fmt.Printf("WARNING could not set remote for branch %q: %v\n", currentBranch, err)
	}

	fmt.Printf("created initial worktree %s\n", currentBranch)

	// Open new branch for editing (maybe)
	noEdit, err := cmd.Flags().GetBool("no-edit")
	if err != nil {
		fmt.Printf("WARNING could not get value of --no-edit flag, opening for editing anyway\n")
		noEdit = false
	}
	if !noEdit {
		err = edit()
		if err != nil {
			fmt.Printf("WARNING could not open branch %q for editing: %v\n", currentBranch, err)
		}
	}

	return nil
}
