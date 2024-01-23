# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21.6
ARG XX_VERSION=1.2.1

ARG DOCKER_VERSION=24.0.6
ARG GOTESTSUM_VERSION=v1.9.0
ARG REGISTRY_VERSION=2.8.0
ARG BUILDKIT_VERSION=v0.11.6
ARG K3S_VERSION=v1.27.9-k3s1

# xx is a helper for cross-compilation
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS golatest

FROM golatest AS gobase
COPY --from=xx / /
RUN apk add --no-cache file git
ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
WORKDIR /src

FROM registry:$REGISTRY_VERSION AS registry
FROM rancher/k3s:${K3S_VERSION} AS k3s
FROM moby/buildkit:$BUILDKIT_VERSION AS buildkit

FROM gobase AS docker
ARG TARGETPLATFORM
ARG DOCKER_VERSION
WORKDIR /opt/docker
RUN DOCKER_ARCH=$(case ${TARGETPLATFORM:-linux/amd64} in \
    "linux/amd64")   echo "x86_64"  ;; \
    "linux/arm/v6")  echo "armel"   ;; \
    "linux/arm/v7")  echo "armhf"   ;; \
    "linux/arm64")   echo "aarch64" ;; \
    "linux/ppc64le") echo "ppc64le" ;; \
    "linux/s390x")   echo "s390x"   ;; \
    *)               echo ""        ;; esac) \
  && echo "DOCKER_ARCH=$DOCKER_ARCH" \
  && wget -qO- "https://download.docker.com/linux/static/stable/${DOCKER_ARCH}/docker-${DOCKER_VERSION}.tgz" | tar xvz --strip 1
RUN ./dockerd --version && ./containerd --version && ./ctr --version && ./runc --version

FROM gobase AS gotestsum
ARG GOTESTSUM_VERSION
ENV GOFLAGS=
RUN --mount=target=/root/.cache,type=cache \
  GOBIN=/out/ go install "gotest.tools/gotestsum@${GOTESTSUM_VERSION}" && \
  /out/gotestsum --version

FROM gobase AS buildx-version
RUN --mount=type=bind,target=. <<EOT
  set -e
  mkdir /buildx-version
  echo -n "$(./hack/git-meta version)" | tee /buildx-version/version
  echo -n "$(./hack/git-meta revision)" | tee /buildx-version/revision
EOT

FROM gobase AS buildx-build
ARG TARGETPLATFORM
RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,from=buildx-version,source=/buildx-version,target=/buildx-version <<EOT
  set -e
  xx-go --wrap
  DESTDIR=/usr/bin VERSION=$(cat /buildx-version/version) REVISION=$(cat /buildx-version/revision) GO_EXTRA_LDFLAGS="-s -w" ./hack/build
  xx-verify --static /usr/bin/docker-buildx
EOT

FROM gobase AS test
ENV SKIP_INTEGRATION_TESTS=1
RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  go test -v -coverprofile=/tmp/coverage.txt -covermode=atomic ./... && \
  go tool cover -func=/tmp/coverage.txt

FROM scratch AS test-coverage
COPY --from=test /tmp/coverage.txt /coverage.txt

FROM scratch AS binaries-unix
COPY --link --from=buildx-build /usr/bin/docker-buildx /buildx

FROM binaries-unix AS binaries-darwin
FROM binaries-unix AS binaries-linux

FROM scratch AS binaries-windows
COPY --link --from=buildx-build /usr/bin/docker-buildx /buildx.exe

FROM binaries-$TARGETOS AS binaries
# enable scanning for this stage
ARG BUILDKIT_SBOM_SCAN_STAGE=true

FROM gobase AS integration-test-base
# https://github.com/docker/docker/blob/master/project/PACKAGERS.md#runtime-dependencies
RUN apk add --no-cache \
      btrfs-progs \
      e2fsprogs \
      e2fsprogs-extra \
      ip6tables \
      iptables \
      openssl \
      shadow-uidmap \
      xfsprogs \
      xz
# k3s deps
RUN apk add --no-cache \
      busybox-binsh \
      cni-plugins \
      cni-plugin-flannel \
      conntrack-tools \
      coreutils \
      dbus \
      findutils \
      ipset
ENV PATH="/usr/libexec/cni:${PATH}"
COPY --link --from=gotestsum /out/gotestsum /usr/bin/
COPY --link --from=registry /bin/registry /usr/bin/
COPY --link --from=docker /opt/docker/* /usr/bin/
COPY --link --from=k3s /bin/k3s /usr/bin/
COPY --link --from=k3s /bin/kubectl /usr/bin/
COPY --link --from=buildkit /usr/bin/buildkitd /usr/bin/
COPY --link --from=buildkit /usr/bin/buildctl /usr/bin/
COPY --link --from=binaries /buildx /usr/bin/
COPY <<-"EOF" /entrypoint.sh
  #!/bin/sh
  set -e
  # cgroup v2: enable nesting
  # https://github.com/moby/moby/blob/v25.0.0/hack/dind#L59-L69
  if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
    # move the processes from the root group to the /init group,
    # otherwise writing subtree_control fails with EBUSY.
    # An error during moving non-existent process (i.e., "cat") is ignored.
    mkdir -p /sys/fs/cgroup/init
    xargs -rn1 < /sys/fs/cgroup/cgroup.procs > /sys/fs/cgroup/init/cgroup.procs || :
    # enable controllers
    sed -e 's/ / +/g' -e 's/^/+/' < /sys/fs/cgroup/cgroup.controllers > /sys/fs/cgroup/cgroup.subtree_control
  fi
  exec "$@"
EOF
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

FROM integration-test-base AS integration-test
COPY . .

# Release
FROM --platform=$BUILDPLATFORM alpine AS releaser
WORKDIR /work
ARG TARGETPLATFORM
RUN --mount=from=binaries \
  --mount=type=bind,from=buildx-version,source=/buildx-version,target=/buildx-version <<EOT
  set -e
  mkdir -p /out
  cp buildx* "/out/buildx-$(cat /buildx-version/version).$(echo $TARGETPLATFORM | sed 's/\//-/g')$(ls buildx* | sed -e 's/^buildx//')"
EOT

FROM scratch AS release
COPY --from=releaser /out/ /

# Shell
FROM docker:$DOCKER_VERSION AS dockerd-release
FROM alpine AS shell
RUN apk add --no-cache iptables tmux git vim less openssh
RUN mkdir -p /usr/local/lib/docker/cli-plugins && ln -s /usr/local/bin/buildx /usr/local/lib/docker/cli-plugins/docker-buildx
COPY ./hack/demo-env/entrypoint.sh /usr/local/bin
COPY ./hack/demo-env/tmux.conf /root/.tmux.conf
COPY --from=dockerd-release /usr/local/bin /usr/local/bin
WORKDIR /work
COPY ./hack/demo-env/examples .
COPY --from=binaries / /usr/local/bin/
VOLUME /var/lib/docker
ENTRYPOINT ["entrypoint.sh"]

FROM binaries
