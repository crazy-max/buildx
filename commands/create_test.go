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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsvToMap(t *testing.T) {
	d := []string{
		"\"tolerations=key=foo,value=bar;key=foo2,value=bar2\",replicas=1",
		"namespace=default",
	}
	r, err := csvToMap(d)

	require.NoError(t, err)

	require.Contains(t, r, "tolerations")
	require.Equal(t, r["tolerations"], "key=foo,value=bar;key=foo2,value=bar2")

	require.Contains(t, r, "replicas")
	require.Equal(t, r["replicas"], "1")

	require.Contains(t, r, "namespace")
	require.Equal(t, r["namespace"], "default")
}
