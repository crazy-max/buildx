# Install

Using `buildx` as a docker CLI plugin requires using Docker 19.03 or newer.
A limited set of functionality works with older versions of Docker when
invoking the binary directly.

## Windows and macOS

Docker Buildx is included in [Docker Desktop](https://docs.docker.com/desktop/)
for Windows and macOS.

## Linux packages

Docker Linux packages also include Docker Buildx when installed using the
[DEB or RPM packages](https://docs.docker.com/engine/install/).

## Manual download

!!! warning "Important"
    This section is for unattended installation of the buildx component. These
    instructions are mostly suitable for testing purposes. We do not recommend
    installing buildx using manual download in production environments as they
    will not be updated automatically with security updates.

    On Windows and macOS, we recommend that you install [Docker Desktop](https://docs.docker.com/desktop/)
    instead. For Linux, we recommend that you follow the [instructions specific for your distribution](#linux-packages).

You can also download the latest binary from the [GitHub releases page](https://github.com/docker/buildx/releases/latest).

Rename the relevant binary and copy it to the destination matching your OS:

=== ":fontawesome-brands-linux: Linux"

    ```shell
    mkdir -p "~/.docker/cli-plugins"
    wget "[[ config.repo_url ]]releases/download/[[ git.tag ]]/buildx-[[ git.tag ]].linux-amd64" -qO "~/.docker/cli-plugins/docker-buildx"
    chmod +x "~/.docker/cli-plugins/docker-buildx"
    ```

    Or copy it into one of these folders for installing it system-wide:

    * `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
    * `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

=== ":fontawesome-brands-apple: MacOS (intel)"

    ```shell
    mkdir -p "~/.docker/cli-plugins"
    wget "[[ config.repo_url ]]releases/download/[[ git.tag ]]/buildx-[[ git.tag ]].darwin-amd64" -qO "~/.docker/cli-plugins/docker-buildx"
    chmod +x "~/.docker/cli-plugins/docker-buildx"
    ```

    Or copy it into one of these folders for installing it system-wide:

    * `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
    * `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

=== ":fontawesome-brands-apple: MacOS (arm)"

    ```shell
    mkdir -p "~/.docker/cli-plugins"
    wget "[[ config.repo_url ]]releases/download/[[ git.tag ]]/buildx-[[ git.tag ]].darwin-arm64" -qO "~/.docker/cli-plugins/docker-buildx"
    chmod +x "~/.docker/cli-plugins/docker-buildx"
    ```

    Or copy it into one of these folders for installing it system-wide:

    * `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
    * `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

=== ":fontawesome-brands-windows: Windows"

    ```powershell
    mkdir -Force "~/.docker/cli-plugins"
    wget "[[ config.repo_url ]]releases/download/[[ git.tag ]]/buildx-[[ git.tag ]].windows-amd64.exe" -OutFile "~/.docker/cli-plugins/docker-buildx.exe"
    ```

    Or copy it into one of these folders for installing it system-wide:

    * `C:\ProgramData\Docker\cli-plugins`
    * `C:\Program Files\Docker\cli-plugins`

=== ":fontawesome-brands-windows: :fontawesome-brands-linux: WSL2"

    ```shell
    mkdir -p "~/.docker/cli-plugins"
    wget "[[ config.repo_url ]]releases/download/[[ git.tag ]]/buildx-[[ git.tag ]].linux-amd64" -qO "~/.docker/cli-plugins/docker-buildx"
    chmod +x "~/.docker/cli-plugins/docker-buildx"
    ```

    Or copy it into one of these folders for installing it system-wide:

    * `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
    * `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

## Dockerfile

Here is how to install and use Buildx inside a Dockerfile through the
[`docker/buildx-bin`](https://hub.docker.com/r/docker/buildx-bin) image:

```Dockerfile
FROM docker
COPY --from=docker/buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx
RUN docker buildx version
```

## Set buildx as the default builder

Running the command [`docker buildx install`](reference/buildx_install.md)
sets up docker builder command as an alias to `docker buildx build`. This
results in the ability to have `docker build` use the current buildx builder.

To remove this alias, run [`docker buildx uninstall`](reference/buildx_uninstall.md).
