package builder

import (
	"context"
	"os"
	"sort"

	"github.com/docker/buildx/driver"
	"github.com/docker/buildx/store"
	"github.com/docker/buildx/store/storeutil"
	"github.com/docker/buildx/util/imagetools"
	"github.com/docker/buildx/util/progress"
	"github.com/docker/cli/cli/command"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Builder represents an active builder object
type Builder struct {
	dockerCli command.Cli

	NodeGroup *store.NodeGroup
	Drivers   []Driver
	Err       error
}

// New initializes a new builder client
func New(dockerCli command.Cli, name string, txn *store.Txn) (_ *Builder, err error) {
	b := &Builder{
		dockerCli: dockerCli,
	}

	if txn == nil {
		// if store instance is nil we create a short-lived instance using
		// the default store and ensure we release it upon this func is completed
		var release func()
		txn, release, err = storeutil.GetStore(dockerCli)
		if err != nil {
			return nil, err
		}
		defer release()
	}

	if name != "" {
		b.NodeGroup, err = storeutil.GetNodeGroup(txn, dockerCli, name)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	b.NodeGroup, err = storeutil.GetCurrentInstance(txn, dockerCli)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetCurrent returns current builder instance.
func GetCurrent(dockerCli command.Cli) (_ *Builder, err error) {
	return New(dockerCli, "", nil)
}

// Validate validates builder context
func (b *Builder) Validate() error {
	if b.NodeGroup.Name == "default" && b.NodeGroup.Name != b.dockerCli.CurrentContext() {
		return errors.Errorf("use `docker --context=default buildx` to switch to default context")
	}
	list, err := b.dockerCli.ContextStore().List()
	if err != nil {
		return err
	}
	for _, l := range list {
		if l.Name == b.NodeGroup.Name && b.NodeGroup.Name != "default" {
			return errors.Errorf("use `docker --context=%s buildx` to switch to context %q", b.NodeGroup.Name, b.NodeGroup.Name)
		}
	}
	return nil
}

// GetImageOpt returns registry auth configuration
func (b *Builder) GetImageOpt() (imagetools.Opt, error) {
	return storeutil.GetImageConfig(b.dockerCli, b.NodeGroup)
}

// Boot bootstrap a builder
func (b *Builder) Boot(ctx context.Context) (bool, error) {
	toBoot := make([]int, 0, len(b.Drivers))
	for idx, d := range b.Drivers {
		if d.Err != nil || d.Driver == nil || d.Info == nil {
			continue
		}
		if d.Info.Status != driver.Running {
			toBoot = append(toBoot, idx)
		}
	}
	if len(toBoot) == 0 {
		return false, nil
	}

	printer := progress.NewPrinter(context.Background(), os.Stderr, os.Stderr, progress.PrinterModeAuto)

	baseCtx := ctx
	eg, _ := errgroup.WithContext(ctx)
	for _, idx := range toBoot {
		func(idx int) {
			eg.Go(func() error {
				pw := progress.WithPrefix(printer, b.NodeGroup.Nodes[idx].Name, len(toBoot) > 1)
				_, err := driver.Boot(ctx, baseCtx, b.Drivers[idx].Driver, pw)
				if err != nil {
					b.Drivers[idx].Err = err
				}
				return nil
			})
		}(idx)
	}

	err := eg.Wait()
	err1 := printer.Wait()
	if err == nil {
		err = err1
	}

	return true, err
}

// Inactive checks if all nodes are inactive for this builder.
func (b *Builder) Inactive() bool {
	for _, d := range b.Drivers {
		if d.Info != nil && d.Info.Status == driver.Running {
			return false
		}
	}
	return true
}

// GetBuilders returns all builders
func GetBuilders(dockerCli command.Cli, txn *store.Txn) ([]*Builder, error) {
	storeng, err := txn.List()
	if err != nil {
		return nil, err
	}

	currentName := "default"
	current, err := storeutil.GetCurrentInstance(txn, dockerCli)
	if err != nil {
		return nil, err
	}
	if current != nil {
		if current.Name == "default" {
			currentName = current.Nodes[0].Endpoint
		} else {
			currentName = current.Name
		}
	}

	currentSet := false

	builders := make([]*Builder, len(storeng))
	seen := make(map[string]struct{})
	for i, ng := range storeng {
		if !currentSet && ng.Name == currentName {
			ng.Current = true
			currentSet = true
		}
		builders[i] = &Builder{
			dockerCli: dockerCli,
			NodeGroup: ng,
		}
		seen[ng.Name] = struct{}{}
	}

	contexts, err := dockerCli.ContextStore().List()
	if err != nil {
		return nil, err
	}
	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Name < contexts[j].Name
	})
	for _, c := range contexts {
		// if a context has the same name as an instance from the store, do not
		// add it to the builders list. An instance from the store takes
		// precedence over context builders.
		if _, ok := seen[c.Name]; ok {
			continue
		}

		defaultNg := false
		if !currentSet && c.Name == currentName {
			defaultNg = true
			currentSet = true
		}

		builders = append(builders, &Builder{
			dockerCli: dockerCli,
			NodeGroup: &store.NodeGroup{
				Name:    c.Name,
				Current: defaultNg,
				Nodes: []store.Node{
					{
						Name:     c.Name,
						Endpoint: c.Name,
					},
				},
			},
		})
	}

	return builders, nil
}
