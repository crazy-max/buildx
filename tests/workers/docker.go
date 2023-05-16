package workers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/util/testutil/integration"
	"github.com/pkg/errors"
)

func InitDockerWorker() {
	integration.Register(&dockerWorker{
		id:         "docker",
		dockerBin:  "docker",
		dockerdBin: "dockerd",
	})
	// defined in Dockerfile
	// e.g. `docker-23.0=/opt/docker-alt-230/bin,docker-20.10=/opt/docker-alt-2010/bin`
	if s := os.Getenv("BUILDX_INTEGRATION_DOCKER_EXTRA"); s != "" {
		entries := strings.Split(s, ",")
		for _, entry := range entries {
			pair := strings.Split(strings.TrimSpace(entry), "=")
			if len(pair) != 2 {
				panic(errors.Errorf("unexpected BUILDX_INTEGRATION_DOCKER_EXTRA: %q", s))
			}
			name, bin := pair[0], pair[1]
			integration.Register(&dockerWorker{
				id:         name,
				dockerBin:  filepath.Join(bin, "docker"),
				dockerdBin: filepath.Join(bin, "dockerd"),
				// override PATH to make sure that the expected version of the binaries are used
				extraEnv: []string{fmt.Sprintf("PATH=%s:%s", bin, os.Getenv("PATH"))},
			})
		}
	}
}

type dockerWorker struct {
	id         string
	dockerBin  string
	dockerdBin string
	extraEnv   []string // e.g. "PATH=/opt/docker-alt-230/bin:/usr/bin:..."
}

func (c dockerWorker) Name() string {
	return c.id
}

func (c dockerWorker) Rootless() bool {
	return false
}

func (c dockerWorker) New(ctx context.Context, cfg *integration.BackendConfig) (b integration.Backend, cl func() error, err error) {
	deferF := &multiCloser{}
	cl = deferF.F()

	defer func() {
		if err != nil {
			deferF.F()()
			cl = nil
		}
	}()

	moby := integration.Moby{
		ID:       c.id,
		Dockerd:  c.dockerdBin,
		ExtraEnv: c.extraEnv,
	}
	bk, bkclose, err := moby.New(ctx, cfg)
	if err != nil {
		return bk, cl, err
	}
	deferF.append(bkclose)

	tmpdir, err := os.MkdirTemp("", "buildxtest_docker")
	if err != nil {
		return nil, nil, err
	}
	deferF.append(func() error { return os.RemoveAll(tmpdir) })

	name := "integration-" + identity.NewID()
	cmd := exec.Command(c.dockerBin, "context", "create",
		name,
		"--config",
		"--docker", "host="+bk.DockerAddress(),
	)
	cmd.Env = append(os.Environ(), "DOCKER_CONTEXT="+bk.DockerAddress())
	if err := cmd.Run(); err != nil {
		return nil, cl, errors.Wrapf(err, "failed to create buildx instance %s", name)
	}
	deferF.append(func() error {
		cmd := exec.Command("docker", "context", "rm", "-f", name)
		return cmd.Run()
	})

	return &backend{
		builder: name,
		context: name,
	}, cl, nil
}
