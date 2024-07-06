package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	runCommand(t, "", "git", "clone",
		fmt.Sprintf("http://localhost:8080/%s", repoName),
		cloneDir,
	)
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
	runCommand(t, repoPath, "git", "init")

	// Add file
	fileName := "foo.txt"
	file, err := os.Create(filepath.Join(repoPath, fileName))
	checkErr(t, err)
	_, err = file.WriteString("hello world")
	checkErr(t, err)
	err = file.Close()
	checkErr(t, err)

	// git add and commit
	runCommand(t, repoPath, "git", "add", fileName)
	runCommand(t, repoPath, "git", "commit", "-m", `"Initial commit"`)

	return reposRoot, repoName
}

func runCommand(t *testing.T, dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf(`error running command %q: %v
output: %s`, strings.Join(append([]string{name}, args...), " "), err, string(out))
	}
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
