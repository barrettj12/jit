package main

import (
	"os"

	"github.com/barrettj12/jit/common"
)

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "clone":
		common.Execute("git-clone", args)
	case "new":
		common.Execute("git-new", args)
	case "pull":
		common.Execute("git-pull", args)
	case "rebase":
		common.Execute("git-rebase", args)
	case "rm", "remove":
		common.Execute("rm.sh", args)
	case "what":
		common.Execute("what.sh", args)
	case "where":
		common.Execute("git-where", args)
	default:
		gitArgs := append([]string{cmd}, args...)
		common.Execute("git", gitArgs)
	}
}
