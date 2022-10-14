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
	case "clone":
		err = cmd.Clone(args)
	case "edit":
		err = cmd.Edit(args)
	case "log":
		err = cmd.Log(args)
	case "new":
		err = cmd.New(args)
	case "pull":
		err = common.Execute(filepath.Join(srcDir, "cmd/git-pull"), args...)
	case "rebase":
		err = common.Execute(filepath.Join(srcDir, "cmd/git-rebase"), args...)
	case "rm", "remove":
		err = common.Execute(filepath.Join(srcDir, "cmd/rm.sh"), args...)
	case "what":
		err = common.Execute(filepath.Join(srcDir, "cmd/what.sh"), args...)
	case "where":
		err = common.Execute(filepath.Join(srcDir, "cmd/git-where"), args...)
	default:
		err = common.Git(command, args...)
	}

	if err != nil {
		fmt.Println("ERROR:", err)
	}
}
