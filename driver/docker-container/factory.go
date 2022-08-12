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
	"fmt"
	"strings"

	"github.com/docker/buildx/driver"
	dockertypes "github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

const prioritySupported = 30
const priorityUnsupported = 70

func init() {
	driver.Register(&factory{})
}

type factory struct {
}

func (*factory) Name() string {
	return "docker-container"
}

func (*factory) Usage() string {
	return "docker-container"
}

func (*factory) Priority(ctx context.Context, endpoint string, api dockerclient.APIClient) int {
	if api == nil {
		return priorityUnsupported
	}
	return prioritySupported
}

func (f *factory) New(ctx context.Context, cfg driver.InitConfig) (driver.Driver, error) {
	if cfg.DockerAPI == nil {
		return nil, errors.Errorf("%s driver requires docker API access", f.Name())
	}
	d := &Driver{factory: f, InitConfig: cfg}
	dockerInfo, err := cfg.DockerAPI.Info(ctx)
	if err != nil {
		return nil, err
	}
	secOpts, err := dockertypes.DecodeSecurityOptions(dockerInfo.SecurityOptions)
	if err != nil {
		return nil, err
	}
	for _, f := range secOpts {
		if f.Name == "userns" {
			d.userNSRemap = true
			break
		}
	}
	for k, v := range cfg.DriverOpts {
		switch {
		case k == "network":
			d.netMode = v
			if v == "host" {
				d.InitConfig.BuildkitFlags = append(d.InitConfig.BuildkitFlags, "--allow-insecure-entitlement=network.host")
			}
		case k == "image":
			d.image = v
		case k == "cgroup-parent":
			d.cgroupParent = v
		case strings.HasPrefix(k, "env."):
			envName := strings.TrimPrefix(k, "env.")
			if envName == "" {
				return nil, errors.Errorf("invalid env option %q, expecting env.FOO=bar", k)
			}
			d.env = append(d.env, fmt.Sprintf("%s=%s", envName, v))
		default:
			return nil, errors.Errorf("invalid driver option %s for docker-container driver", k)
		}
	}

	return d, nil
}

func (f *factory) AllowsInstances() bool {
	return true
}
