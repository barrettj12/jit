package git

import (
	"fmt"
	"strings"
)

func CurrentBranch() (string, error) {
	out, err := internalExec("", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Create a new branch `name` based on `base`.
func CreateBranch(name, base string) error {
	_, err := internalExec("", "branch", name, base)
	return err
}

// Retrieves the push target for the specified branch. You can use branch = ""
// for the current branch.
// A return value of "" means no upstream is set.
func PushTarget(branch string) (string, error) {
	out, err := internalExec("", "rev-parse", "--abbrev-ref", fmt.Sprintf("%s@{push}", branch))
	if err == nil {
		return strings.TrimSpace(out), nil
	}
	if IsNoUpstreamConfiguredErr(err) {
		return "", nil
	}
	return "", err
}

type PushArgs struct {
	Remote      string // remote repository to push to
	Branch      string // branch to push
	SetUpstream bool   // should the upstream be set on a successful push
}

func Push(opts PushArgs) error {
	args := []string{"push"}
	if opts.SetUpstream {
		args = append(args, "-u")
	}
	if opts.Remote != "" {
		args = append(args, opts.Remote)
	}
	if opts.Branch != "" {
		args = append(args, opts.Branch)
	}

	_, err := internalExec("", args...)
	return err
}

func Pull(dir string) error {
	_, err := internalExec(dir, "pull")
	return err
}

func SetUpstream(dir, localBranch, remote, remoteBranch string) error {
	_, err := internalExec(dir, "branch", "-u",
		fmt.Sprintf("%s/%s", remote, remoteBranch), localBranch)
	return err
}
