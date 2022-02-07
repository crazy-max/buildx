# buildx inspect

```
docker buildx inspect [NAME]
```

<!---MARKER_GEN_START-->
Inspect current builder instance

### Options

| Name | Type | Description |
| --- | --- | --- |
| [`--bootstrap`](#bootstrap) |  | Ensure builder has booted before inspecting |
| [`--builder`](#builder) | `string` | Override the configured builder instance |


<!---MARKER_GEN_END-->

## Description

Shows information about the current or specified builder.

## Examples

### <a name="bootstrap"></a> Ensure that the builder is running before inspecting (`--bootstrap`)

Use the `--bootstrap` option to ensure that the builder is running before
inspecting it. If the driver is `docker-container`, then `--bootstrap` starts
the buildkit container and waits until it is operational. Bootstrapping is
automatically done during build, and therefore not necessary. The same BuildKit
container is used during the lifetime of the associated builder node (as
displayed in `buildx ls`).

### <a name="builder"></a> Override the configured builder instance (`--builder`)

Same as [`buildx --builder`](buildx.md#builder).

### Get information about a builder instance

By default, `inspect` shows information about the current builder. Specify the
name of the builder to inspect to get information about that builder.
The following example shows information about a builder instance named
`mybuilder`:

```shell
docker buildx inspect mybuilder
```
```text
Name:   builder
Driver: docker-container

Nodes:
Name:      builder0
Endpoint:  unix:///var/run/docker.sock
Status:    running
Flags:     --debug --allow-insecure-entitlement security.insecure --allow-insecure-entitlement network.host
Platforms: linux/amd64, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6

Name:      mac-mini-m1
Endpoint:  tcp://mac-mini-m1:2376
Status:    running
Platforms: linux/arm64*, linux/amd64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6

Name:      sifive
Endpoint:  ssh://sifive@1.2.3.4
Status:    running
Platforms: linux/riscv64*
```
