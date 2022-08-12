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


FROM alpine:3.14 AS gen
RUN apk add --no-cache git
WORKDIR /src
RUN --mount=type=bind,target=. <<EOT
#!/usr/bin/env bash
set -e
mkdir /out
# see also ".mailmap" for how email addresses and names are deduplicated
{
  echo "# This file lists all individuals having contributed content to the repository."
  echo "# For how it is generated, see hack/dockerfiles/authors.Dockerfile."
  echo
  git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf
} > /out/AUTHORS
cat /out/AUTHORS
EOT

FROM scratch AS update
COPY --from=gen /out /

FROM gen AS validate
RUN --mount=type=bind,target=.,rw <<EOT
set -e
git add -A
cp -rf /out/* .
if [ -n "$(git status --porcelain -- AUTHORS)" ]; then
  echo >&2 'ERROR: Authors result differs. Please update with "make authors"'
  git status --porcelain -- AUTHORS
  exit 1
fi
EOT
