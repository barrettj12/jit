package cmd

import (
	"fmt"
	"strings"

	"github.com/barrettj12/jit/common"
)

func Fork(args []string) error {
	basePath, err := common.RepoBasePath()
	if err != nil {
		return err
	}
	split := strings.Split(basePath, "/")
	user := split[len(split)-2]
	repo := split[len(split)-1]
	return fork(user, repo)
}

func fork(user, repo string) error {
	// Create fork
	err := common.Execute("gh", "repo", "fork",
		fmt.Sprintf("%s/%s", user, repo), "--clone=false")
	if err != nil {
		return fmt.Errorf("error creating fork: %w", err)
	}

	// Add as remote
	ghUser := common.GitHubUser()
	return addRemote(ghUser, "")
}
