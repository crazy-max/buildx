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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetAfter(t *testing.T) {
	m := New()

	m.Set("foo", "bar")
	m.Set("bar", "baz")

	ctx := context.TODO()
	v, err := m.Get(ctx, "foo", "bar")
	require.NoError(t, err)

	require.Equal(t, 2, len(v))
	require.Equal(t, "bar", v["foo"])
	require.Equal(t, "baz", v["bar"])

	v, err = m.Get(ctx, "foo")
	require.NoError(t, err)
	require.Equal(t, 1, len(v))
	require.Equal(t, "bar", v["foo"])
}

func TestTimeout(t *testing.T) {
	m := New()

	m.Set("foo", "bar")

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	_, err := m.Get(ctx, "bar")
	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestBlocking(t *testing.T) {
	m := New()

	m.Set("foo", "bar")

	go func() {
		time.Sleep(100 * time.Millisecond)
		m.Set("bar", "baz")
		time.Sleep(50 * time.Millisecond)
		m.Set("baz", "abc")
	}()

	ctx := context.TODO()
	v, err := m.Get(ctx, "foo", "bar", "baz")
	require.NoError(t, err)
	require.Equal(t, 3, len(v))
	require.Equal(t, "bar", v["foo"])
	require.Equal(t, "baz", v["bar"])
	require.Equal(t, "abc", v["baz"])
}
