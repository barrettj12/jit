package git_test

import (
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/testutil"
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
