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
	"net/url"
	"path/filepath"
	"strings"

	// import connhelpers for special url schemes
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer"
	_ "github.com/moby/buildkit/client/connhelper/kubepod"
	_ "github.com/moby/buildkit/client/connhelper/ssh"

	"github.com/docker/buildx/driver"
	util "github.com/docker/buildx/driver/remote/util"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

const prioritySupported = 20
const priorityUnsupported = 90

func init() {
	driver.Register(&factory{})
}

type factory struct {
}

func (*factory) Name() string {
	return "remote"
}

func (*factory) Usage() string {
	return "remote"
}

func (*factory) Priority(ctx context.Context, endpoint string, api dockerclient.APIClient) int {
	if util.IsValidEndpoint(endpoint) != nil {
		return priorityUnsupported
	}
	return prioritySupported
}

func (f *factory) New(ctx context.Context, cfg driver.InitConfig) (driver.Driver, error) {
	if len(cfg.Files) > 0 {
		return nil, errors.Errorf("setting config file is not supported for remote driver")
	}
	if len(cfg.BuildkitFlags) > 0 {
		return nil, errors.Errorf("setting buildkit flags is not supported for remote driver")
	}

	d := &Driver{
		factory:    f,
		InitConfig: cfg,
	}

	tls := &tlsOpts{}
	tlsEnabled := false
	for k, v := range cfg.DriverOpts {
		switch k {
		case "servername":
			tls.serverName = v
			tlsEnabled = true
		case "cacert":
			if !filepath.IsAbs(v) {
				return nil, errors.Errorf("non-absolute path '%s' provided for %s", v, k)
			}
			tls.caCert = v
			tlsEnabled = true
		case "cert":
			if !filepath.IsAbs(v) {
				return nil, errors.Errorf("non-absolute path '%s' provided for %s", v, k)
			}
			tls.cert = v
			tlsEnabled = true
		case "key":
			if !filepath.IsAbs(v) {
				return nil, errors.Errorf("non-absolute path '%s' provided for %s", v, k)
			}
			tls.key = v
			tlsEnabled = true
		default:
			return nil, errors.Errorf("invalid driver option %s for remote driver", k)
		}
	}

	if tlsEnabled {
		if tls.serverName == "" {
			// guess servername as hostname of target address
			uri, err := url.Parse(cfg.EndpointAddr)
			if err != nil {
				return nil, err
			}
			tls.serverName = uri.Hostname()
		}
		missing := []string{}
		if tls.caCert == "" {
			missing = append(missing, "cacert")
		}
		if tls.cert == "" {
			missing = append(missing, "cert")
		}
		if tls.key == "" {
			missing = append(missing, "key")
		}
		if len(missing) > 0 {
			return nil, errors.Errorf("tls enabled, but missing keys %s", strings.Join(missing, ", "))
		}
		d.tlsOpts = tls
	}

	return d, nil
}

func (f *factory) AllowsInstances() bool {
	return true
}
