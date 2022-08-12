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

package progress

import (
	"strings"

	"github.com/moby/buildkit/client"
)

func WithPrefix(w Writer, pfx string, force bool) Writer {
	return &prefixed{
		Writer: w,
		pfx:    pfx,
		force:  force,
	}
}

type prefixed struct {
	Writer
	pfx   string
	force bool
}

func (p *prefixed) Write(v *client.SolveStatus) {
	if p.force {
		for _, v := range v.Vertexes {
			v.Name = addPrefix(p.pfx, v.Name)
		}
	}
	p.Writer.Write(v)
}

func addPrefix(pfx, name string) string {
	if strings.HasPrefix(name, "[") {
		return "[" + pfx + " " + name[1:]
	}
	return "[" + pfx + "] " + name
}
