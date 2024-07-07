package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func CurrentBranch() (string, error) {
	// git rev-parse --abbrev-ref HEAD
	out, err := internalExec("", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

type CloneArgs struct {
	RepoURL    string // URL of repo to clone
	CloneDir   string // directory to clone the repo into
	Bare       bool   // whether to make a bare clone
	OriginName string // the name to give to the origin remote
}

func Clone(opts CloneArgs) error {
	args := []string{"clone"}
	if opts.Bare {
		args = append(args, "--bare")
	}
	if opts.OriginName != "" {
		args = append(args, "--origin", opts.OriginName)
	}
	args = append(args, opts.RepoURL)
	if opts.CloneDir != "" {
		args = append(args, opts.CloneDir)
	}

	_, err := internalExec("", args...)
	return err
}

type AddWorktreeArgs struct {
	Dir          string // directory to run the command in
	WorktreePath string // path for the new worktree
	Branch       string // branch to check out in the new worktree (optional)
}

func AddWorktree(opts AddWorktreeArgs) error {
	args := []string{"worktree", "add", opts.WorktreePath}
	if opts.Branch != "" {
		args = append(args, opts.Branch)
	}
	_, err := internalExec(opts.Dir, args...)
	return err
}

// Create a new branch `name` based on `base`.
func CreateBranch(name, base string) error {
	_, err := internalExec("", "branch", name, base)
	return err
}

func RemoteExists(dir, remote string) (bool, error) {
	_, err := internalExec(dir, "remote", "get-url", remote)
	if err == nil {
		return true, nil
	}
	if IsNoSuchRemoteErr(err) {
		return false, nil
	}
	return false, err
}

func Fetch(dir, remote, branch string) error {
	_, err := internalExec(dir, "fetch", remote, branch)
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

func SetConfig(dir, key, value string) error {
	_, err := internalExec(dir, "config", key, value)
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

// Runs git with the given args, returning stdout and/or any error.
func internalExec(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr

	var runErr error
	runErr = cmd.Run() // this error contains the exit code

	// handle errors
	if runErr != nil {
		// Read stderr for error info
		errInfo := stderr.String()
		return "", fmt.Errorf("%s\n%w", errInfo, runErr)
	}

	return stdout.String(), nil
}
