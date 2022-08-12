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

ifneq (, $(BUILDX_BIN))
	export BUILDX_CMD = $(BUILDX_BIN)
else ifneq (, $(shell docker buildx version))
	export BUILDX_CMD = docker buildx
else ifneq (, $(shell which buildx))
	export BUILDX_CMD = $(which buildx)
else
	$(error "Buildx is required: https://github.com/docker/buildx#installing")
endif

export BIN_OUT = ./bin
export RELEASE_OUT = ./release-out

shell:
	./hack/shell

binaries:
	$(BUILDX_CMD) bake binaries

binaries-cross:
	$(BUILDX_CMD) bake binaries-cross

install: binaries
	mkdir -p ~/.docker/cli-plugins
	install bin/buildx ~/.docker/cli-plugins/docker-buildx

release:
	./hack/release

validate-all: lint test validate-vendor validate-docs

lint:
	$(BUILDX_CMD) bake lint

test:
	$(BUILDX_CMD) bake test

validate-vendor:
	$(BUILDX_CMD) bake validate-vendor

validate-docs:
	$(BUILDX_CMD) bake validate-docs

validate-authors:
	$(BUILDX_CMD) bake validate-authors

test-driver:
	./hack/test-driver

vendor:
	./hack/update-vendor

docs:
	./hack/update-docs

authors:
	$(BUILDX_CMD) bake update-authors

mod-outdated:
	$(BUILDX_CMD) bake mod-outdated

.PHONY: shell binaries binaries-cross install release validate-all lint validate-vendor validate-docs validate-authors vendor docs authors
