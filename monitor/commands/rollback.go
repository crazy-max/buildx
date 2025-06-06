package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/buildx/build"
	"github.com/docker/buildx/monitor/types"
)

type RollbackCmd struct {
	m types.Monitor

	invokeConfig *build.InvokeConfig
	stdout       io.WriteCloser
}

func NewRollbackCmd(m types.Monitor, invokeConfig *build.InvokeConfig, stdout io.WriteCloser) types.Command {
	return &RollbackCmd{m, invokeConfig, stdout}
}

func (cm *RollbackCmd) Info() types.CommandInfo {
	return types.CommandInfo{
		Name:        "rollback",
		HelpMessage: "re-runs the interactive container with the step's rootfs contents",
		HelpMessageLong: `
Usage:
  rollback [FLAGS] [COMMAND] [ARG...]

Flags:
  --init Run the container with the initial rootfs of that step.

COMMAND and ARG... will be executed in the container.
`,
	}
}

func (cm *RollbackCmd) Exec(ctx context.Context, args []string) error {
	cfg := cm.invokeConfig
	if len(args) >= 2 {
		cmds := args[1:]
		if cmds[0] == "--init" {
			cfg.Initial = true
			cmds = cmds[1:]
		}
		if len(cmds) > 0 {
			cfg.Entrypoint = []string{cmds[0]}
			cfg.Cmd = cmds[1:]
			cfg.NoCmd = false
		}
	}
	id := cm.m.Rollback(ctx, cfg)
	fmt.Fprintf(cm.stdout, "Interactive container was restarted with process %q. Press Ctrl-a-c to switch to the new container\n", id)
	return nil
}
