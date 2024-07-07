package main

import (
	"fmt"
	"github.com/barrettj12/jit/common/testutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGitServer(t *testing.T) {
	// Create a directory to store Git repos
	reposRoot, err := os.MkdirTemp("", "repos")
	testutil.CheckErr(t, err)
	t.Cleanup(func() {
		err := os.RemoveAll(reposRoot)
		if err != nil {
			t.Logf("error cleaning up dir %q: %v", reposRoot, err)
		}
	})

	repoPath := testutil.SetupTestRepo(t, reposRoot)
	repoName := filepath.Base(repoPath)

	err = os.Setenv("GIT_PROJECT_ROOT", reposRoot)
	testutil.CheckErr(t, err)
	// run Git server in background
	go main()

	cloneDir, err := os.MkdirTemp("", "clone")
	t.Cleanup(func() {
		err := os.RemoveAll(cloneDir)
		if err != nil {
			t.Logf("error cleaning up dir %q: %v", cloneDir, err)
		}
	})

	testutil.RunCommand(t, "", "git", "clone",
		fmt.Sprintf("http://localhost:8080/%s", repoName),
		cloneDir,
	)
}
