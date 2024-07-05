package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitServer(t *testing.T) {
	reposRoot, repoName := setupTestRepo(t)

	err := os.Setenv("GIT_PROJECT_ROOT", reposRoot)
	checkErr(t, err)
	// run Git server in background
	go main()

	cloneDir, err := os.MkdirTemp("", "clone")
	t.Cleanup(func() {
		err := os.RemoveAll(cloneDir)
		if err != nil {
			t.Logf("error cleaning up dir %q: %v", cloneDir, err)
		}
	})

	cloneCmd := exec.Command("git", "clone",
		fmt.Sprintf("http://localhost:8080/%s", repoName),
		cloneDir,
	)
	err = cloneCmd.Run()
	checkErr(t, err)
}

func setupTestRepo(t *testing.T) (string, string) {
	// Create a directory to store Git repos
	reposRoot, err := os.MkdirTemp("", "repos")
	checkErr(t, err)
	t.Cleanup(func() {
		err := os.RemoveAll(reposRoot)
		if err != nil {
			t.Logf("error cleaning up dir %q: %v", reposRoot, err)
		}
	})

	// Create a test repo inside the directory
	repoName := "testrepo"
	repoPath := filepath.Join(reposRoot, repoName)
	err = os.Mkdir(repoPath, 0777)
	checkErr(t, err)

	// Initialise git repo
	gitInitCmd := exec.Command("git", "init")
	gitInitCmd.Dir = repoPath
	err = gitInitCmd.Run()
	checkErr(t, err)

	// Add file
	fileName := "foo.txt"
	file, err := os.Create(filepath.Join(repoPath, fileName))
	checkErr(t, err)
	_, err = file.WriteString("hello world")
	checkErr(t, err)
	err = file.Close()
	checkErr(t, err)

	// git add and commit
	gitAddCmd := exec.Command("git", "add", fileName)
	gitAddCmd.Dir = repoPath
	err = gitAddCmd.Run()
	checkErr(t, err)

	gitCommitCmd := exec.Command("git", "commit", "-m", `"Initial commit"`)
	gitCommitCmd.Dir = repoPath
	err = gitCommitCmd.Run()
	checkErr(t, err)

	return reposRoot, repoName
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
