# Registry mirror

You can define a registry mirror to use for your builds by providing a [BuildKit daemon configuration](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md)
while creating a builder with the [`--config` flags](../reference/buildx_create.md#config).

```toml title="/etc/buildkitd.toml"
debug = true
[registry."docker.io"]
  mirrors = ["mirror.gcr.io"]
```

!!! tip
    `debug = true` has been added to be able to debug requests in the BuildKit
    daemon and see if the mirror is effectively used.

Now let's [create a `docker-container` builder](../reference/buildx_create.md)
that will use this BuildKit configuration:

```shell
docker buildx create --name "mybuilder" --use --driver "docker-container" --config "/etc/buildkitd.toml"
```

Boot and [inspect `mybuilder`](../reference/buildx_inspect.md):

```shell
docker buildx inspect --bootstrap
```

Build a simple image:

```shell
docker buildx build --load . -f-<<EOF
FROM alpine
RUN echo "hello world"
EOF
```

Now let's check the builder container logs:

```text title="docker logs buildx_buildkit_mybuilder0" linenums="1" hl_lines="34 35 41"
time="2022-02-06T17:47:46Z" level=info msg="auto snapshotter: using overlayfs"
time="2022-02-06T17:47:46Z" level=warning msg="using host network as the default"
time="2022-02-06T17:47:46Z" level=info msg="found worker \"3a7s6b28bvvd2o9zwmg61k5yi\", labels=map[org.mobyproject.buildkit.worker.executor:oci org.mobyproject.buildkit.worker.hostname:b1b7845900bb org.mobyproject.buildkit.worker.snapshotter:overlayfs], platforms=[linux/amd64 linux/arm64 linux/riscv64 linux/ppc64le linux/s390x linux/386 linux/mips64le linux/mips64 linux/arm/v7 linux/arm/v6]"
time="2022-02-06T17:47:46Z" level=warning msg="skipping containerd worker, as \"/run/containerd/containerd.sock\" does not exist"
time="2022-02-06T17:47:46Z" level=info msg="found 1 workers, default=\"3a7s6b28bvvd2o9zwmg61k5yi\""
time="2022-02-06T17:47:46Z" level=warning msg="currently, only the default worker can be used."
time="2022-02-06T17:47:46Z" level=info msg="running server on /run/buildkit/buildkitd.sock"
time="2022-02-06T17:47:46Z" level=debug msg="session started" spanID=179ab8b4f2b6cb44 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:46Z" level=debug msg="session finished: <nil>" spanID=179ab8b4f2b6cb44 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:46Z" level=debug msg="session started" spanID=3887255ba87bc5a3 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="new ref for local: leywczg6isrvgr9p3d3vd63sx" span="[internal] load .dockerignore"
time="2022-02-06T17:47:47Z" level=debug msg="new ref for local: zqfzklfi3ik0imvbthgftikmn" span="[internal] load build definition from Dockerfile"
time="2022-02-06T17:47:47Z" level=debug msg="diffcopy took: 6.036047ms" span="[internal] load .dockerignore"
time="2022-02-06T17:47:47Z" level=debug msg="diffcopy took: 2.85119ms" span="[internal] load build definition from Dockerfile"
time="2022-02-06T17:47:47Z" level=debug msg="saved leywczg6isrvgr9p3d3vd63sx as local.sharedKey:context:context-.dockerignore:registry-mirror:94055aefadf84244" span="[internal] load .dockerignore"
time="2022-02-06T17:47:47Z" level=debug msg="saved zqfzklfi3ik0imvbthgftikmn as local.sharedKey:dockerfile:dockerfile:registry-mirror:94055aefadf84244" span="[internal] load build definition from Dockerfile"
time="2022-02-06T17:47:47Z" level=debug msg=resolving spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="do request" request.header.accept="application/vnd.docker.distribution.manifest.v2+json, application/vnd.docker.distribution.manifest.list.v2+json, application/vnd.oci.image.manifest.v1+json, application/vnd.oci.image.index.v1+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=HEAD spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="fetch response received" response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.content-length=1638 response.header.content-type=application/vnd.docker.distribution.manifest.list.v2+json response.header.date="Sun, 06 Feb 2022 17:47:52 GMT" response.header.docker-content-digest="sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300" response.header.docker-distribution-api-version=registry/2.0 response.header.server="Docker Registry" response.header.x-frame-options=SAMEORIGIN response.header.x-xss-protection=0 response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg=resolved desc.digest="sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="do request" request.header.accept="application/vnd.docker.distribution.manifest.list.v2+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="fetch response received" response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.content-length=1638 response.header.content-type=application/vnd.docker.distribution.manifest.list.v2+json response.header.date="Sun, 06 Feb 2022 17:47:53 GMT" response.header.docker-content-digest="sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300" response.header.docker-distribution-api-version=registry/2.0 response.header.server="Docker Registry" response.header.x-frame-options=SAMEORIGIN response.header.x-xss-protection=0 response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="do request" request.header.accept="application/vnd.docker.distribution.manifest.v2+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="do request" request.header.accept="application/vnd.docker.distribution.manifest.v2+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="fetch response received" response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.content-length=528 response.header.content-type=application/vnd.docker.distribution.manifest.v2+json response.header.date="Sun, 06 Feb 2022 17:47:53 GMT" response.header.docker-content-digest="sha256:2689e157117d2da668ad4699549e55eba1ceb79cb7862368b30919f0488213f4" response.header.docker-distribution-api-version=registry/2.0 response.header.server="Docker Registry" response.header.x-frame-options=SAMEORIGIN response.header.x-xss-protection=0 response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:47Z" level=debug msg="do request" request.header.accept="application/vnd.docker.container.image.v1+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.content-length=528 response.header.content-type=application/vnd.docker.distribution.manifest.v2+json response.header.date="Sun, 06 Feb 2022 17:47:53 GMT" response.header.docker-content-digest="sha256:e7d88de73db3d3fd9b2d63aa7f447a10fd0220b7cbf39803c803f2af9ba256b3" response.header.docker-distribution-api-version=registry/2.0 response.header.server="Docker Registry" response.header.x-frame-options=SAMEORIGIN response.header.x-xss-protection=0 response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="do request" request.header.accept="application/vnd.docker.container.image.v1+json, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.accept-ranges=bytes response.header.age=1356 response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.cache-control="public, max-age=3600" response.header.content-length=1469 response.header.content-type=application/octet-stream response.header.date="Sun, 06 Feb 2022 17:25:17 GMT" response.header.etag="\"774380abda8f4eae9a149e5d5d3efc83\"" response.header.expires="Sun, 06 Feb 2022 18:25:17 GMT" response.header.last-modified="Wed, 24 Nov 2021 21:07:57 GMT" response.header.server=UploadServer response.header.x-goog-generation=1637788077652182 response.header.x-goog-hash="crc32c=V3DSrg==" response.header.x-goog-hash.1="md5=d0OAq9qPTq6aFJ5dXT78gw==" response.header.x-goog-metageneration=1 response.header.x-goog-storage-class=STANDARD response.header.x-goog-stored-content-encoding=identity response.header.x-goog-stored-content-length=1469 response.header.x-guploader-uploadid=ADPycduqQipVAXc3tzXmTzKQ2gTT6CV736B2J628smtD1iDytEyiYCgvvdD8zz9BT1J1sASUq9pW_ctUyC4B-v2jvhIxnZTlKg response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.accept-ranges=bytes response.header.age=760 response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.cache-control="public, max-age=3600" response.header.content-length=1471 response.header.content-type=application/octet-stream response.header.date="Sun, 06 Feb 2022 17:35:13 GMT" response.header.etag="\"35d688bd15327daafcdb4d4395e616a8\"" response.header.expires="Sun, 06 Feb 2022 18:35:13 GMT" response.header.last-modified="Wed, 24 Nov 2021 21:07:12 GMT" response.header.server=UploadServer response.header.x-goog-generation=1637788032100793 response.header.x-goog-hash="crc32c=aWgRjA==" response.header.x-goog-hash.1="md5=NdaIvRUyfar8201DleYWqA==" response.header.x-goog-metageneration=1 response.header.x-goog-storage-class=STANDARD response.header.x-goog-stored-content-encoding=identity response.header.x-goog-stored-content-length=1471 response.header.x-guploader-uploadid=ADPycdtR-gJYwC7yHquIkJWFFG8FovDySvtmRnZBqlO3yVDanBXh_VqKYt400yhuf0XbQ3ZMB9IZV2vlcyHezn_Pu3a1SMMtiw response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg=fetch spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="do request" request.header.accept="application/vnd.docker.image.rootfs.diff.tar.gzip, */*" request.header.user-agent=containerd/1.5.8+unknown request.method=GET spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="fetch response received" response.header.accept-ranges=bytes response.header.age=1356 response.header.alt-svc="h3=\":443\"; ma=2592000,h3-29=\":443\"; ma=2592000,h3-Q050=\":443\"; ma=2592000,h3-Q046=\":443\"; ma=2592000,h3-Q043=\":443\"; ma=2592000,quic=\":443\"; ma=2592000; v=\"46,43\"" response.header.cache-control="public, max-age=3600" response.header.content-length=2818413 response.header.content-type=application/octet-stream response.header.date="Sun, 06 Feb 2022 17:25:17 GMT" response.header.etag="\"1d55e7be5a77c4a908ad11bc33ebea1c\"" response.header.expires="Sun, 06 Feb 2022 18:25:17 GMT" response.header.last-modified="Wed, 24 Nov 2021 21:07:06 GMT" response.header.server=UploadServer response.header.x-goog-generation=1637788026431708 response.header.x-goog-hash="crc32c=ZojF+g==" response.header.x-goog-hash.1="md5=HVXnvlp3xKkIrRG8M+vqHA==" response.header.x-goog-metageneration=1 response.header.x-goog-storage-class=STANDARD response.header.x-goog-stored-content-encoding=identity response.header.x-goog-stored-content-length=2818413 response.header.x-guploader-uploadid=ADPycdsebqxiTBJqZ0bv9zBigjFxgQydD2ESZSkKchpE0ILlN9Ibko3C5r4fJTJ4UR9ddp-UBd-2v_4eRpZ8Yo2llW_j4k8WhQ response.status="200 OK" spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="using pigz for decompression"
time="2022-02-06T17:47:48Z" level=debug msg="diff applied" d=37.803805ms digest="sha256:59bf1c3509f33515622619af21ed55bbe26d24913cedbca106468a5fb37a50c3" media=application/vnd.docker.image.rootfs.diff.tar.gzip size=2818413 spanID=9460e5b6e64cec91 traceID=b162d3040ddf86d6614e79c66a01a577
time="2022-02-06T17:47:48Z" level=debug msg="> creating koar2j6yg893h6m3h5nmdrwwx [/bin/sh -c echo \"hello world\"]" span="[2/2] RUN echo \"hello world\""
time="2022-02-06T17:47:48Z" level=debug msg="Using double walk diff for /tmp/containerd-mount239997430 from /tmp/containerd-mount395065907"
time="2022-02-06T17:47:48Z" level=debug msg="session finished: <nil>" spanID=3887255ba87bc5a3 traceID=b162d3040ddf86d6614e79c66a01a577
```

As you can see on highlighted lines, requests come from the Google registry mirror.
