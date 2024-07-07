package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	"github.com/barrettj12/jit/common"
)

var newCmd = &cobra.Command{
	Use:   "new <branch> [based-on]",
	Short: "Create a new branch",
	RunE:  NewV2,
}

// TODO: this works but why are we getting "detached HEAD" ?
// TODO: need to be careful with branch vs remote/branch
func NewV2(cmd *cobra.Command, args []string) error {
	git := newGitProvider()

	branch, err := common.ReqArg(args, 0, "Which branch do you want to create/get?")
	if err != nil {
		return err
	}

	if len(args) < 2 {
		// Resolve an existing branch name
		branchSpec, ok := git.ResolveBranch(branch)
		if ok {
			return git.AddWorktree(branchSpec)
		}
		fmt.Printf("could not resolve branch %q\n", branch)
	}

	// If we're here, either:
	// - The user specified a base branch for the new worktree
	// - We couldn't resolve the provided arg to an existing branch
	// So we'll be creating a new branch based on an existing branch.

	base, err := common.ReqArg(args, 1, "Which branch should this be based on?")
	if err != nil {
		return err
	}

	fmt.Printf("Pulling branch %q...\n", base)
	err = pull(base)
	if err != nil {
		return err
	}

	err = git.AddBranch(branch, base)
	if err != nil {
		return err
	}

	err = git.AddWorktree(branch)
	if err != nil {
		return err
	}
	return Edit(nil, []string{branch})
}

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
// TODO: move to separate package
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
			fmt.Printf("WARNING: could not add remote %q: %v", remote, err)
			return "", false
		}
	}

	// Fetch branch if necessary
	err = g.Fetch(remote, remoteBranch)
	if err != nil {
		fmt.Printf("WARNING: could not fetch remote branch %q: %v", branch, err)
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
