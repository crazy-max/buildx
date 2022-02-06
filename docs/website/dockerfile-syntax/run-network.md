# Network modes `RUN --network=...`

```dockerfile
# syntax=docker/dockerfile:1.3
```

`RUN --network` allows control over which networking environment the command
is run in.

## `RUN --network=none`

The command is run with no network access (`lo` is still available, but is
isolated to this process)

### Example: isolating external effects

```dockerfile
# syntax=docker/dockerfile:1.3
FROM python:3.6
ADD mypackage.tgz wheels/
RUN --network=none pip install --find-links wheels mypackage
```

`pip` will only be able to install the packages provided in the tarfile, which
can be controlled by an earlier build stage.

## `RUN --network=host`

The command is run in the host's network environment (similar to
`docker build --network=host`, but on a per-instruction basis)

!!! caution
    The use of `--network=host` is protected by the `network.host` entitlement,
    which needs to be enabled when starting the buildkitd daemon
    (`--allow-insecure-entitlement network.host`) and on the build request
    (`--allow network.host`).

## `RUN --network=default`

Equivalent to not supplying a flag at all, the command is run in the default
network for the build.
