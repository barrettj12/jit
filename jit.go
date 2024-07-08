package main

import (
	"fmt"
	"github.com/barrettj12/jit/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	// Set up global flags
	baseCmd.PersistentFlags().StringP("repo", "R", "", "repo to execute commands in")

	cmd.AddSubcommands(baseCmd)

	err := baseCmd.Execute()
	if err != nil {
		fmt.Printf("ERROR %s\n", err)
		os.Exit(1)
	}
}

// baseCmd represents the base command when called without any subcommands
var baseCmd = &cobra.Command{
	Use: "jit",
}
