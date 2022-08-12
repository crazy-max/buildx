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
	"time"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/identity"
	"github.com/opencontainers/go-digest"
)

type Logger func(*client.SolveStatus)

type SubLogger interface {
	Wrap(name string, fn func() error) error
	Log(stream int, dt []byte)
	SetStatus(*client.VertexStatus)
}

func Wrap(name string, l Logger, fn func(SubLogger) error) (err error) {
	dgst := digest.FromBytes([]byte(identity.NewID()))
	tm := time.Now()
	l(&client.SolveStatus{
		Vertexes: []*client.Vertex{{
			Digest:  dgst,
			Name:    name,
			Started: &tm,
		}},
	})

	defer func() {
		tm2 := time.Now()
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		l(&client.SolveStatus{
			Vertexes: []*client.Vertex{{
				Digest:    dgst,
				Name:      name,
				Started:   &tm,
				Completed: &tm2,
				Error:     errMsg,
			}},
		})
	}()

	return fn(&subLogger{dgst, l})
}

type subLogger struct {
	dgst   digest.Digest
	logger Logger
}

func (sl *subLogger) Wrap(name string, fn func() error) (err error) {
	tm := time.Now()
	sl.logger(&client.SolveStatus{
		Statuses: []*client.VertexStatus{{
			Vertex:    sl.dgst,
			ID:        name,
			Timestamp: time.Now(),
			Started:   &tm,
		}},
	})

	defer func() {
		tm2 := time.Now()
		sl.logger(&client.SolveStatus{
			Statuses: []*client.VertexStatus{{
				Vertex:    sl.dgst,
				ID:        name,
				Timestamp: time.Now(),
				Started:   &tm,
				Completed: &tm2,
			}},
		})
	}()

	return fn()
}

func (sl *subLogger) Log(stream int, dt []byte) {
	sl.logger(&client.SolveStatus{
		Logs: []*client.VertexLog{{
			Vertex:    sl.dgst,
			Stream:    stream,
			Data:      dt,
			Timestamp: time.Now(),
		}},
	})
}

func (sl *subLogger) SetStatus(st *client.VertexStatus) {
	st.Vertex = sl.dgst
	sl.logger(&client.SolveStatus{
		Statuses: []*client.VertexStatus{st},
	})
}
