package common

import (
	"fmt"
	"regexp"
)

var reErrGitNotARepo = regexp.MustCompile("not a git repository")

func ExecGit(dir string, args ...string) (string, error) {
	res := Exec(ExecArgs{
		Cmd:  "git",
		Args: args,
		Dir:  dir,
	})

	// handle errors
	if res.RunError != nil {
		// Read stderr for error info
		errInfo := res.Stderr

		if reErrGitNotARepo.MatchString(errInfo) {
			return "", ErrGitNotARepo
		} else {
			return "", fmt.Errorf("%s\n%s", errInfo, res.RunError)
		}
	}

	return res.Stdout, nil
}

// Git errors
type GitError string

func (e GitError) Error() string {
	return string(e)
}

const (
	ErrGitNotARepo = GitError("current dir is not inside a git repo")
)
