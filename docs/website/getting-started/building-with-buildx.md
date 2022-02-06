# Building with Buildx

Buildx is a Docker CLI plugin that extends the `docker build` command with the
full support of the features provided by [Moby BuildKit](https://github.com/moby/buildkit)
builder toolkit. It provides the same user experience as `docker build` with
many new features like creating scoped builder instances and building against
multiple nodes concurrently.

After installation, buildx can be accessed through the `docker buildx` command
with Docker 19.03. [`docker buildx build`](../reference/buildx_build.md) is the
command for starting a new build. With Docker versions older than 19.03 buildx
binary can be called directly to access the `docker buildx` subcommands.

```shell
$ docker buildx build .
[+] Building 8.4s (23/32)
 => ...
```

Buildx will always build using the BuildKit engine and does not require
`DOCKER_BUILDKIT=1` environment variable for starting builds.

The `docker buildx build` command supports features available for `docker build`,
including features such as outputs configuration, inline build caching, and
specifying target platform. In addition, Buildx also supports new features that
are not yet available for regular `docker build` like building manifest lists,
distributed caching, and exporting build results to OCI image tarballs.

Buildx is supposed to be flexible and can be run in different configurations
that are exposed through a driver concept. Currently, we support a
[`docker` driver](../reference/buildx_create.md#docker-driver) that uses
the BuildKit library bundled into the Docker daemon binary, a
[`docker-container` driver](../reference/buildx_create.md#docker-container-driver)
that automatically launches BuildKit inside a Docker container and a
[`kubernetes` driver](../reference/buildx_create.md#kubernetes-driver) to
spin up pods with defined BuildKit container image to build your images. We
plan to add more drivers in the future.

The user experience of using buildx is very similar across drivers, but there
are some features that are not currently supported by the `docker` driver,
because the BuildKit library bundled into docker daemon currently uses a
different storage component. In contrast, all images built with `docker` driver
are automatically added to the `docker images` view by default, whereas when
using other drivers the method for outputting an image needs to be selected
with `--output`.
