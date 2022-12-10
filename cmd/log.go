package cmd

//import (
//	"github.com/barrettj12/jit/common"
//	"github.com/juju/gnuflag"
//)
//
//func Log(args []string) error {
//	expand := gnuflag.Bool("expand", false, "show full log entries instead of oneline")
//	gnuflag.Parse(true)
//
//	gitArgs := gnuflag.Args()
//	if !*expand {
//		gitArgs = append(gitArgs, "--oneline")
//	}
//	return common.Git("log", gitArgs...)
//}
