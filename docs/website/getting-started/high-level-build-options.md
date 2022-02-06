# High-level build options

Buildx also aims to provide support for high-level build concepts that go beyond
invoking a single build command. We want to support building all the images in
your application together and let the users define project specific reusable
build flows that can then be easily invoked by anyone.

BuildKit efficiently handles multiple concurrent build requests and
de-duplicating work. The build commands can be combined with general-purpose
command runners (for example, `make`). However, these tools generally invoke
builds in sequence  and therefore cannot leverage the full potential of BuildKit
parallelization, or combine BuildKit’s output for the user. For this use case,
we have added a command called [`docker buildx bake`](../reference/buildx_bake.md).

The `bake` command supports building images from compose files, similar to
[`docker-compose build`](https://docs.docker.com/compose/reference/build/),
but allowing all the services to be built concurrently as part of a single
request.

There is also support for custom build rules from HCL/JSON files allowing
better code reuse and different target groups. The design of bake is in very
early stages and we are looking for feedback from users.
