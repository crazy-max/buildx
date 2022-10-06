package mocks

//go:generate mockery --srcpkg=github.com/docker/cli/cli/command --name=Cli --structname=DockerCli --filename=docker_cli.go --with-expecter --output=.
//go:generate mockery --srcpkg=github.com/docker/docker/client --name=APIClient --structname=DockerAPI --filename=docker_api.go --with-expecter --output=.
