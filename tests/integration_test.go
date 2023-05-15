package tests

import (
	"testing"

	"github.com/docker/buildx/tests/workers"
	"github.com/moby/buildkit/util/testutil/integration"
)

func init() {
	if integration.IsTestDockerd() {
		workers.InitDockerWorker()
		workers.InitDockerContainerWorker()
	} else {
		workers.InitRemoteWorker()
	}
}

func TestIntegration(t *testing.T) {
	var tests []func(t *testing.T, sb integration.Sandbox)
	tests = append(tests, buildTests...)
	tests = append(tests, inspectTests...)
	tests = append(tests, lsTests...)
	testIntegration(t, tests...)
}

func testIntegration(t *testing.T, funcs ...func(t *testing.T, sb integration.Sandbox)) {
	mirroredImages := integration.OfficialImages("busybox:latest", "alpine:latest")
	mirroredImages["moby/buildkit:buildx-stable-1"] = "docker.io/moby/buildkit:buildx-stable-1"
	mirrors := integration.WithMirroredImages(mirroredImages)

	tests := integration.TestFuncs(funcs...)
	integration.Run(t, tests, mirrors)
}
