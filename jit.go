package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/barrettj12/jit/cmd"
	"github.com/barrettj12/jit/common"
)

func main() {
	// Get source directory
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine source directory")
	}
	srcDir := filepath.Dir(file)

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	// case "echo":
	// 	err = common.Execute("echo", args)
	case "add-remote":
		err = cmd.AddRemote(args)
	case "apply":
		err = cmd.Apply(args)
	case "clone":
		err = cmd.Clone(args)
	case "edit":
		err = cmd.Edit(args)
	case "fetch":
		err = cmd.Fetch(args)
	case "fork":
		err = cmd.Fork(args)
	case "log":
		err = cmd.Log(args)
	case "new":
		err = cmd.NewV2(args)
	case "publish":
		err = cmd.Publish(args)
	case "pull":
		err = cmd.Pull(args)
	case "rebase":
		err = common.Execute(filepath.Join(srcDir, "cmd/git-rebase"), args...)
	case "rm", "remove":
		err = cmd.Remove(args)
	case "what":
		err = cmd.What(args)
	case "where":
		err = common.Execute(filepath.Join(srcDir, "cmd/git-where"), args...)
	default:
		err = common.Git(command, args...)
	}

	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1) // TODO: pass through exit code
	}
}
