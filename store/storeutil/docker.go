package storeutil

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/docker"
	"github.com/docker/docker/client"
)

// ClientForEndpoint returns a docker client for an endpoint
func ClientForEndpoint(dockerCli command.Cli, name string) (client.APIClient, error) {
	dem, err := GetDockerEndpoint(dockerCli, name)
	if err == nil && dem != nil {
		ep, err := docker.WithTLSData(dockerCli.ContextStore(), name, *dem)
		if err != nil {
			return nil, err
		}
		clientOpts, err := ep.ClientOpts()
		if err != nil {
			return nil, err
		}
		return client.NewClientWithOpts(clientOpts...)
	}
	ep := docker.Endpoint{
		EndpointMeta: docker.EndpointMeta{
			Host: name,
		},
	}
	clientOpts, err := ep.ClientOpts()
	if err != nil {
		return nil, err
	}
	return client.NewClientWithOpts(clientOpts...)
}
