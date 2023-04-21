package common

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// Nice interface for running command line methods
type ExecArgs struct {
	Cmd  string
	Args []string
	Dir  string

	Stdout, Stderr io.Writer

	Background bool
}

type ExecResult struct {
	RunError error
	// output
	Stdout, Stderr, Combined string
}

// TODO: have a way to globally change the default dir
func Exec(args ExecArgs) ExecResult {
	cmd := exec.Command(args.Cmd, args.Args...)
	cmd.Dir = args.Dir

	// Set up stdin/out/err
	cmd.Stdin = os.Stdin

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combined := &bytes.Buffer{}

	outWriters := []io.Writer{stdout, combined}
	if args.Stdout != nil {
		outWriters = append(outWriters, args.Stdout)
	}
	cmd.Stdout = io.MultiWriter(outWriters...)

	errWriters := []io.Writer{stderr, combined}
	if args.Stderr != nil {
		errWriters = append(errWriters, args.Stderr)
	}
	cmd.Stderr = io.MultiWriter(errWriters...)

	var runErr error
	if args.Background {
		runErr = cmd.Start()
	} else {
		runErr = cmd.Run() // this error contains the exit code
	}

	return ExecResult{
		RunError: runErr,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Combined: combined.String(),
	}
}
