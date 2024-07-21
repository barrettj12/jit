package url

import (
	"net/url"
	"path"
	"strings"
)

// RemoteRepo represents a remote Git repository on GitHub, GitLab, etc
type RemoteRepo interface {
	URL() string
}

// Raw is a raw URL.
type Raw string

func (u Raw) URL() string {
	return string(u)
}

// Nil represents an unspecified URL.
var Nil RemoteRepo = Raw("")

func IsNil(r RemoteRepo) bool {
	return r == nil || r.URL() == ""
}

// GitHubRepo represents the URL to a repo on GitHub, in the form
//
//	https://github.com/user/repo
type GitHubRepo string

// GitHubURL converts the given path components into a GitHub URL.
//
//	"user"                         -> "https://github.com/user"
//	"user/repo"                    -> "https://github.com/user/repo"
//	"user", "repo"                 -> "https://github.com/user/repo"
//	"https://github.com/user/repo" -> "https://github.com/user/repo"
func GitHubURL(c ...string) GitHubRepo {
	// TODO: validation?
	if len(c) > 0 {
		u, err := url.Parse(c[0])
		if err == nil {
			c[0] = u.Path
		}
	}
	pathComponents := append([]string{"github.com"}, c...)
	return GitHubRepo("https://" + path.Join(pathComponents...))
}

func (r GitHubRepo) URL() string {
	return string(r)
}

func (r GitHubRepo) Owner() string {
	parsed, _ := url.Parse(string(r))
	split := strings.Split(parsed.Path, "/")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func (r GitHubRepo) RepoName() string {
	parsed, _ := url.Parse(string(r))
	split := strings.Split(parsed.Path, "/")
	if len(split) > 2 {
		return split[2]
	}
	return ""
}
