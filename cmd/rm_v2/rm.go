package rm_v2

import (
	"fmt"
	"strings"

	"github.com/barrettj12/jit/common"
)

func RemoveV2(args []string) error {
	branch, err := common.ReqArg(args, 0, "Which branch would you like to remove?")
	if err != nil {
		return err
	}

	localBranch, err := resolveBranch(branch)
	if err != nil {
		return err
	}

	deleteUpstreamBranch(localBranch)
	deleteWorktree(localBranch)
	deleteLocalBranch(localBranch)
	return nil
}

// This is based on cmd/gitProvider.ResolveBranch, but with subtle differences
// (it only considers local branches and doesn't try to fetch).
func resolveBranch(branch string) (string, error) {
	// Try to resolve as a local branch
	_, err := common.ExecGit("", "rev-parse", branch)
	if err == nil {
		return branch, nil
	}

	// Try to resolve remote:branch or remote/branch
	split := strings.SplitN(branch, ":", 2)
	if len(split) < 2 {
		split = strings.SplitN(branch, "/", 2)
		if len(split) < 2 {
			// Can't parse this branch ref
			return "", fmt.Errorf("branch not found")
		}
	}

	remoteBranch := split[1]
	_, err = common.ExecGit("", "rev-parse", remoteBranch)
	if err == nil {
		return remoteBranch, nil
	}

	return "", fmt.Errorf("branch not found")
}

func deleteUpstreamBranch(branch string) {

}

func deleteWorktree(branch string) {

}

func deleteLocalBranch(branch string) {

}
