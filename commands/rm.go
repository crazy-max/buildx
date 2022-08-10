package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/buildx/builder"
	"github.com/docker/buildx/store"
	"github.com/docker/buildx/store/storeutil"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type rmOptions struct {
	builder     string
	keepState   bool
	keepDaemon  bool
	allInactive bool
	force       bool
}

const (
	rmInactiveWarning = `WARNING! This will remove all builders that are not in running state. Are you sure you want to continue?`
)

func runRm(dockerCli command.Cli, in rmOptions) error {
	ctx := appcontext.Context()

	if in.allInactive && !in.force && !command.PromptForConfirmation(dockerCli.In(), dockerCli.Out(), rmInactiveWarning) {
		return nil
	}

	txn, release, err := storeutil.GetStore(dockerCli)
	if err != nil {
		return err
	}
	defer release()

	if in.allInactive {
		return rmAllInactive(ctx, txn, dockerCli, in)
	}

	b, err := builder.New(dockerCli, in.builder, txn)
	if err != nil {
		return err
	}
	if b.NodeGroup == nil {
		return nil
	}
	if err = b.LoadDrivers(ctx, false, ""); err != nil {
		return err
	}

	if cb := b.GetContextName(); cb != "" {
		return errors.Errorf("context builder cannot be removed, run `docker context rm %s` to remove this context", cb)
	}

	err1 := rm(ctx, b.Drivers, in)
	if err := txn.Remove(b.NodeGroup.Name); err != nil {
		return err
	}
	if err1 != nil {
		return err1
	}

	_, _ = fmt.Fprintf(dockerCli.Err(), "%s removed\n", b.NodeGroup.Name)
	return nil
}

func rmCmd(dockerCli command.Cli, rootOpts *rootOptions) *cobra.Command {
	var options rmOptions

	cmd := &cobra.Command{
		Use:   "rm [NAME]",
		Short: "Remove a builder instance",
		Args:  cli.RequiresMaxArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.builder = rootOpts.builder
			if len(args) > 0 {
				if options.allInactive {
					return errors.New("cannot specify builder name when --all-inactive is set")
				}
				options.builder = args[0]
			}
			return runRm(dockerCli, options)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&options.keepState, "keep-state", false, "Keep BuildKit state")
	flags.BoolVar(&options.keepDaemon, "keep-daemon", false, "Keep the buildkitd daemon running")
	flags.BoolVar(&options.allInactive, "all-inactive", false, "Remove all inactive builders")
	flags.BoolVarP(&options.force, "force", "f", false, "Do not prompt for confirmation")

	return cmd
}

func rm(ctx context.Context, drivers []builder.Driver, in rmOptions) (err error) {
	for _, di := range drivers {
		if di.Driver == nil {
			continue
		}
		// Do not stop the buildkitd daemon when --keep-daemon is provided
		if !in.keepDaemon {
			if err := di.Driver.Stop(ctx, true); err != nil {
				return err
			}
		}
		if err := di.Driver.Rm(ctx, true, !in.keepState, !in.keepDaemon); err != nil {
			return err
		}
		if di.Err != nil {
			err = di.Err
		}
	}
	return err
}

func rmAllInactive(ctx context.Context, txn *store.Txn, dockerCli command.Cli, in rmOptions) error {
	builders, err := builder.GetBuilders(dockerCli, txn)
	if err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	eg, _ := errgroup.WithContext(timeoutCtx)
	for _, b := range builders {
		func(b *builder.Builder) {
			eg.Go(func() error {
				if err := b.LoadDrivers(timeoutCtx, true, ""); err != nil {
					return errors.Wrapf(err, "cannot load %s", b.NodeGroup.Name)
				}
				if cb := b.GetContextName(); cb != "" {
					return errors.Errorf("context builder cannot be removed, run `docker context rm %s` to remove this context", cb)
				}
				if b.NodeGroup.Dynamic {
					return nil
				}
				if b.Inactive() {
					rmerr := rm(ctx, b.Drivers, in)
					if err := txn.Remove(b.NodeGroup.Name); err != nil {
						return err
					}
					_, _ = fmt.Fprintf(dockerCli.Err(), "%s removed\n", b.NodeGroup.Name)
					return rmerr
				}
				return nil
			})
		}(b)
	}

	return eg.Wait()
}
