package cmd

import (
	"fmt"
	"github.com/barrettj12/jit/common"
	"os"
)

func What(args []string) error {
	localBranch, err := common.ReqArg(args, 0, "Branch to get info for:")
	if err != nil {
		return err
	}

	remote, remoteBranch, err := common.PushLoc(localBranch)
	if err != nil {
		return err
	}

	res := common.Exec(common.ExecArgs{
		Cmd: "gh",
		Args: []string{
			"pr", "view", fmt.Sprintf("%s:%s", remote, remoteBranch),
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
