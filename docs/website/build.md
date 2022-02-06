# Build

## Plugin

```shell
# Buildx 0.6+
docker buildx bake "https://github.com/docker/buildx.git"
mkdir -p ~/.docker/cli-plugins
mv ./bin/buildx ~/.docker/cli-plugins/docker-buildx

# Docker 19.03+
DOCKER_BUILDKIT=1 docker build --platform=local -o . "https://github.com/docker/buildx.git"
mkdir -p ~/.docker/cli-plugins
mv buildx ~/.docker/cli-plugins/docker-buildx

# Local
git clone https://github.com/docker/buildx.git && cd buildx
make install
```

## Docs website

```shell
# Build website and output to ./site
docker buildx bake update-website

# Runs website and watch for changes
docker buildx bake base-website
docker run --rm -it -p 8000:8000 -v $(pwd):/docs buildx-website:local
# Open http://localhost:8000 in your browser
```
