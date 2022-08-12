// Copyright 2022 Docker Buildx authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package remote

import (
	"context"
	"time"

	"github.com/docker/buildx/driver"
	"github.com/docker/buildx/util/progress"
	"github.com/moby/buildkit/client"
)

type Driver struct {
	factory driver.Factory
	driver.InitConfig
	*tlsOpts
}

type tlsOpts struct {
	serverName string
	caCert     string
	cert       string
	key        string
}

func (d *Driver) Bootstrap(ctx context.Context, l progress.Logger) error {
	for i := 0; ; i++ {
		info, err := d.Info(ctx)
		if err != nil {
			return err
		}
		if info.Status != driver.Inactive {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if i > 10 {
				i = 10
			}
			time.Sleep(time.Duration(i) * time.Second)
		}
	}
}

func (d *Driver) Info(ctx context.Context) (*driver.Info, error) {
	c, err := d.Client(ctx)
	if err != nil {
		return &driver.Info{
			Status: driver.Inactive,
		}, nil
	}

	if _, err := c.ListWorkers(ctx); err != nil {
		return &driver.Info{
			Status: driver.Inactive,
		}, nil
	}

	return &driver.Info{
		Status: driver.Running,
	}, nil
}

func (d *Driver) Version(ctx context.Context) (string, error) {
	return "", nil
}

func (d *Driver) Stop(ctx context.Context, force bool) error {
	return nil
}

func (d *Driver) Rm(ctx context.Context, force, rmVolume, rmDaemon bool) error {
	return nil
}

func (d *Driver) Client(ctx context.Context) (*client.Client, error) {
	opts := []client.ClientOpt{}
	if d.tlsOpts != nil {
		opts = append(opts, client.WithCredentials(d.tlsOpts.serverName, d.tlsOpts.caCert, d.tlsOpts.cert, d.tlsOpts.key))
	}

	return client.New(ctx, d.InitConfig.EndpointAddr, opts...)
}

func (d *Driver) Features() map[driver.Feature]bool {
	return map[driver.Feature]bool{
		driver.OCIExporter:    true,
		driver.DockerExporter: false,
		driver.CacheExport:    true,
		driver.MultiPlatform:  true,
	}
}

func (d *Driver) Factory() driver.Factory {
	return d.factory
}

func (d *Driver) IsMobyDriver() bool {
	return false
}

func (d *Driver) Config() driver.InitConfig {
	return d.InitConfig
}
