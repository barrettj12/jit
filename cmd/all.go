package cmd

import "github.com/spf13/cobra"

func AddSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		addRemoteCmd,
		applyCmd,
		newCloneCmd(),
		editCmd,
		fetchCmd,
		forkCmd,
		logCmd,
		newCmd,
		publishCmd,
		pullCmd,
		rebaseCmd,
		removeCmd,
		whatCmd,
		whereCmd,
	)
}
