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

package logutil

import (
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewFilter(levels []logrus.Level, filters ...string) logrus.Hook {
	dl := logrus.New()
	dl.SetOutput(io.Discard)
	return &logsFilter{
		levels:        levels,
		filters:       filters,
		discardLogger: dl,
	}
}

type logsFilter struct {
	levels        []logrus.Level
	filters       []string
	discardLogger *logrus.Logger
}

func (d *logsFilter) Levels() []logrus.Level {
	return d.levels
}

func (d *logsFilter) Fire(entry *logrus.Entry) error {
	for _, f := range d.filters {
		if strings.Contains(entry.Message, f) {
			entry.Logger = d.discardLogger
			return nil
		}
	}
	return nil
}
