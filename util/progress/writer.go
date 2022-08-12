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

type Writer interface {
	Write(*client.SolveStatus)
	ValidateLogSource(digest.Digest, interface{}) bool
	ClearLogSource(interface{})
}

func Write(w Writer, name string, f func() error) {
	dgst := digest.FromBytes([]byte(identity.NewID()))
	tm := time.Now()

	vtx := client.Vertex{
		Digest:  dgst,
		Name:    name,
		Started: &tm,
	}

	w.Write(&client.SolveStatus{
		Vertexes: []*client.Vertex{&vtx},
	})

	err := f()

	tm2 := time.Now()
	vtx2 := vtx
	vtx2.Completed = &tm2
	if err != nil {
		vtx2.Error = err.Error()
	}
	w.Write(&client.SolveStatus{
		Vertexes: []*client.Vertex{&vtx2},
	})
}

func NewChannel(w Writer) (chan *client.SolveStatus, chan struct{}) {
	ch := make(chan *client.SolveStatus)
	done := make(chan struct{})
	go func() {
		for {
			v, ok := <-ch
			if !ok {
				close(done)
				w.ClearLogSource(done)
				return
			}

			if len(v.Logs) > 0 {
				logs := make([]*client.VertexLog, 0, len(v.Logs))
				for _, l := range v.Logs {
					if w.ValidateLogSource(l.Vertex, done) {
						logs = append(logs, l)
					}
				}
				v.Logs = logs
			}

			w.Write(v)
		}
	}()
	return ch, done
}
