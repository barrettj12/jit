package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common/git"
	"github.com/barrettj12/jit/common/types"
	"github.com/spf13/cobra"
	"os"

	"github.com/barrettj12/jit/common"
)

// TODO: rename this to `info` ?
var whatCmd = &cobra.Command{
	Use:   "what <branch>",
	Short: "Get info about a branch",
	RunE:  What,
}

func What(cmd *cobra.Command, args []string) error {
	localBranch, err := common.ReqArg(args, 0, "Branch to get info for:")
	if err != nil {
		return err
	}

	remoteBranch, err := git.PushTarget(types.LocalBranch(localBranch))
	if err != nil {
		fmt.Printf("WARNING: couldn't get remote: %v\n", err)
		// Just default to the local branch name
		defaultRemote, _ := common.DefaultRemote()
		remoteBranch = types.RemoteBranch{
			Remote: defaultRemote,
			Branch: localBranch,
		}
		fmt.Printf("assuming remote branch is %q\n\n", remoteBranch)
	}

	res := common.Exec(common.ExecArgs{
		Cmd: "gh",
		Args: []string{
			"pr", "view", fmt.Sprintf("%s:%s", remoteBranch.Remote, remoteBranch.Branch),
			"--json", "title,state,headRefName,baseRefName,url", "-t", `
{{.title}}
{{.state}}: {{.headRefName}} -> {{.baseRefName}}
{{.url}}
`[1:],
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	return res.RunError
}
