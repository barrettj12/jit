package url

import (
	"github.com/barrettj12/jit/common/testutil"
	"testing"
)

func TestGitHubRepo(t *testing.T) {
	tests := []struct {
		url        GitHubRepo
		user, repo string
	}{{
		url:  "https://github.com/user/repo",
		user: "user",
		repo: "repo",
	}, {
		url:  "https://github.com/user",
		user: "user",
		repo: "",
	}}

	for _, test := range tests {
		testutil.AssertEqual(t, test.url.Owner(), test.user)
		testutil.AssertEqual(t, test.url.RepoName(), test.repo)
	}
}

func TestMakeGitHubURL(t *testing.T) {
	tests := []struct {
		input []string
		url   GitHubRepo
	}{{
		input: []string{"user"},
		url:   GitHubRepo("https://github.com/user"),
	}, {
		input: []string{"user", "repo"},
		url:   GitHubRepo("https://github.com/user/repo"),
	}, {
		input: []string{"user/repo"},
		url:   GitHubRepo("https://github.com/user/repo"),
	}, {
		input: []string{"https://github.com/user/repo"},
		url:   GitHubRepo("https://github.com/user/repo"),
	}}

	for _, test := range tests {
		url := GitHubURL(test.input...)
		testutil.AssertEqual(t, url, test.url)
	}
}

func TestIsNil(t *testing.T) {
	tests := []RemoteRepo{nil, Nil, Raw("")}
	for _, url := range tests {
		testutil.AssertEqual(t, IsNil(url), true)
	}
}
