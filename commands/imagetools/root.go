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
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

type RootOptions struct {
	Builder *string
}

func RootCmd(dockerCli command.Cli, opts RootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "imagetools",
		Short: "Commands to work on images in registry",
	}

	cmd.AddCommand(
		createCmd(dockerCli, opts),
		inspectCmd(dockerCli, opts),
	)

	return cmd
}
