package commands

import (
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func debugShellCmd(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:    "debug-shell",
		Short:  "Start a monitor",
		Hidden: true,
	}
}
