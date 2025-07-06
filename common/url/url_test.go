package url

import (
	"testing"

	"github.com/barrettj12/jit/common/testutil"
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

func TestParseURL(t *testing.T) {
	tests := []struct {
		input    []string
		url      string
		owner    string
		repoName string
		hostedBy RepoSource
	}{{
		input:    []string{"user"},
		url:      "https://github.com/user",
		owner:    "user",
		repoName: "",
		hostedBy: GitHub,
	}, {
		input:    []string{"user", "repo"},
		url:      "https://github.com/user/repo",
		owner:    "user",
		repoName: "repo",
		hostedBy: GitHub,
	}, {
		input:    []string{"user/repo"},
		url:      "https://github.com/user/repo",
		owner:    "user",
		repoName: "repo",
		hostedBy: GitHub,
	}, {
		input:    []string{"https://github.com/user/repo"},
		url:      "https://github.com/user/repo",
		owner:    "user",
		repoName: "repo",
		hostedBy: GitHub,
	}, {
		input:    []string{"https://gitlab.com/user/repo"},
		url:      "https://gitlab.com/user/repo",
		owner:    "user",
		repoName: "repo",
		hostedBy: GitLab,
	}}

	for _, test := range tests {
		url := URL(test.input...)
		testutil.AssertEqual(t, url.URL(), test.url)
		testutil.AssertEqual(t, url.Owner(), test.owner)
		testutil.AssertEqual(t, url.RepoName(), test.repoName)
		testutil.AssertEqual(t, url.HostedBy(), test.hostedBy)
	}
}
