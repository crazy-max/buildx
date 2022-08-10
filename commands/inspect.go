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

	b, err := builder.New(dockerCli, in.builder, nil)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	err = b.LoadDrivers(timeoutCtx, true, "")

	var bootBuilder *builder.Builder
	if in.bootstrap {
		var ok bool
		ok, err = b.Boot(ctx)
		if err != nil {
			return err
		}
		bootBuilder = b
		if ok {
			err = b.LoadDrivers(timeoutCtx, true, "")
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintf(w, "Name:\t%s\n", b.NodeGroup.Name)
	fmt.Fprintf(w, "Driver:\t%s\n", b.NodeGroup.Driver)

	if err != nil {
		fmt.Fprintf(w, "Error:\t%s\n", err.Error())
	} else if b.Err != nil {
		fmt.Fprintf(w, "Error:\t%s\n", b.Err.Error())
	}
	if err == nil {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Nodes:")

		for i, n := range b.NodeGroup.Nodes {
			if i != 0 {
				fmt.Fprintln(w, "")
			}
			fmt.Fprintf(w, "Name:\t%s\n", n.Name)
			fmt.Fprintf(w, "Endpoint:\t%s\n", n.Endpoint)

			var driverOpts []string
			for k, v := range n.DriverOpts {
				driverOpts = append(driverOpts, fmt.Sprintf("%s=%q", k, v))
			}
			if len(driverOpts) > 0 {
				fmt.Fprintf(w, "Driver Options:\t%s\n", strings.Join(driverOpts, " "))
			}

			if err := b.Drivers[i].Err; err != nil {
				fmt.Fprintf(w, "Error:\t%s\n", err.Error())
			} else if bootBuilder != nil && len(bootBuilder.Drivers) > i && bootBuilder.Drivers[i].Err != nil {
				fmt.Fprintf(w, "Error:\t%s\n", bootBuilder.Drivers[i].Err.Error())
			} else {
				fmt.Fprintf(w, "Status:\t%s\n", b.Drivers[i].Info.Status)
				if len(n.Flags) > 0 {
					fmt.Fprintf(w, "Flags:\t%s\n", strings.Join(n.Flags, " "))
				}
				fmt.Fprintf(w, "Platforms:\t%s\n", strings.Join(platformutil.FormatInGroups(n.Platforms, b.Drivers[i].Platforms), ", "))
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
