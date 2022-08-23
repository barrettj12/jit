package common

import (
	"fmt"
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

// Returns the value of env variable JIT_DIR
func JitDir() (string, error) {
	path, ok := os.LookupEnv("JIT_DIR")
	if !ok {
		return "", fmt.Errorf("env var JIT_DIR not set")
	}
	return path, nil
}
