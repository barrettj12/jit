package common

import (
	"bufio"
	"fmt"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
func JitDir() (path.JitDir, error) {
	jitDir, ok := os.LookupEnv("JIT_DIR")
	if !ok {
		return "", fmt.Errorf("env var JIT_DIR not set")
	}
	return path.JitDir(jitDir), nil
}

func DefaultRepoBasePath(user, repo string) (path.Repo, error) {
	jitDir, err := JitDir()
	if err != nil {
		return "", fmt.Errorf("getting jit dir: %w", err)
	}

	return path.RepoPath(jitDir, user, repo), nil
}

// ReqArg will first see if args[i] has been defined.
// If so, it will return this value.
// If not, it will prompt the user to enter a value.
func ReqArg(args []string, i int, prompt string) (string, error) {
	if len(args) > i {
		return args[i], nil
	}
	// Argument was not defined, so prompt the user for input
	return Prompt(prompt)
}

// Prompt the user to enter a value.
func Prompt(prompt string) (string, error) {
	nonInteractive := os.Getenv("JIT_NONINTERACTIVE")
	parsed, err := strconv.ParseBool(nonInteractive)
	if err == nil && parsed {
		panic("internal error: common.Prompt called with JIT_NONINTERACTIVE=1")
	}

	sc := bufio.NewScanner(os.Stdin)
	fmt.Printf("%v ", prompt)
	sc.Scan()
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}
	return sc.Text(), nil
}

// Returns the root folder of the current repo
func RepoBasePath() (path.Repo, error) {
	stdout, err := ExecGit(path.CurrentDir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", err
	}

	basepath := filepath.Dir(stdout)
	return path.Repo(basepath), nil
}

// Returns the absolute file path to the given branch/worktree
// TODO: replace this with LookupWorktreeForBranch
func WorktreePath(branch string) (string, error) {
	// Get path to new worktree
	gitDir, err := RepoBasePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir.Path(), branch), nil
}

var ErrUpstreamNotFound = fmt.Errorf("upstream not found")

// Returns push location (remote, branch) for the given branch
// TODO: replace this with git.PushTarget
func PushLoc(localBranch string) (remote, remoteBranch string, err error) {
	stdout, err := ExecGit(path.CurrentDir, "for-each-ref", "--format='%(push:short)'",
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

func DefaultRemote() types.RemoteName {
	return types.RemoteName(GitHubUser())
}

// Fetches the given branches.
// If remote == "", it will fetch all branches.
// If branch == "", it will fetch all branches for the given remote.
// TODO: replace with git.Fetch
func Fetch(remote, branch string) error {
	args := []string{"fetch"}
	if remote != "" {
		args = append(args, remote)
		if branch != "" {
			args = append(args, branch)
		}
	}

	baseDir, err := RepoBasePath()
	if err != nil {
		return err
	}
	_, err = ExecGit(baseDir, args...)
	return err
}
