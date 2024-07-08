package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config key[=value]",
		Short: "set or get a config value",
		RunE:  Config,
	}

	// Set flags
	cmd.Flags().Bool("global", false, "set/get from global config")

	return cmd
}

func Config(cmd *cobra.Command, args []string) error {
	// Three modes:
	// - config key - get key
	// - config key value - set key to value
	// - config key=value - set key to value

	if len(args) > 1 {
		return setConfig(args[0], args[1])
	}
	if strings.Contains(args[0], "=") {
		split := strings.SplitN(args[0], "=", 0)
		return setConfig(split[0], split[1])
	}
	return getConfig(args[0])
}

func getConfig(key string) error {}

func setConfig(key, value string) error {}
