package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/buildx/build"
	"github.com/docker/buildx/monitor/types"
	"github.com/pkg/errors"
)

type ExecCmd struct {
	m types.Monitor

	invokeConfig *build.InvokeConfig
	stdout       io.WriteCloser
}

func NewExecCmd(m types.Monitor, invokeConfig *build.InvokeConfig, stdout io.WriteCloser) types.Command {
	return &ExecCmd{m, invokeConfig, stdout}
}

func (cm *ExecCmd) Info() types.CommandInfo {
	return types.CommandInfo{
		Name:        "exec",
		HelpMessage: "execute a process in the interactive container",
		HelpMessageLong: `
Usage:
  exec COMMAND [ARG...]

COMMAND and ARG... will be executed in the container.
`,
	}
}

func (cm *ExecCmd) Exec(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.Errorf("command must be passed")
	}
	cfg := &build.InvokeConfig{
		Entrypoint: []string{args[1]},
		Cmd:        args[2:],
		NoCmd:      false,
		// TODO: support other options as well via flags
		Env:  cm.invokeConfig.Env,
		User: cm.invokeConfig.User,
		Cwd:  cm.invokeConfig.Cwd,
		Tty:  true,
	}
	pid := cm.m.Exec(ctx, cfg)
	fmt.Fprintf(cm.stdout, "Process %q started. Press Ctrl-a-c to switch to that process.\n", pid)
	return nil
}
