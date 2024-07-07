package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/barrettj12/jit/common"
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
	split := strings.Split(basePath, "/")
	user := split[len(split)-2]
	repo := split[len(split)-1]
	return fork(user, repo, "")
}

func fork(user, repo, dir string) error {
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
	return addRemote(ghUser, "")
}
