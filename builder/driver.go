package builder

import (
	"context"

	"github.com/docker/buildx/driver"
	ctxkube "github.com/docker/buildx/driver/kubernetes/context"
	"github.com/docker/buildx/store"
	"github.com/docker/buildx/store/storeutil"
	"github.com/docker/buildx/util/dockerutil"
	"github.com/docker/buildx/util/imagetools"
	"github.com/docker/buildx/util/platformutil"
	"github.com/moby/buildkit/util/grpcerrors"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
)

type Driver struct {
	store.Node
	Driver      driver.Driver
	Info        *driver.Info
	Platforms   []ocispecs.Platform
	ImageOpt    imagetools.Opt
	ProxyConfig map[string]string
	Version     string
	Err         error
}

// Drivers returns drivers for this builder.
func (b *Builder) Drivers() []Driver {
	return b.drivers
}

// LoadDrivers loads and returns drivers for this builder.
// TODO: this should be a method on a Driver object and lazy load data for each driver.
func (b *Builder) LoadDrivers(ctx context.Context, withData bool) (_ []Driver, err error) {
	eg, _ := errgroup.WithContext(ctx)
	b.drivers = make([]Driver, len(b.NodeGroup.Nodes))

	defer func() {
		if b.err == nil && err != nil {
			b.err = err
		}
	}()

	factory, err := b.Factory(ctx)
	if err != nil {
		return nil, err
	}

	imageopt, err := b.ImageOpt()
	if err != nil {
		return nil, err
	}

	for i, n := range b.NodeGroup.Nodes {
		func(i int, n store.Node) {
			eg.Go(func() error {
				di := Driver{
					Node:        n,
					ProxyConfig: storeutil.GetProxyConfig(b.opts.dockerCli),
				}
				defer func() {
					b.drivers[i] = di
				}()

				dockerapi, err := dockerutil.NewClientAPI(b.opts.dockerCli, n.Endpoint)
				if err != nil {
					di.Err = err
					return nil
				}

				contextStore := b.opts.dockerCli.ContextStore()

				var kcc driver.KubeClientConfig
				kcc, err = ctxkube.ConfigFromContext(n.Endpoint, contextStore)
				if err != nil {
					// err is returned if n.Endpoint is non-context name like "unix:///var/run/docker.sock".
					// try again with name="default".
					// FIXME(@AkihiroSuda): n should retain real context name.
					kcc, err = ctxkube.ConfigFromContext("default", contextStore)
					if err != nil {
						logrus.Error(err)
					}
				}

				tryToUseKubeConfigInCluster := false
				if kcc == nil {
					tryToUseKubeConfigInCluster = true
				} else {
					if _, err := kcc.ClientConfig(); err != nil {
						tryToUseKubeConfigInCluster = true
					}
				}
				if tryToUseKubeConfigInCluster {
					kccInCluster := driver.KubeClientConfigInCluster{}
					if _, err := kccInCluster.ClientConfig(); err == nil {
						logrus.Debug("using kube config in cluster")
						kcc = kccInCluster
					}
				}

				d, err := driver.GetDriver(ctx, "buildx_buildkit_"+n.Name, factory, n.Endpoint, dockerapi, imageopt.Auth, kcc, n.Flags, n.Files, n.DriverOpts, n.Platforms, b.opts.contextPathHash)
				if err != nil {
					di.Err = err
					return nil
				}
				di.Driver = d
				di.ImageOpt = imageopt

				if withData {
					if err := di.loadData(ctx); err != nil {
						di.Err = err
					}
				}
				return nil
			})
		}(i, n)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// TODO: This should be done in the routine loading driver data
	if withData {
		kubernetesDriverCount := 0
		for _, d := range b.drivers {
			if d.Info != nil && len(d.Info.DynamicNodes) > 0 {
				kubernetesDriverCount++
			}
		}

		isAllKubernetesDrivers := len(b.drivers) == kubernetesDriverCount
		if isAllKubernetesDrivers {
			var drivers []Driver
			var dynamicNodes []store.Node
			for _, di := range b.drivers {
				// dynamic nodes are used in Kubernetes driver.
				// Kubernetes' pods are dynamically mapped to BuildKit Nodes.
				if di.Info != nil && len(di.Info.DynamicNodes) > 0 {
					for i := 0; i < len(di.Info.DynamicNodes); i++ {
						diClone := di
						if pl := di.Info.DynamicNodes[i].Platforms; len(pl) > 0 {
							diClone.Platforms = pl
						}
						drivers = append(drivers, di)
					}
					dynamicNodes = append(dynamicNodes, di.Info.DynamicNodes...)
				}
			}

			// not append (remove the static nodes in the store)
			b.NodeGroup.Nodes = dynamicNodes
			b.drivers = drivers
			b.NodeGroup.Dynamic = true
		}
	}

	return b.drivers, nil
}

func (d *Driver) loadData(ctx context.Context) error {
	if d.Driver == nil {
		return nil
	}
	info, err := d.Driver.Info(ctx)
	if err != nil {
		return err
	}
	d.Info = info
	if d.Info.Status == driver.Running {
		driverClient, err := d.Driver.Client(ctx)
		if err != nil {
			return err
		}
		workers, err := driverClient.ListWorkers(ctx)
		if err != nil {
			return errors.Wrap(err, "listing workers")
		}
		for _, w := range workers {
			d.Platforms = append(d.Platforms, w.Platforms...)
		}
		d.Platforms = platformutil.Dedupe(d.Platforms)
		inf, err := driverClient.Info(ctx)
		if err != nil {
			if st, ok := grpcerrors.AsGRPCStatus(err); ok && st.Code() == codes.Unimplemented {
				d.Version, err = d.Driver.Version(ctx)
				if err != nil {
					return errors.Wrap(err, "getting version")
				}
			}
		} else {
			d.Version = inf.BuildkitVersion.Version
		}
	}
	return nil
}
