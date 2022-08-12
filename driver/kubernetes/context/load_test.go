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

package context

import (
	"os"
	"testing"

	"github.com/docker/cli/cli/command"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultContextInitializer(t *testing.T) {
	os.Setenv("KUBECONFIG", "./fixtures/test-kubeconfig")
	defer os.Unsetenv("KUBECONFIG")
	ctx, err := command.ResolveDefaultContext(&cliflags.CommonOptions{}, command.DefaultContextStoreConfig())
	require.NoError(t, err)
	assert.Equal(t, "default", ctx.Meta.Name)
	assert.Equal(t, "zoinx", ctx.Meta.Endpoints[KubernetesEndpoint].(EndpointMeta).DefaultNamespace)
}
