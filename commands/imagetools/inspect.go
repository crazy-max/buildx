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
	"github.com/docker/buildx/store"
	"github.com/docker/buildx/store/storeutil"
	"github.com/docker/buildx/util/imagetools"
	"github.com/docker/cli-docs-tool/annotation"
	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type inspectOptions struct {
	builder string
	format  string
	raw     bool
}

func runInspect(dockerCli command.Cli, in inspectOptions, name string) error {
	ctx := appcontext.Context()

	if in.format != "" && in.raw {
		return errors.Errorf("format and raw cannot be used together")
	}

	txn, release, err := storeutil.GetStore(dockerCli)
	if err != nil {
		return err
	}
	defer release()

	var ng *store.NodeGroup

	if in.builder != "" {
		ng, err = storeutil.GetNodeGroup(txn, dockerCli, in.builder)
		if err != nil {
			return err
		}
	} else {
		ng, err = storeutil.GetCurrentInstance(txn, dockerCli)
		if err != nil {
			return err
		}
	}

	imageopt, err := storeutil.GetImageConfig(dockerCli, ng)
	if err != nil {
		return err
	}

	p, err := imagetools.NewPrinter(ctx, imageopt, name, in.format)
	if err != nil {
		return err
	}

	return p.Print(in.raw, dockerCli.Out())
}

func inspectCmd(dockerCli command.Cli, rootOpts RootOptions) *cobra.Command {
	var options inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] NAME",
		Short: "Show details of an image in the registry",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.builder = *rootOpts.Builder
			return runInspect(dockerCli, options, args[0])
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&options.format, "format", "", "Format the output using the given Go template")
	flags.SetAnnotation("format", annotation.DefaultValue, []string{`"{{.Manifest}}"`})

	flags.BoolVar(&options.raw, "raw", false, "Show original, unformatted JSON manifest")

	return cmd
}
