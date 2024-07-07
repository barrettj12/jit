package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
	"path/filepath"

	"github.com/barrettj12/jit/common"
)

var addRemoteCmd = &cobra.Command{
	Use:   "add-remote <user>",
	Short: "Add a remote to a repo",
	RunE:  AddRemote,
}

func AddRemote(cmd *cobra.Command, args []string) error {
	remoteName, err := common.ReqArg(args, 0, "Which remote would you like to add?")
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

	err := git.AddRemote(name, url)
	if err != nil {
		return fmt.Errorf("error adding remote: %w", err)
	}
	fmt.Printf("added remote %s (%s)\n", name, url)
	return nil
}
