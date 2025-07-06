package url

import (
	"net/url"
	"path"
	"strings"
)

// RemoteRepo represents a remote Git repository on GitHub, GitLab, etc
type RemoteRepo interface {
	URL() string
	Owner() string
	RepoName() string
	HostedBy() RepoSource
}

// Raw is a raw URL.
type Raw string

func (u Raw) URL() string {
	return string(u)
}

func (u Raw) Owner() string {
	// Assume the first url component is the owner
	parsed, _ := url.Parse(string(u))
	split := strings.Split(parsed.Path, "/")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func (u Raw) RepoName() string {
	// Assume the second url component is the repo name
	parsed, _ := url.Parse(string(u))
	split := strings.Split(parsed.Path, "/")
	if len(split) > 2 {
		return split[2]
	}
	return ""
}

func (u Raw) HostedBy() RepoSource {
	parsed, err := url.Parse(string(u))
	if err == nil {
		switch parsed.Host {
		case "github.com":
			return GitHub
		case "gitlab.com":
			return GitLab
		default:
			return UnknownSite
		}
	}
	return UnknownSite
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

func (r GitHubRepo) HostedBy() RepoSource {
	return GitHub
}

// URL converts the given path components into a repo URL. If the website is
// not specified, GitHub will be assumed.
//
//	"user"                         -> "https://github.com/user"
//	"user/repo"                    -> "https://github.com/user/repo"
//	"user", "repo"                 -> "https://github.com/user/repo"
//	"https://server.com/user/repo" -> "https://server.com/user/repo"
func URL(c ...string) RemoteRepo {
	if len(c) == 0 {
		return Nil
	}
	parsed, err := url.Parse(c[0])
	if err != nil {
		return GitHubURL(c...)
	}
	if parsed.Host == "" {
		// Assume GitHub
		return GitHubURL(c...)
	}
	return Raw(c[0])
}

// RepoSource represents a possible hosting site for a Git repo, e.g. GitHub,
// GitLab, private website, ...
type RepoSource string

const (
	GitHub      RepoSource = "github"
	GitLab      RepoSource = "gitlab"
	UnknownSite RepoSource = "unknown"
)
