# buildx prune

```
docker buildx prune
```

<!---MARKER_GEN_START-->
Remove build cache

### Options

| Name | Type | Description |
| --- | --- | --- |
| `-a`, `--all` |  | Remove all unused images, not just dangling ones |
| [`--builder`](#builder) | `string` | Override the configured builder instance |
| `--filter` | `filter` | Provide filter values (e.g., `until=24h`) |
| `-f`, `--force` |  | Do not prompt for confirmation |
| `--keep-storage` | `bytes` | Amount of disk space to keep for cache |
| `--verbose` |  | Provide a more verbose output |


<!---MARKER_GEN_END-->

## Description

Removes build cache for a builder instance.

```shell
docker buildx prune
```
```text
WARNING! This will remove all dangling build cache. Are you sure you want to continue? [y/N] y
ID                                              RECLAIMABLE     SIZE            LAST ACCESSED
i3xv48uv6a3a3s4qh6yqf63jr*                      true            4.096kB         13 minutes ago
u67cbi3i43j8ltpyvhzriusbb*                      true            47.43MB         13 minutes ago
oidamimg5erfc4wiglaounyqg*                      true            8.192kB         13 minutes ago
0r8dwcy8vo2z548ddzmpfpj2l*                      true            20.48kB         13 minutes ago
b77tzd3n6nteljbpyf9glqa72*                      true            163MB           13 minutes ago
pvh3rzg32qi80hi6ohdr77aed*                      true            47.47MB         13 minutes ago
w2dsxanl5v9i8gafv6sz0dwmx                       true            91.81kB         13 minutes ago
06d5sl21ika6grg4itrdufz6c*                      true            401.4MB         13 minutes ago
kchwj7d09dhn2h01i46l3t0mr*                      true            8.192kB         13 minutes ago
gitl5usms3o1830qp0zdepmdi                       true            8.192kB         13 minutes ago
k8holwbg1cipgn6e7cqaoprya                       true            21.09MB         13 minutes ago
uq9rq5g1vlxp849fgsgpqgu1u                       true            81.92kB         13 minutes ago
qafkabp2y3v36nejgmn2bp3ih                       true            16.54kB         13 minutes ago
y1l3pgcnjg5eoubuaqsxptahj                       true            454.6MB         14 minutes ago
4szeerohx1wam17wy0dd68t8q                       true            12.44kB         14 minutes ago
hnzrtedw3bixp0d6en28zepu2                       true            1.716MB         14 minutes ago
jowerukjv6wy2iorlavbxkn28                       true            9.053MB         14 minutes ago
Total:  1.146GB
```

## Examples

### <a name="builder"></a> Override the configured builder instance (`--builder`)

Same as [`buildx --builder`](buildx.md#builder).
