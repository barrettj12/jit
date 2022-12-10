package common

import (
	"fmt"
	"regexp"
)

var GitNotARepoErr = regexp.MustCompile("not a git repository")

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

		if GitNotARepoErr.MatchString(errInfo) {
			return "", fmt.Errorf("current dir is not inside a git repo")
		} else {
			return "", fmt.Errorf("%s\n%s", errInfo, res.RunError)
		}
	}

	return res.Stdout, nil
}
