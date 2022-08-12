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

package commands

import (
	"os"

	"github.com/docker/buildx/util/cobrautil"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type installOptions struct {
}

func runInstall(dockerCli command.Cli, in installOptions) error {
	dir := config.Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "could not create docker config")
	}

	cfg, err := config.Load(dir)
	if err != nil {
		return err
	}

	if cfg.Aliases == nil {
		cfg.Aliases = map[string]string{}
	}
	cfg.Aliases["builder"] = "buildx"

	if err := cfg.Save(); err != nil {
		return errors.Wrap(err, "could not write docker config")
	}
	return nil
}

func installCmd(dockerCli command.Cli) *cobra.Command {
	var options installOptions

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install buildx as a 'docker builder' alias",
		Args:  cli.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(dockerCli, options)
		},
		Hidden: true,
	}

	// hide builder persistent flag for this command
	cobrautil.HideInheritedFlags(cmd, "builder")

	return cmd
}
