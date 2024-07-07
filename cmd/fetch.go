package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"github.com/barrettj12/jit/common/git"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch <branch>",
	Short: "Fetch a remote branch",
	RunE:  Fetch,
}

// Accepts the following formats:
//
//	fetch <remote> <branch>
//	fetch <remote>:<branch>
//	fetch <remote>/<branch>
func Fetch(cmd *cobra.Command, args []string) error {
	gitp := newGitProvider()

	remote, err := common.ReqArg(args, 0, "What branch would you like to fetch?")
	if err != nil {
		return err
	}

	if len(args) < 2 {
		_, resolved := gitp.ResolveBranch(remote)
		if resolved {
			return nil
		}
	}

	// Assume remote was just remote name - get branch name
	branch, err := common.ReqArg(args, 0, fmt.Sprintf(
		"What branch would you like to fetch from remote %q?", remote))
	if err != nil {
		return err
	}

	// Add remote if necessary
	remoteExists, err := git.RemoteExists("", remote)
	if err != nil {
		return fmt.Errorf("could not determine if remote %q exists: %w", remote, err)
	}
	if !remoteExists {
		err = addRemote(remote, "")
		if err != nil {
			return fmt.Errorf("could not add remote %q: %w", remote, err)
		}
	}

	// Fetch branch if necessary
	err = git.Fetch("", remote, branch)
	if err != nil {
		return fmt.Errorf("could not fetch remote branch %q: %w", branch, err)
	}

	return nil
}
