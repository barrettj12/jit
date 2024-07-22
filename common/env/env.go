package env

import (
	"fmt"
	"github.com/barrettj12/jit/common/path"
	"os"
	"strconv"
)

// JitDir returns the value of the environment variable JIT_DIR, which should
// point to the base directory where Jit will clone repos to.
func JitDir() (path.JitDir, error) {
	jitDir := os.Getenv("JIT_DIR")
	if jitDir == "" {
		return "", fmt.Errorf("env var JIT_DIR not set")
	}
	return path.JitDir(jitDir), nil
}

// GitHubUser returns the value of the environment variable GH_USER, which
// should contain the user's GitHub username.
func GitHubUser() (string, error) {
	ghUser := os.Getenv("GH_USER")
	if ghUser == "" {
		return "", fmt.Errorf("env var GH_USER not set")
	}
	return ghUser, nil
}

// Debug returns the value of the environment variable JIT_DEBUG, which
// indicates that Jit should print any Git commands that it runs.
func Debug() bool {
	debug := os.Getenv("JIT_DEBUG")
	parsed, err := strconv.ParseBool(debug)
	if err != nil {
		return false
	}
	return parsed
}

// NonInteractive returns the value of the environment variable
// JIT_NONINTERACTIVE, which indicates that Jit should not prompt for user
// input in the current environment.
func NonInteractive() bool {
	nonInteractive := os.Getenv("JIT_NONINTERACTIVE")
	parsed, err := strconv.ParseBool(nonInteractive)
	if err != nil {
		return false
	}
	return parsed
}
