package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/docker/buildx/builder"
	"github.com/docker/buildx/util/platformutil"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"
)

type inspectOptions struct {
	bootstrap bool
	builder   string
}

func runInspect(dockerCli command.Cli, in inspectOptions) error {
	ctx := appcontext.Context()

	b, err := builder.New(dockerCli,
		builder.WithName(in.builder),
		builder.WithSkippedValidation(),
	)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	drivers, err := b.LoadDrivers(timeoutCtx, true)
	if in.bootstrap {
		var ok bool
		ok, err = b.Boot(ctx)
		if err != nil {
			return err
		}
		if ok {
			drivers, err = b.LoadDrivers(timeoutCtx, true)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Name:\t%s\n", b.Name)
	fmt.Fprintf(w, "Driver:\t%s\n", b.Driver)

	if err != nil {
		fmt.Fprintf(w, "Error:\t%s\n", err.Error())
	} else if b.Err() != nil {
		fmt.Fprintf(w, "Error:\t%s\n", b.Err().Error())
	}
	if err == nil {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Nodes:")

		for i, d := range drivers {
			if i != 0 {
				fmt.Fprintln(w, "")
			}
			fmt.Fprintf(w, "Name:\t%s\n", d.Name)
			fmt.Fprintf(w, "Endpoint:\t%s\n", d.Endpoint)

			var driverOpts []string
			for k, v := range d.DriverOpts {
				driverOpts = append(driverOpts, fmt.Sprintf("%s=%q", k, v))
			}
			if len(driverOpts) > 0 {
				fmt.Fprintf(w, "Driver Options:\t%s\n", strings.Join(driverOpts, " "))
			}

			if err := d.Err; err != nil {
				fmt.Fprintf(w, "Error:\t%s\n", err.Error())
			} else {
				fmt.Fprintf(w, "Status:\t%s\n", drivers[i].Info.Status)
				if len(d.Flags) > 0 {
					fmt.Fprintf(w, "Flags:\t%s\n", strings.Join(d.Flags, " "))
				}
				if drivers[i].Version != "" {
					fmt.Fprintf(w, "Buildkit:\t%s\n", drivers[i].Version)
				}
				fmt.Fprintf(w, "Platforms:\t%s\n", strings.Join(platformutil.FormatInGroups(d.Node.Platforms, d.Platforms), ", "))
			}
		}
	}

	w.Flush()

	return nil
}

func inspectCmd(dockerCli command.Cli, rootOpts *rootOptions) *cobra.Command {
	var options inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [NAME]",
		Short: "Inspect current builder instance",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.builder = rootOpts.builder
			if len(args) > 0 {
				options.builder = args[0]
			}
			return runInspect(dockerCli, options)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&options.bootstrap, "bootstrap", false, "Ensure builder has booted before inspecting")

	return cmd
}
