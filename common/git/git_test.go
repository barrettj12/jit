package git_test

import (
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/testutil"
	"path/filepath"
	"testing"
)

func TestRemoteExists(t *testing.T) {
	repoPath := testutil.SetupTestRepo(t, "")
	remoteExists, err := git.RemoteExists(repoPath, "myremote")
	testutil.CheckErr(t, err)
	testutil.AssertEqual(t, remoteExists, false)

	testutil.RunCommand(t, repoPath, "git", "remote", "add", "myremote", "https://github.com/myremote/myrepo")
	remoteExists, err = git.RemoteExists(repoPath, "myremote")
	testutil.CheckErr(t, err)
	testutil.AssertEqual(t, remoteExists, true)
}

// Check that if a branch doesn't exist, git.AddWorktree returns a helpful
// error.
func TestAddWorktreeBranchDoesntExist(t *testing.T) {
	repoPath := testutil.SetupTestRepo(t, "")
	err := git.AddWorktree(git.AddWorktreeArgs{
		Branch:       "mybranch",
		WorktreePath: filepath.Join(repoPath, "mybranch"),
		Dir:          repoPath,
	})
	testutil.AssertNotEqual(t, err, nil)
	testutil.AssertEqual(t, err.Error(), `branch "mybranch" doesn't exist`)
}
