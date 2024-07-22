package common

import (
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
)

// GitRemoteFromURL finds the name of the Git remote which has the given URL.
func GitRemoteFromURL(repo url.GitHubRepo) (types.RemoteName, error) {
	remotes, err := git.ListRemotes()
	if err != nil {
		return "", fmt.Errorf("getting remotes: %v", err)
	}

	for _, remote := range remotes {
		if remote.PushURL.Owner() == repo.Owner() {
			return remote.Name, nil
		}
	}
	return "", fmt.Errorf("no remote found matching %q", repo)
}
