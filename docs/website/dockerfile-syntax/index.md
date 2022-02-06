# BuildKit Dockerfile syntax

This page documents new BuildKit-only commands added to the [Dockerfile frontend](https://github.com/moby/buildkit/tree/master/frontend/dockerfile).

!!! info
    Original documentation can be found in [Dockerfile frontend module](https://github.com/moby/buildkit/blob/master/frontend/dockerfile/docs/syntax.md)

## Using external Dockerfile frontend

BuildKit supports loading frontends dynamically from container images. Images
for Dockerfile frontends are available at [`docker/dockerfile`](https://hub.docker.com/r/docker/dockerfile/tags/) repository.

To use the external frontend, the first line of your Dockerfile needs to be
`# syntax=docker/dockerfile:1.3` pointing to the specific image you want to use.

BuildKit also ships with Dockerfile frontend builtin, but it is recommended to
use an external image to make sure that all users use the same version on the
builder and to pick up bugfixes automatically without waiting for a new version
of BuildKit or Docker engine.

The images are published on two channels: *latest* and *labs*. The latest
channel uses semver versioning while labs uses an [incrementing number](https://github.com/moby/buildkit/issues/528).
This means the labs channel may remove a feature without incrementing the major
component of a version and you may want to pin the image to a specific revision.
Even when syntaxes change in between releases on labs channel, the old versions
are guaranteed to be backward compatible.

## Syntax

* [Build mounts `RUN --mount=...`](run-mount.md)
* [Network modes `RUN --network=...`](run-mount.md)
* [Security context `RUN --security=...`](run-security.md)
* [Here-Documents](heredocs.md)

## Built-in build args

| Arg                                   | Type    | Description |
|---------------------------------------|---------|-------------|
| `BUILDKIT_CACHE_MOUNT_NS`             | String  | Set optional cache ID namespace. |
| `BUILDKIT_CONTEXT_KEEP_GIT_DIR`       | Bool    | Trigger git context to keep the `.git` directory. |
| `BUILDKIT_INLINE_BUILDINFO_ATTRS`     | Bool    | Inline build info attributes in image config or not. |
| `BUILDKIT_INLINE_CACHE`               | Bool    | Inline cache metadata to image config or not. |
| `BUILDKIT_MULTI_PLATFORM`             | Bool    | Opt into determnistic output regardless of multi-platform output or not. |
| `BUILDKIT_SANDBOX_HOSTNAME`           | String  | Set the hostname (default `buildkitsandbox`) |
| `BUILDKIT_SYNTAX`                     | String  | Set frontend image |
