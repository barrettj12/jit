package common

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Execute(script string, args ...string) error {
	cmd := exec.Command(script, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Git(cmd string, args ...string) error {
	gitArgs := append([]string{cmd}, args...)
	return Execute("git", gitArgs...)
}

// Returns the value of env variable JIT_DIR
func JitDir() (string, error) {
	path, ok := os.LookupEnv("JIT_DIR")
	if !ok {
		return "", fmt.Errorf("env var JIT_DIR not set")
	}
	return path, nil
}

// ReqArg will first see if args[i] has been defined.
// If so, it will return this value.
// If not, it will prompt the user to enter a value.
func ReqArg(args []string, i int, prompt string) (string, error) {
	if len(args) > i {
		return args[i], nil
	}

	// Argument was not defined, so prompt the user for input
	sc := bufio.NewScanner(os.Stdin)
	fmt.Printf("%v ", prompt)
	sc.Scan()
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}
	return sc.Text(), nil
}

// Returns the root folder of the current repo
func RepoBasePath() (string, error) {
	stdout, err := ExecGit("", "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", err
	}

	path := filepath.Dir(stdout)
	return path, nil
}

// Returns the absolute file path to the given branch/worktree
func WorktreePath(branch string) (string, error) {
	// Get path to new worktree
	gitDir, err := RepoBasePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, branch), nil
}

var ErrUpstreamNotFound = fmt.Errorf("upstream not found")

// Returns push location (remote, branch) for the given branch
func PushLoc(localBranch string) (remote, remoteBranch string, err error) {
	stdout, err := ExecGit("", "for-each-ref", "--format='%(push:short)'",
		fmt.Sprintf("refs/heads/%s", localBranch))
	if err != nil {
		return "", "", err
	}

	pushloc := strings.Trim(stdout, "'\n")
	if pushloc == "" {
		return "", "", ErrUpstreamNotFound
	}
	split := strings.Split(pushloc, "/")
	return split[0], split[1], nil
}

func GitHubUser() string {
	return os.Getenv("GH_USER")
}
