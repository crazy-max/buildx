# syntax=docker/dockerfile:1

ARG GO_VERSION=1.19
ARG MOCKERY_VERSION=v2.14.0

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src

FROM vektra/mockery:${MOCKERY_VERSION} AS mockery
FROM base AS generate
RUN --mount=type=bind,target=.,rw \
    --mount=from=mockery,source=/usr/local/bin/mockery,target=/usr/local/bin/mockery <<EOT
  set -e
  go generate ./...
  mkdir -p /out/util
  cp -Rf util/mocks /out/util
EOT

FROM scratch AS update
COPY --from=generate /out /

FROM generate AS validate
RUN apk add --no-cache git
RUN --mount=type=bind,target=.,rw <<EOT
  set -e
  git add -A
  cp -rf /out/* .
  diff=$(git status --porcelain -- util/mocks)
  if [ -n "$diff" ]; then
    echo >&2 'ERROR: Vendor result differs. Please vendor your package with "make generate"'
    echo "$diff"
    exit 1
  fi
EOT
