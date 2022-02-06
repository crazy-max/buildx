# Working with builder instances

By default, buildx will initially use the [`docker` driver](../reference/buildx_create.md)
if it is supported, providing a very similar user experience to the native
`docker build`. Note that you must use a local shared daemon to build your applications.

Buildx allows you to create new instances of isolated builders. This can be
used for getting a scoped environment for your CI builds that does not change
the state of the shared daemon or for isolating the builds for different
projects. You can create a new instance for a set of remote nodes, forming a
build farm, and quickly switch between them.

You can create new instances using the [`docker buildx create`](../reference/buildx_create.md)
command. This creates a new builder instance with a single node based on your
current configuration.

To use a remote node you can specify the `DOCKER_HOST` or the remote context name
while creating the new builder. After creating a new instance, you can manage its
lifecycle using the [`docker buildx inspect`](../reference/buildx_inspect.md),
[`docker buildx stop`](../reference/buildx_stop.md), and
[`docker buildx rm`](../reference/buildx_rm.md) commands. To list all
available builders, use [`buildx ls`](../reference/buildx_ls.md). After
creating a new builder you can also append new nodes to it.

To switch between different builders, use [`docker buildx use <name>`](../reference/buildx_use.md).
After running this command, the build commands will automatically use this
builder.

Docker also features a [`docker context`](https://docs.docker.com/engine/reference/commandline/context/)
command that can be used for giving names for remote Docker API endpoints.
Buildx integrates with `docker context` so that all of your contexts
automatically get a default builder instance. While creating a new builder
instance or when adding a node to it, you can also set the context name as the
target.
