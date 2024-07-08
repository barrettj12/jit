package cmd

import (
	"fmt"
	"strings"

	"github.com/barrettj12/jit/common"
)

// TODO: all this needs to be replaced with the common/git package

// GitProvider provides Go bindings for git commands.
type GitProvider interface {
	// Branch methods

	// AddBranch creates a new branch with the specified base.
	AddBranch(branch, base string) error

	// ResolveBranch resolves the provided string into an unambiguously specified
	// branch, pulling from remotes if necessary.
	ResolveBranch(string) (string, bool)

	// Worktree methods

	// AddWorktree adds a worktree based on the specified branch.
	AddWorktree(branch string) error

	// Remote methods

	// AddRemote adds the given remote.
	AddRemote(remoteName, url string) error

	// GetRemote retrieves the URL for the given remote.
	GetRemote(remoteName string) (url string, err error)

	// Other methods

	// Fetch fetches the given branch from the given remote.
	Fetch(remote, branch string) error
}

func newGitProvider() GitProvider {
	return &gitProvider{}
}

// gitProvider implements GitProvider.
type gitProvider struct{}

func (g *gitProvider) AddBranch(branch, base string) error {
	_, err := common.ExecGit("", "branch", branch, base)
	return err
}

func (g *gitProvider) ResolveBranch(branch string) (string, bool) {
	// Try to resolve as a local branch
	_, err := common.ExecGit("", "rev-parse", branch)
	if err == nil {
		return branch, true
	}

	// Try to resolve remote:branch or remote/branch
	split := strings.SplitN(branch, ":", 2)
	if len(split) < 2 {
		split = strings.SplitN(branch, "/", 2)
		if len(split) < 2 {
			// Can't parse this branch ref
			return "", false
		}
	}

	remote := split[0]
	remoteBranch := split[1]

	// Add remote if necessary
	_, err = g.GetRemote(remote)
	if err != nil {
		err = g.AddRemote(remote, "")
		if err != nil {
			fmt.Printf("WARNING: could not add remote %q: %v\n", remote, err)
			return "", false
		}
	}

	// Fetch branch if necessary
	err = g.Fetch(remote, remoteBranch)
	if err != nil {
		fmt.Printf("WARNING: could not fetch remote branch %q: %v\n", branch, err)
		return "", false
	}

	return fmt.Sprintf("%s/%s", remote, remoteBranch), true
}

func (g *gitProvider) AddWorktree(branch string) error {
	path, err := common.WorktreePath(branch)
	if err != nil {
		return err
	}
	return common.Git("worktree", "add", path, branch)

}

func (g *gitProvider) AddRemote(remoteName, url string) error {
	// TODO: fix up the dependency here
	return addRemote(remoteName, url)
}

func (g *gitProvider) GetRemote(remoteName string) (url string, err error) {
	return common.ExecGit("", "remote", "get-url", remoteName)
}

func (g *gitProvider) Fetch(remote, branch string) error {
	_, err := common.ExecGit("", "fetch", remote, branch)
	return err
}
