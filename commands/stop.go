package commands

import (
	"context"

	"github.com/docker/buildx/builder"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"
)

type stopOptions struct {
	builder string
}

func runStop(dockerCli command.Cli, in stopOptions) error {
	ctx := appcontext.Context()

	b, err := builder.New(dockerCli,
		builder.WithName(in.builder),
		builder.WithSkippedValidation(),
	)
	if err != nil {
		return err
	}
	drivers, err := b.LoadDrivers(ctx, false)
	if err != nil {
		return err
	}

	return stop(ctx, drivers)
}

func stopCmd(dockerCli command.Cli, rootOpts *rootOptions) *cobra.Command {
	var options stopOptions

	cmd := &cobra.Command{
		Use:   "stop [NAME]",
		Short: "Stop builder instance",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.builder = rootOpts.builder
			if len(args) > 0 {
				options.builder = args[0]
			}
			return runStop(dockerCli, options)
		},
	}

	return cmd
}

func stop(ctx context.Context, drivers []builder.Driver) (err error) {
	for _, di := range drivers {
		if di.Driver != nil {
			if err := di.Driver.Stop(ctx, true); err != nil {
				return err
			}
		}
		if di.Err != nil {
			err = di.Err
		}
	}
	return err
}
