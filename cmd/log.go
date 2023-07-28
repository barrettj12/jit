package cmd

import "github.com/barrettj12/jit/common"

func Log(args []string) error {
	gitArgs := append([]string{"--oneline"}, args...)
	return common.Git("log", gitArgs...)
}
