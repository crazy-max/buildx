# buildx imagetools inspect

```
docker buildx imagetools inspect [OPTIONS] NAME
```

<!---MARKER_GEN_START-->
Show details of image in the registry

### Options

| Name | Type | Description |
| --- | --- | --- |
| [`--builder`](#builder) | `string` | Override the configured builder instance |
| [`--raw`](#raw) |  | Show original JSON manifest |


<!---MARKER_GEN_END-->

## Description

Show details of an image in the registry.

```shell
docker buildx imagetools inspect alpine
```
```text
Name:      docker.io/library/alpine:latest
MediaType: application/vnd.docker.distribution.manifest.list.v2+json
Digest:    sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300

Manifests:
  Name:      docker.io/library/alpine:latest@sha256:e7d88de73db3d3fd9b2d63aa7f447a10fd0220b7cbf39803c803f2af9ba256b3
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/amd64

  Name:      docker.io/library/alpine:latest@sha256:e047bc2af17934d38c5a7fa9f46d443f1de3a7675546402592ef805cfa929f9d
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm/v6

  Name:      docker.io/library/alpine:latest@sha256:8483ecd016885d8dba70426fda133c30466f661bb041490d525658f1aac73822
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm/v7

  Name:      docker.io/library/alpine:latest@sha256:c74f1b1166784193ea6c8f9440263b9be6cae07dfe35e32a5df7a31358ac2060
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm64/v8

  Name:      docker.io/library/alpine:latest@sha256:2689e157117d2da668ad4699549e55eba1ceb79cb7862368b30919f0488213f4
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/386

  Name:      docker.io/library/alpine:latest@sha256:2042a492bcdd847a01cd7f119cd48caa180da696ed2aedd085001a78664407d6
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/ppc64le

  Name:      docker.io/library/alpine:latest@sha256:49e322ab6690e73a4909f787bcbdb873631264ff4a108cddfd9f9c249ba1d58e
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/s390x
```

## Examples

### <a name="builder"></a> Override the configured builder instance (`--builder`)

Same as [`buildx --builder`](buildx.md#builder).

### <a name="raw"></a> Show original, unformatted JSON manifest (`--raw`)

Use the `--raw` option to print the original JSON bytes instead of the formatted
output.
