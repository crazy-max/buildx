# buildx du

```
docker buildx du
```

<!---MARKER_GEN_START-->
Disk usage

### Options

| Name | Type | Description |
| --- | --- | --- |
| [`--builder`](#builder) | `string` | Override the configured builder instance |
| `--filter` | `filter` | Provide filter values |
| `--verbose` |  | Provide a more verbose output |


<!---MARKER_GEN_END-->

## Description

Display cache disk usage for a builder instance.

```shell
docker buildx du
```
```text
ID                                              RECLAIMABLE     SIZE            LAST ACCESSED
y1l3pgcnjg5eoubuaqsxptahj                       true            454.6MB         35 seconds ago
06d5sl21ika6grg4itrdufz6c*                      true            401.4MB         8 seconds ago
b77tzd3n6nteljbpyf9glqa72*                      true            163MB           8 seconds ago
pvh3rzg32qi80hi6ohdr77aed*                      true            47.47MB         8 seconds ago
u67cbi3i43j8ltpyvhzriusbb*                      true            47.43MB         8 seconds ago
28zmf7pgnwewgerpjoiiqmhox                       true            30.4MB          8 seconds ago
k8holwbg1cipgn6e7cqaoprya                       true            21.09MB         8 seconds ago
jowerukjv6wy2iorlavbxkn28                       true            9.053MB         35 seconds ago
hnzrtedw3bixp0d6en28zepu2                       true            1.716MB         35 seconds ago
w2dsxanl5v9i8gafv6sz0dwmx                       true            91.81kB         8 seconds ago
uq9rq5g1vlxp849fgsgpqgu1u                       true            81.92kB         8 seconds ago
0r8dwcy8vo2z548ddzmpfpj2l*                      true            20.48kB         8 seconds ago
qafkabp2y3v36nejgmn2bp3ih                       true            16.54kB         8 seconds ago
4szeerohx1wam17wy0dd68t8q                       true            12.44kB         35 seconds ago
oidamimg5erfc4wiglaounyqg*                      true            8.192kB         8 seconds ago
kchwj7d09dhn2h01i46l3t0mr*                      true            8.192kB         8 seconds ago
gitl5usms3o1830qp0zdepmdi                       true            8.192kB         8 seconds ago
i3xv48uv6a3a3s4qh6yqf63jr*                      true            4.096kB         8 seconds ago
Reclaimable:    1.176GB
Total:          1.176GB
```

## Examples

### <a name="builder"></a> Override the configured builder instance (`--builder`)

Same as [`buildx --builder`](buildx.md#builder).
