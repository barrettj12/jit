package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// SetupTestRepo creates a simple Git repository for testing in a temporary
// folder. If a dir is passed in, the test repo will be created inside that
// dir, otherwise it will be placed in a temporary folder. In any case, the
// full path to the created test repo will be returned.
func SetupTestRepo(t *testing.T, dir string) (path string) {
	repoPath, err := os.MkdirTemp(dir, "testrepo")
	CheckErr(t, err)
	t.Cleanup(func() {
		err := os.RemoveAll(repoPath)
		if err != nil {
			t.Logf("error cleaning up dir %q: %v", repoPath, err)
		}
	})

	// Initialise git repo
	RunCommand(t, repoPath, "git", "init")

	// Add file
	fileName := "foo.txt"
	file, err := os.Create(filepath.Join(repoPath, fileName))
	CheckErr(t, err)
	_, err = file.WriteString("hello world")
	CheckErr(t, err)
	err = file.Close()
	CheckErr(t, err)

	// git add and commit
	RunCommand(t, repoPath, "git", "add", fileName)
	RunCommand(t, repoPath, "git", "commit", "-m", `"Initial commit"`)
	return repoPath
}
