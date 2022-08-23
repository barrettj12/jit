package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/go-git/go-git/v5"
	"net/url"
	"path/filepath"
)

// Clone clones the provided repo, using the workflow described in
// https://morgan.cugerone.com/blog/how-to-use-git-worktree-and-in-a-clean-way/
// TODO: ensure fork is created
func Clone(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide URL to clone")
	}
	urlStr := args[0]

	// Use JIT_DIR to find clone path
	jitDir, err := common.JitDir()
	if err != nil {
		return err
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing url: %w", err)
	}
	cloneDir := filepath.Join(jitDir, u.Path)

	// Make a bare repo
	_, err = git.PlainClone(filepath.Join(cloneDir, ".git"), true, &git.CloneOptions{
		URL:        urlStr,
		RemoteName: "source",
	})
	if err != nil {
		return fmt.Errorf("error cloning repo: %w", err)
	}

	// Success - print message to user
	fmt.Printf(`
Successfully cloned repo %v into %v
Copy remote branches using
    jit get <remote-branch>
or create new branches using
    jit new <branch> <base>
`[1:], urlStr, cloneDir)
	return nil
}
