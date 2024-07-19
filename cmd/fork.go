package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/path"
	"github.com/barrettj12/jit/common/types"
	"github.com/barrettj12/jit/common/url"
	"github.com/spf13/cobra"
	"os"
)

var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Fork the current Git repo",
	RunE:  Fork,
}

func Fork(cmd *cobra.Command, args []string) error {
	basePath, err := common.RepoBasePath()
	if err != nil {
		return err
	}
	return fork(path.CurrentDir, basePath.Owner(), basePath.RepoName())
}

// TODO: move to common
func fork(dir path.Dir, user, repo string) error {
	// Create fork
	res := common.Exec(common.ExecArgs{
		Cmd:    "gh",
		Args:   []string{"repo", "fork", fmt.Sprintf("%s/%s", user, repo), "--clone=false"},
		Dir:    dir,
		Stdout: os.Stdout,
	})
	if res.RunError != nil {
		return fmt.Errorf(`error creating fork: %w
stderr: %s`, res.RunError, res.Stderr)
	}

	// Add as remote
	ghUser := common.GitHubUser()
	remoteName := types.RemoteName(ghUser)
	return addRemote(remoteName, url.Nil)
}
