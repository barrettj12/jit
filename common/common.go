package common

import (
	"os"
	"os/exec"
)

func Execute(script string, args []string) error {
	cmd := exec.Command(script, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Git(cmd string, args []string) error {
	gitArgs := append([]string{cmd}, args...)
	return Execute("git", gitArgs)
}
