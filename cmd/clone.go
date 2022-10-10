package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"net/url"
	"path/filepath"
	"strings"
)

// Clone clones the provided repo, using the workflow described in
// https://morgan.cugerone.com/blog/how-to-use-git-worktree-and-in-a-clean-way/
// TODO: ensure fork is created
func Clone(args []string) error {
	var user, repo string

	var arg1 string
	err := common.ReqArg(args, 0, "URL of Git repo to clone:", &arg1)
	if err != nil {
		return err
	}
	user, repo, err = parseGitRepoURL(arg1)
	if err != nil {
		return err
	}

	if repo == "" {
		// arg1 was just "user"
		err := common.ReqArg(args, 1, fmt.Sprint("Which of %s's repos do you want?", user), &repo)
		if err != nil {
			return err
		}
	}

	// Reconstruct repository URL
	repoURL := fmt.Sprintf("https://github.com/%s/%s", user, repo)

	// Use JIT_DIR to find clone path
	jitDir, err := common.JitDir()
	if err != nil {
		return err
	}
	cloneDir := filepath.Join(jitDir, user, repo)

	// Make a bare repo
	// git clone --bare <repository> <directory>
	err = common.Git("clone", []string{"--bare", repoURL, filepath.Join(cloneDir, ".git")})
	if err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	// Add remote "fork"

	// Success - print message to user
	fmt.Printf(`
Successfully cloned repo %s/%s into %v
Copy remote branches using
    jit get <remote-branch>
or create new branches using
    jit new <branch> <base>
`[1:], user, repo, cloneDir)
	return nil
}

// Parses the given argument to a GitHub user & repo.
// Valid inputs are
//
//	"user"                         -> "user", "",     nil
//	"user/repo"                    -> "user", "repo", nil
//	"https://github.com/user/repo" -> "user", "repo", nil
func parseGitRepoURL(raw string) (string, string, error) {
	u, err := url.Parse(raw) // only matches full URL with scheme
	if err != nil {
		return "", "", err
	}

	switch u.Host {
	case "":
	// raw is not a URL
	case "github.com":
		raw = u.Path[1:] // strip leading '/'
	default:
		return "", "", fmt.Errorf("host %s not supported", u.Host)
	}

	// raw is now "user" or "user/repo"
	parts := strings.Split(raw, "/")
	switch len(parts) {
	case 1:
		// "user"
		return parts[0], "", nil
	case 2:
		// "user/repo"
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("invalid Git repo URL %s", raw)
	}
}
