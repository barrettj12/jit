package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

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

func SetConfig(dir, key, value string) error {
	_, err := internalExec(dir, "config", key, value)
	return err
}

func Apply(path string) error {
	_, err := internalExec("", "apply", path)
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
