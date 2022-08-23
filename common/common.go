package common

import (
	"bufio"
	"bytes"
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

// ReqArg will first see if args[i] has been defined
// If so, it will set *ptr to this value
// If not, it will prompt the user to enter a value
func ReqArg(args []string, i int, prompt string, ptr *string) error {
	if len(args) > i {
		*ptr = args[i]
		return nil
	}

	// Argument was not defined, so prompt the user for input
	sc := bufio.NewScanner(os.Stdin)
	fmt.Printf("%v ", prompt)
	sc.Scan()
	if err := sc.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	*ptr = sc.Text()
	return nil
}

// Returns the root folder of the current repo
func RepoBasePath() (string, error) {
	b := bytes.Buffer{}
	cmd := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
