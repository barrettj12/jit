package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/barrettj12/jit/common"
)

func AddRemote(args []string) error {
	remoteName, err := common.ReqArg(args, 0, "Which branch would you like to remove?")
	if err != nil {
		return err
	}

	var remoteURL string
	if len(args) >= 2 {
		remoteURL = args[1]
	}
	return addRemote(remoteName, remoteURL)
}

func addRemote(name, url string) error {
	if url == "" {
		repoPath, err := common.RepoBasePath()
		if err != nil {
			return err
		}
		repoName := filepath.Base(repoPath)
		url = githubURL(name, repoName)
	}

	_, err := common.ExecGit("", "remote", "add", name, url)
	if err != nil {
		return fmt.Errorf("error adding remote: %w", err)
	}
	fmt.Printf("added remote %s (%s)\n", name, url)
	return nil
}
