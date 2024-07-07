package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
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

	return cmd
}

// Clone clones the provided repo, using the workflow described in
// https://morgan.cugerone.com/blog/how-to-use-git-worktree-and-in-a-clean-way/
func Clone(cmd *cobra.Command, args []string) error {
	var user, repo string

	arg1, err := common.ReqArg(args, 0, "URL of Git repo to clone:")
	if err != nil {
		return err
	}
	user, repo, err = parseGitRepoURL(arg1)
	if err != nil {
		return err
	}

	if repo == "" {
		// arg1 was just "user"
		repo, err = common.ReqArg(args, 1, fmt.Sprintf("Which of %s's repos do you want?", user))
		if err != nil {
			return err
		}
	}

	// Reconstruct repository URL
	repoURL := githubURL(user, repo)

	// Use JIT_DIR to find clone path
	jitDir, err := common.JitDir()
	if err != nil {
		return err
	}
	cloneDir := filepath.Join(jitDir, user, repo)

	// Clone the repo
	err = git.Clone(git.CloneArgs{
		RepoURL:    repoURL,
		CloneDir:   filepath.Join(cloneDir, ".git"),
		Bare:       true,
		OriginName: user,
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
		err = fork(user, repo, cloneDir)
		if err != nil {
			return err
		}
	}

	// Create new worktree tracking HEAD of source remote
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	// TODO: this should use common code with the new command so that we are
	//   correctly setting up the upstream branch, etc.
	err = git.AddWorktree(cloneDir, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to create initial worktree: %w", err)
	}

	fmt.Printf("created initial worktree %s\n", currentBranch)
	return nil
}

// Parses the given argument to a GitHub user & repo.
// Valid inputs are
//
//	"user"                         -> "user", "",     nil
//	"user/repo"                    -> "user", "repo", nil
//	"https://github.com/user/repo" -> "user", "repo", nil
func parseGitRepoURL(raw string) (string, string, error) {
	u, err := url.Parse(raw) // only matches full URL with scheme
	if err != nil {
		return "", "", err
	}

	switch u.Host {
	case "":
	// raw is not a URL
	case "github.com":
		raw = u.Path[1:] // strip leading '/'
	default:
		return "", "", fmt.Errorf("host %s not supported", u.Host)
	}

	// raw is now "user" or "user/repo"
	parts := strings.Split(raw, "/")
	switch len(parts) {
	case 1:
		// "user"
		return parts[0], "", nil
	case 2:
		// "user/repo"
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("invalid Git repo URL %s", raw)
	}
}

func githubURL(user, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s", user, repo)
}
