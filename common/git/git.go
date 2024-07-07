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

	fmt.Println(args)
	_, err := internalExec("", args...)
	return err
}

func AddWorktree(dir, name string) error {
	_, err := internalExec(dir, "worktree", "add", name)
	return err
}

// Set the upstream of localBranch to remote:remoteBranch.
func SetUpstream(localBranch, remote, remoteBranch string) error {
	// TODO: implement
	return nil
}

// Create a new branch `name` based on `base`.
func CreateBranch(name, base string) error {
	// TODO: implement
	return nil
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
