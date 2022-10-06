package commands

import (
	"testing"

	"github.com/docker/buildx/util/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRunVersion(t *testing.T) {
	api := mocks.NewDockerAPI(t)
	cli := mocks.NewDockerCli(t)
	cli.EXPECT().Client().Return(api)

	assert.NoError(t, runVersion(cli))
}
