package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

	_, err := internalExec(internalExecArgs{args: args})
	return err
}

func SetConfig(dir, key, value string) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"config", key, value},
		dir:  dir,
	})
	return err
}

func Apply(path string) error {
	_, err := internalExec(internalExecArgs{
		args: []string{"apply", path},
	})
	return err
}

type internalExecArgs struct {
	args         []string // args to feed to git
	dir          string   // directory to run the command in
	attachStderr bool     // if true, attach cmd stderr to os.Stderr
}

// Runs git with the given args, returning stdout and/or any error.
func internalExec(opts internalExecArgs) (string, error) {
	cmd := exec.Command("git", opts.args...)
	cmd.Dir = opts.dir

	// Handle stdout/stderr
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout

	stderrBuffer := &bytes.Buffer{}
	if opts.attachStderr {
		cmd.Stderr = io.MultiWriter(stderrBuffer, os.Stderr)
	} else {
		cmd.Stderr = stderrBuffer
	}

	var runErr error
	runErr = cmd.Run() // this error contains the exit code

	// handle errors
	if runErr != nil {
		// Read stderr for error info
		errInfo := stderrBuffer.String()
		return "", fmt.Errorf("%s\n%w", errInfo, runErr)
	}

	return stdout.String(), nil
}
