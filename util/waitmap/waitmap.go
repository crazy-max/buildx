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

package waitmap

import (
	"context"
	"sync"
)

type Map struct {
	mu sync.RWMutex
	m  map[string]interface{}
	ch map[string]chan struct{}
}

func New() *Map {
	return &Map{
		m:  make(map[string]interface{}),
		ch: make(map[string]chan struct{}),
	}
}

func (m *Map) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[key] = value

	if ch, ok := m.ch[key]; ok {
		if ch != nil {
			close(ch)
		}
	}
	m.ch[key] = nil
}

func (m *Map) Get(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return map[string]interface{}{}, nil
	}

	if len(keys) > 1 {
		out := make(map[string]interface{})
		for _, key := range keys {
			mm, err := m.Get(ctx, key)
			if err != nil {
				return nil, err
			}
			out[key] = mm[key]
		}
		return out, nil
	}

	key := keys[0]
	m.mu.Lock()
	ch, ok := m.ch[key]
	if !ok {
		ch = make(chan struct{})
		m.ch[key] = ch
	}

	if ch != nil {
		m.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ch:
			m.mu.Lock()
		}
	}

	res := m.m[key]
	m.mu.Unlock()

	return map[string]interface{}{key: res}, nil
}
