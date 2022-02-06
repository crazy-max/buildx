# Build mounts `RUN --mount=...`

To use this flag set Dockerfile version to at least `1.2`

```dockerfile
# syntax=docker/dockerfile:1.3
```

`RUN --mount` allows you to create mounts that process running as part of the
build can access. This can be used to bind files from other part of the build
without copying, accessing build secrets or ssh-agent sockets, or creating cache
locations to speed up your build.

## `RUN --mount=type=bind`

!!! info
    `type=bind` is the default mount type being used if not provided

This mount type allows binding directories (read-only) in the context or in an
image to the build container.

| Option               | Description |
|----------------------|-------------|
| `target`[^1]         | Mount path. |
| `source`             | Source path in the `from`. Defaults to the root of the `from`. |
| `from`               | Build stage or image name for the root of the source. Defaults to the build context. |
| `rw`,`readwrite`     | Allow writes on the mount. Written data will be discarded. |

## `RUN --mount=type=cache`

This mount type allows the build container to cache directories for compilers
and package managers.

| Option              | Description |
|---------------------|-------------|
| `id`                | Optional ID to identify separate/different caches. Defaults to value of `target`. |
| `target`[^1]        | Mount path. |
| `ro`,`readonly`     | Read-only if set. |
| `sharing`           | One of `shared`, `private`, or `locked`. Defaults to `shared`. A `shared` cache mount can be used concurrently by multiple writers. `private` creates a new mount if there are multiple writers. `locked` pauses the second writer until the first one releases the mount. |
| `from`              | Build stage to use as a base of the cache mount. Defaults to empty directory. |
| `source`            | Subpath in the `from` to mount. Defaults to the root of the `from`. |
| `mode`              | File mode for new cache directory in octal. Default `0755`. |
| `uid`               | User ID for new cache directory. Default `0`. |
| `gid`               | Group ID for new cache directory. Default `0`. |

Contents of the cache directories persists between builder invocations without
invalidating the instruction cache. Cache mounts should only be used for better
performance. Your build should work with any contents of the cache directory as
another build may overwrite the files or GC may clean it if more storage space
is needed.

### Example: cache Go packages

```dockerfile
# syntax=docker/dockerfile:1.3
FROM golang
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build ...
```

### Example: cache apt packages

```dockerfile
# syntax=docker/dockerfile:1.3
FROM ubuntu
RUN rm -f /etc/apt/apt.conf.d/docker-clean; echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache
RUN --mount=type=cache,target=/var/cache/apt \
  --mount=type=cache,target=/var/lib/apt \
  apt update && apt-get --no-install-recommends install -y gcc
```

## `RUN --mount=type=tmpfs`

This mount type allows mounting tmpfs in the build container.

| Option              | Description |
|---------------------|-------------|
| `target`[^1]        | Mount path. |
| `size`              | Specify an upper limit on the size of the filesystem. |

## `RUN --mount=type=secret`

This mount type allows the build container to access secure files such as
private keys without baking them into the image.

| Option              | Description |
|---------------------|-------------|
| `id`                | ID of the secret. Defaults to basename of the target path. |
| `target`            | Mount path. Defaults to `/run/secrets/` + `id`. |
| `required`          | If set to `true`, the instruction errors out when the secret is unavailable. Defaults to `false`. |
| `mode`              | File mode for secret file in octal. Default `0400`. |
| `uid`               | User ID for secret file. Default `0`. |
| `gid`               | Group ID for secret file. Default `0`. |

### Example: access to S3

```dockerfile
# syntax=docker/dockerfile:1.3
FROM python:3
RUN pip install awscli
RUN --mount=type=secret,id=aws,target=/root/.aws/credentials \
  aws s3 cp s3://... ...
```

```shell
$ docker buildx build --secret id=aws,src=$HOME/.aws/credentials .
```

## `RUN --mount=type=ssh`

This mount type allows the build container to access SSH keys via SSH agents,
with support for passphrases.

| Option              | Description |
|---------------------|-----------|
| `id`                | ID of SSH agent socket or key. Defaults to "default". |
| `target`            | SSH agent socket path. Defaults to `/run/buildkit/ssh_agent.${N}`. |
| `required`          | If set to `true`, the instruction errors out when the key is unavailable. Defaults to `false`. |
| `mode`              | File mode for socket in octal. Default `0600`. |
| `uid`               | User ID for socket. Default `0`. |
| `gid`               | Group ID for socket. Default `0`. |

### Example: access to Gitlab

```dockerfile
# syntax=docker/dockerfile:1.3
FROM alpine
RUN apk add --no-cache openssh-client
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
RUN --mount=type=ssh \
  ssh -q -T git@gitlab.com 2>&1 | tee /hello
# "Welcome to GitLab, @GITLAB_USERNAME_ASSOCIATED_WITH_SSHKEY" should be printed here
# with the type of build progress is defined as `plain`.
```

```shell
$ eval $(ssh-agent)
$ ssh-add ~/.ssh/id_rsa
(Input your passphrase here)
$ docker buildx build --ssh default=$SSH_AUTH_SOCK .
```

You can also specify a path to `*.pem` file on the host directly instead of `$SSH_AUTH_SOCK`.
However, pem files with passphrases are not supported.

[^1]: Value required
