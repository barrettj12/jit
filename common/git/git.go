package git

import (
	"bytes"
	"fmt"
	"github.com/barrettj12/jit/common/env"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CloneArgs struct {
	Repo       url.RemoteRepo   // repo to clone
	CloneDir   path.Dir         // directory to clone the repo into
	Bare       bool             // whether to make a bare clone
	OriginName types.RemoteName // the name to give to the origin remote
}

func Clone(opts CloneArgs) error {
	args := []string{"clone"}
	if opts.Bare {
		args = append(args, "--bare")
	}
	if opts.OriginName != "" {
		args = append(args, "--origin", string(opts.OriginName))
	}
	args = append(args, opts.Repo.URL())
	if out := path.Path(opts.CloneDir); out != "" {
		args = append(args, out)
	}

	_, err := internalExec(internalExecArgs{args: args})
	return err
}

func SetConfig(dir path.Dir, key, value string) error {
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

type RebaseArgs struct {
	Base        types.LocalBranch // base branch/ref to rebase against
	Interactive bool
	Env         []string
}

func Rebase(opts RebaseArgs) error {
	args := []string{"rebase"}
	if opts.Interactive {
		args = append(args, "-i")
	}
	args = append(args, string(opts.Base))

	_, err := internalExec(internalExecArgs{
		args: args,
		env:  opts.Env,
	})
	return err
}

func MergeBase(branch1, branch2 types.LocalBranch) (types.LocalBranch, error) {
	out, err := internalExec(internalExecArgs{
		args: []string{"merge-base", string(branch1), string(branch2)},
	})
	if err != nil {
		return "", err
	}
	return types.LocalBranch(strings.TrimSpace(out)), nil
}

type MergeArgs struct {
	Branch types.LocalBranch // branch to merge into current
	Squash bool
}

func Merge(opts MergeArgs) error {
	args := []string{"merge"}
	if opts.Squash {
		args = append(args, "--squash")
	}
	args = append(args, string(opts.Branch))

	_, err := internalExec(internalExecArgs{
		args:         args,
		attachStdout: true,
	})
	return err
}

type CommitArgs struct {
	Message string // commit message to use
}

func Commit(opts CommitArgs) error {
	args := []string{"commit"}
	if opts.Message != "" {
		args = append(args, "-m", opts.Message)
	}

	_, err := internalExec(internalExecArgs{
		args: args,
	})
	return err
}

type internalExecArgs struct {
	args         []string // args to feed to git
	dir          path.Dir // directory to run the command in
	attachStdout bool     // if true, attach cmd stdout to os.Stdout
	attachStderr bool     // if true, attach cmd stderr to os.Stderr
	env          []string // environment variable key=value pairs
}

// Runs git with the given args, returning stdout and/or any error.
var internalExec = func(opts internalExecArgs) (string, error) {
	cmd := exec.Command("git", opts.args...)
	cmd.Dir = path.Path(opts.dir)
	cmd.Env = append(cmd.Environ(), opts.env...)

	// Handle stdout/stderr
	stdoutBuffer := &bytes.Buffer{}
	if opts.attachStdout {
		cmd.Stdout = io.MultiWriter(stdoutBuffer, os.Stdout)
	} else {
		cmd.Stdout = stdoutBuffer
	}

	stderrBuffer := &bytes.Buffer{}
	if opts.attachStderr {
		cmd.Stderr = io.MultiWriter(stderrBuffer, os.Stderr)
	} else {
		cmd.Stderr = stderrBuffer
	}

	if env.Debug() {
		fmt.Println(cmd.String())
	}

	var runErr error
	runErr = cmd.Run() // this error contains the exit code

	// handle errors
	if runErr != nil {
		// Read stderr for error info
		errInfo := stderrBuffer.String()
		return "", fmt.Errorf("%s\n%w", errInfo, runErr)
	}

	return stdoutBuffer.String(), nil
}
