package common

import (
	"os"
	"os/exec"
)

func Execute(script string, args []string) error {
	cmd := &exec.Cmd{
		Path:   script,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
