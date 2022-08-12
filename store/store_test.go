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

package store

import (
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestEmptyStartup(t *testing.T) {
	t.Parallel()
	tmpdir, err := os.MkdirTemp("", "buildx-store")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	s, err := New(tmpdir)
	require.NoError(t, err)

	txn, close, err := s.Txn()
	require.NoError(t, err)
	defer close()

	ng, err := txn.Current("foo")
	require.NoError(t, err)
	require.Nil(t, ng)
}

func TestNodeLocking(t *testing.T) {
	t.Parallel()
	tmpdir, err := os.MkdirTemp("", "buildx-store")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	s, err := New(tmpdir)
	require.NoError(t, err)

	_, release, err := s.Txn()
	require.NoError(t, err)

	ready := make(chan struct{})

	go func() {
		_, release, err := s.Txn()
		require.NoError(t, err)
		release()
		close(ready)
	}()

	select {
	case <-time.After(100 * time.Millisecond):
	case <-ready:
		require.Fail(t, "transaction should have waited")
	}

	release()
	select {
	case <-time.After(200 * time.Millisecond):
		require.Fail(t, "transaction should have completed")
	case <-ready:
	}
}

func TestNodeManagement(t *testing.T) {
	t.Parallel()
	tmpdir, err := os.MkdirTemp("", "buildx-store")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	s, err := New(tmpdir)
	require.NoError(t, err)

	txn, release, err := s.Txn()
	require.NoError(t, err)
	defer release()

	err = txn.Save(&NodeGroup{
		Name:   "foo/bar",
		Driver: "driver",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid name")

	err = txn.Save(&NodeGroup{
		Name:   "mybuild",
		Driver: "mydriver",
	})
	require.NoError(t, err)

	ng, err := txn.NodeGroupByName("mybuild")
	require.NoError(t, err)
	require.Equal(t, "mybuild", ng.Name)
	require.Equal(t, "mydriver", ng.Driver)

	_, err = txn.NodeGroupByName("mybuild2")
	require.Error(t, err)
	require.True(t, os.IsNotExist(errors.Cause(err)))

	err = txn.Save(&NodeGroup{
		Name:   "mybuild2",
		Driver: "mydriver2",
	})
	require.NoError(t, err)

	ng, err = txn.NodeGroupByName("mybuild2")
	require.NoError(t, err)
	require.Equal(t, "mybuild2", ng.Name)
	require.Equal(t, "mydriver2", ng.Driver)

	// update existing
	err = txn.Save(&NodeGroup{
		Name:   "mybuild",
		Driver: "mydriver-mod",
	})
	require.NoError(t, err)

	ng, err = txn.NodeGroupByName("mybuild")
	require.NoError(t, err)
	require.Equal(t, "mybuild", ng.Name)
	require.Equal(t, "mydriver-mod", ng.Driver)

	ngs, err := txn.List()
	require.NoError(t, err)
	require.Equal(t, 2, len(ngs))

	// test setting current
	err = txn.SetCurrent("foo", "mybuild", false, false)
	require.NoError(t, err)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)

	ng, err = txn.Current("bar")
	require.NoError(t, err)
	require.Nil(t, ng)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.Nil(t, ng)

	// set with default
	err = txn.SetCurrent("foo", "mybuild", false, true)
	require.NoError(t, err)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)

	ng, err = txn.Current("bar")
	require.NoError(t, err)
	require.Nil(t, ng)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)

	err = txn.SetCurrent("foo", "mybuild2", false, true)
	require.NoError(t, err)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild2", ng.Name)

	err = txn.SetCurrent("bar", "mybuild", false, false)
	require.NoError(t, err)

	ng, err = txn.Current("bar")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild2", ng.Name)

	// set global
	err = txn.SetCurrent("foo", "mybuild2", true, false)
	require.NoError(t, err)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild2", ng.Name)

	ng, err = txn.Current("bar")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild2", ng.Name)

	err = txn.SetCurrent("bar", "mybuild", false, false)
	require.NoError(t, err)

	ng, err = txn.Current("bar")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.Nil(t, ng)

	err = txn.SetCurrent("bar", "mybuild", false, true)
	require.NoError(t, err)

	err = txn.SetCurrent("foo", "mybuild2", false, false)
	require.NoError(t, err)

	// test removal
	err = txn.Remove("mybuild2")
	require.NoError(t, err)

	_, err = txn.NodeGroupByName("mybuild2")
	require.Error(t, err)
	require.True(t, os.IsNotExist(errors.Cause(err)))

	ng, err = txn.Current("foo")
	require.NoError(t, err)
	require.Nil(t, ng)

	ng, err = txn.Current("bar")
	require.NoError(t, err)
	require.NotNil(t, ng)
	require.Equal(t, "mybuild", ng.Name)
}
