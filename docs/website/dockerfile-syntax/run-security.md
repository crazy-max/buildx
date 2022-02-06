# Security context `RUN --security=...`

To use this flag, set Dockerfile version to `labs` channel.

```dockerfile
# syntax=docker/dockerfile:1.3-labs
```

## `RUN --security=insecure`

With `--security=insecure`, builder runs the command without sandbox in insecure
mode, which allows to run flows requiring elevated privileges (e.g. containerd).
This is equivalent to running `docker run --privileged`.

!!! caution
    In order to access this feature, entitlement
    `security.insecure` should be enabled when starting the buildkitd daemon
    (`--allow-insecure-entitlement security.insecure`) and for a build request
    (`--allow security.insecure`).

### Example: check entitlements

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM ubuntu
RUN --security=insecure cat /proc/self/status | grep CapEff
```
```text
#84 0.093 CapEff:	0000003fffffffff
```

## `RUN --security=sandbox`

Default sandbox mode can be activated via `--security=sandbox`, but that is no-op.
