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

package platformutil

import (
	"strings"

	"github.com/containerd/containerd/platforms"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

func Parse(platformsStr []string) ([]specs.Platform, error) {
	if len(platformsStr) == 0 {
		return nil, nil
	}
	out := make([]specs.Platform, 0, len(platformsStr))
	for _, s := range platformsStr {
		parts := strings.Split(s, ",")
		if len(parts) > 1 {
			p, err := Parse(parts)
			if err != nil {
				return nil, err
			}
			out = append(out, p...)
			continue
		}
		p, err := parse(s)
		if err != nil {
			return nil, err
		}
		out = append(out, platforms.Normalize(p))
	}
	return out, nil
}

func parse(in string) (specs.Platform, error) {
	if strings.EqualFold(in, "local") {
		return platforms.DefaultSpec(), nil
	}
	return platforms.Parse(in)
}

func Dedupe(in []specs.Platform) []specs.Platform {
	m := map[string]struct{}{}
	out := make([]specs.Platform, 0, len(in))
	for _, p := range in {
		p := platforms.Normalize(p)
		key := platforms.Format(p)
		if _, ok := m[key]; ok {
			continue
		}
		m[key] = struct{}{}
		out = append(out, p)
	}
	return out
}

func FormatInGroups(gg ...[]specs.Platform) []string {
	m := map[string]struct{}{}
	out := make([]string, 0, len(gg))
	for i, g := range gg {
		for _, p := range g {
			p := platforms.Normalize(p)
			key := platforms.Format(p)
			if _, ok := m[key]; ok {
				continue
			}
			m[key] = struct{}{}
			v := platforms.Format(p)
			if i == 0 {
				v += "*"
			}
			out = append(out, v)
		}
	}
	return out
}

func Format(in []specs.Platform) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	for _, p := range in {
		out = append(out, platforms.Format(p))
	}
	return out
}
