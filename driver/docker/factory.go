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

package docker

import (
	"context"

	"github.com/docker/buildx/driver"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

const prioritySupported = 10
const priorityUnsupported = 99

func init() {
	driver.Register(&factory{})
}

type factory struct {
}

func (*factory) Name() string {
	return "docker"
}

func (*factory) Usage() string {
	return "docker"
}

func (*factory) Priority(ctx context.Context, endpoint string, api dockerclient.APIClient) int {
	if api == nil {
		return priorityUnsupported
	}

	c, err := api.DialHijack(ctx, "/grpc", "h2c", nil)
	if err != nil {
		return priorityUnsupported
	}
	c.Close()

	return prioritySupported
}

func (f *factory) New(ctx context.Context, cfg driver.InitConfig) (driver.Driver, error) {
	if cfg.DockerAPI == nil {
		return nil, errors.Errorf("docker driver requires docker API access")
	}
	if len(cfg.Files) > 0 {
		return nil, errors.Errorf("setting config file is not supported for docker driver, use dockerd configuration file")
	}

	return &Driver{factory: f, InitConfig: cfg}, nil
}

func (f *factory) AllowsInstances() bool {
	return false
}
