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

package buildflags

import (
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
)

func ParseEntitlements(in []string) ([]entitlements.Entitlement, error) {
	out := make([]entitlements.Entitlement, 0, len(in))
	for _, v := range in {
		switch v {
		case "security.insecure":
			out = append(out, entitlements.EntitlementSecurityInsecure)
		case "network.host":
			out = append(out, entitlements.EntitlementNetworkHost)
		default:
			return nil, errors.Errorf("invalid entitlement: %v", v)
		}
	}
	return out, nil
}
