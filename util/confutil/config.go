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

package confutil

import (
	"os"
	"path"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ConfigDir will look for correct configuration store path;
// if `$BUILDX_CONFIG` is set - use it, otherwise use parent directory
// of Docker config file (i.e. `${DOCKER_CONFIG}/buildx`)
func ConfigDir(dockerCli command.Cli) string {
	if buildxConfig := os.Getenv("BUILDX_CONFIG"); buildxConfig != "" {
		logrus.Debugf("using config store %q based in \"$BUILDX_CONFIG\" environment variable", buildxConfig)
		return buildxConfig
	}

	buildxConfig := filepath.Join(filepath.Dir(dockerCli.ConfigFile().Filename), "buildx")
	logrus.Debugf("using default config store %q", buildxConfig)
	return buildxConfig
}

// DefaultConfigFile returns the default BuildKit configuration file path
func DefaultConfigFile(dockerCli command.Cli) (string, bool) {
	f := path.Join(ConfigDir(dockerCli), "buildkitd.default.toml")
	if _, err := os.Stat(f); err == nil {
		return f, true
	}
	return "", false
}

// loadConfigTree loads BuildKit config toml tree
func loadConfigTree(fp string) (*toml.Tree, error) {
	f, err := os.Open(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to load config from %s", fp)
	}
	defer f.Close()
	t, err := toml.LoadReader(f)
	if err != nil {
		return t, errors.Wrap(err, "failed to parse config")
	}
	return t, nil
}
