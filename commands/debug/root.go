package debug

import (
	"github.com/docker/buildx/util/cobrautil"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

// DebugConfig is a user-specified configuration for the debugger.
type DebugConfig struct {
	// InvokeFlag is a flag to configure the launched debugger and the commaned executed on the debugger.
	InvokeFlag string

	// OnFlag is a flag to configure the timing of launching the debugger.
	OnFlag string
}

// DebuggableCmd is a command that supports debugger with recognizing the user-specified DebugConfig.
type DebuggableCmd interface {
	// NewDebugger returns the new *cobra.Command with support for the debugger with recognizing DebugConfig.
	NewDebugger(*DebugConfig) *cobra.Command
}

func RootCmd(dockerCli command.Cli, children ...DebuggableCmd) *cobra.Command {
	var progressMode string
	var options DebugConfig

	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Start debugger",
	}
	cobrautil.MarkCommandExperimental(cmd)

	flags := cmd.Flags()
	flags.StringVar(&options.InvokeFlag, "invoke", "", "Launch a monitor with executing specified command")
	flags.StringVar(&options.OnFlag, "on", "error", "When to launch the monitor ([always, error])")
	flags.StringVar(&progressMode, "progress", "auto", `Set type of progress output ("auto", "plain", "tty", "rawjson") for the monitor. Use plain to show container output`)

	cobrautil.MarkFlagsExperimental(flags, "invoke", "on")

	for _, c := range children {
		cmd.AddCommand(c.NewDebugger(&options))
	}

	return cmd
}
