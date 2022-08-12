# syntax=docker/dockerfile:1.4

# Copyright 2022 Docker Buildx authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


ARG GO_VERSION=1.18
ARG MODOUTDATED_VERSION=v0.8.0

FROM golang:${GO_VERSION}-alpine AS base
RUN apk add --no-cache git rsync
WORKDIR /src

FROM base AS vendored
RUN --mount=target=/context \
  --mount=target=.,type=tmpfs  \
  --mount=target=/go/pkg/mod,type=cache <<EOT
set -e
rsync -a /context/. .
go mod tidy -compat=1.17
go mod vendor
mkdir /out
cp -r go.mod go.sum vendor /out
EOT

FROM scratch AS update
COPY --from=vendored /out /out

FROM vendored AS validate
RUN --mount=target=/context \
  --mount=target=.,type=tmpfs <<EOT
set -e
rsync -a /context/. .
git add -A
rm -rf vendor
cp -rf /out/* .
if [ -n "$(git status --porcelain -- go.mod go.sum vendor)" ]; then
  echo >&2 'ERROR: Vendor result differs. Please vendor your package with "make vendor"'
  git status --porcelain -- go.mod go.sum vendor
  exit 1
fi
EOT

FROM psampaz/go-mod-outdated:${MODOUTDATED_VERSION} AS go-mod-outdated
FROM base AS outdated
RUN --mount=target=.,ro \
  --mount=target=/go/pkg/mod,type=cache \
  --mount=from=go-mod-outdated,source=/home/go-mod-outdated,target=/usr/bin/go-mod-outdated \
  go list -mod=readonly -u -m -json all | go-mod-outdated -update -direct
