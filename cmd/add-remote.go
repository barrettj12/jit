package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
	"github.com/spf13/cobra"
)

var addRemoteCmd = &cobra.Command{
	Use:   "add-remote <user>",
	Short: "Add a remote to a repo",
	RunE:  AddRemote,
}

func AddRemote(cmd *cobra.Command, args []string) error {
	name, err := common.ReqArg(args, 0, "Which remote would you like to add?")
	if err != nil {
		return err
	}

	var repo string
	if len(args) >= 2 {
		repo = args[1]
	}
	return addRemote(types.RemoteName(name), url.Raw(repo))
}

// TODO: move to common
func addRemote(name types.RemoteName, repo url.RemoteRepo) error {
	if url.IsNil(repo) {
		repoPath, err := common.RepoBasePath()
		if err != nil {
			return fmt.Errorf("getting repo base path: %w", err)
		}
		repo = url.GitHubURL(string(name), repoPath.RepoName())
	}

	err := git.AddRemote(name, repo)
	if err != nil {
		return fmt.Errorf("error adding remote: %w", err)
	}
	fmt.Printf("added remote %s (%s)\n", name, repo)
	return nil
}
