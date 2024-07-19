package types

import (
	"fmt"
	url "github.com/barrettj12/jit/common/url"
	"strings"
)

// GitHubBranch represents a remote branch on GitHub.
type GitHubBranch struct {
	RepoURL url.GitHubRepo
	Branch  string
}

// GitHub formats remote branches as remote:branch
func ParseGitHubBranch(s string) GitHubBranch {
	split := strings.SplitN(s, ":", 2)
	return GitHubBranch{
		RepoURL: url.GitHubURL(split[0]),
		Branch:  split[1],
	}
}

func (b GitHubBranch) String() string {
	return fmt.Sprintf("%s:%s", b.RepoURL.Owner(), b.Branch)
}
