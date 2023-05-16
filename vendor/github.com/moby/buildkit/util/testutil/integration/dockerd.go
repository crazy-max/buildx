package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/moby/buildkit/cmd/buildkitd/config"
	"github.com/moby/buildkit/util/testutil/dockerd"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// InitDockerdWorker registers a dockerd worker with the global registry.
func InitDockerdWorker() {
	dockerdUnsupported := []string{
		FeatureCacheExport,
		FeatureCacheImport,
		FeatureCacheBackendAzblob,
		FeatureCacheBackendGha,
		FeatureCacheBackendLocal,
		FeatureCacheBackendRegistry,
		FeatureCacheBackendS3,
		FeatureDirectPush,
		FeatureImageExporter,
		FeatureMultiCacheExport,
		FeatureMultiPlatform,
		FeatureOCIExporter,
		FeatureOCILayout,
		FeatureProvenance,
		FeatureSBOM,
		FeatureSecurityMode,
		FeatureCNINetwork,
	}
	Register(&Moby{
		ID:          "dockerd",
		Dockerd:     "dockerd",
		IsRootless:  false,
		Unsupported: dockerdUnsupported,
	})

	dockerdContainerdUnsupported := []string{
		FeatureSecurityMode,
		FeatureCNINetwork,
	}
	Register(&Moby{
		ID:                    "dockerd-containerd",
		Dockerd:               "dockerd",
		IsRootless:            false,
		ContainerdSnapshotter: true,
		Unsupported:           dockerdContainerdUnsupported,
	})

	// e.g. `docker-23.0=/opt/docker-alt-230/bin,docker-20.10=/opt/docker-alt-2010/bin`
	if s := os.Getenv("BUILDKIT_INTEGRATION_DOCKERD_EXTRA"); s != "" {
		entries := strings.Split(s, ",")
		for _, entry := range entries {
			pair := strings.Split(strings.TrimSpace(entry), "=")
			if len(pair) != 2 {
				panic(errors.Errorf("unexpected BUILDKIT_INTEGRATION_DOCKERD_EXTRA: %q", s))
			}
			name, bin := pair[0], pair[1]
			Register(&Moby{
				ID:         name,
				Dockerd:    filepath.Join(bin, "dockerd"),
				IsRootless: false,
				// override PATH to make sure that the expected version of the binaries are used
				ExtraEnv:    []string{fmt.Sprintf("PATH=%s:%s", bin, os.Getenv("PATH"))},
				Unsupported: dockerdUnsupported,
			})
			Register(&Moby{
				ID:                    name + "-containerd",
				Dockerd:               filepath.Join(bin, "dockerd"),
				IsRootless:            false,
				ContainerdSnapshotter: true,
				// override PATH to make sure that the expected version of the binaries are used
				ExtraEnv:    []string{fmt.Sprintf("PATH=%s:%s", bin, os.Getenv("PATH"))},
				Unsupported: dockerdContainerdUnsupported,
			})
		}
	}
}

type Moby struct {
	ID                    string
	Dockerd               string
	IsRootless            bool
	ContainerdSnapshotter bool
	ExtraEnv              []string
	Unsupported           []string
}

func (c Moby) Name() string {
	return c.ID
}

func (c Moby) Rootless() bool {
	return c.IsRootless
}

func (c Moby) New(ctx context.Context, cfg *BackendConfig) (b Backend, cl func() error, err error) {
	if err := requireRoot(); err != nil {
		return nil, nil, err
	}

	bkcfg, err := config.LoadFile(cfg.ConfigFile)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to load buildkit config file %s", cfg.ConfigFile)
	}

	dcfg := dockerd.Config{
		Features: map[string]bool{
			"containerd-snapshotter": c.ContainerdSnapshotter,
		},
	}
	if reg, ok := bkcfg.Registries["docker.io"]; ok && len(reg.Mirrors) > 0 {
		for _, m := range reg.Mirrors {
			dcfg.Mirrors = append(dcfg.Mirrors, "http://"+m)
		}
	}
	if bkcfg.Entitlements != nil {
		for _, e := range bkcfg.Entitlements {
			switch e {
			case "network.host":
				dcfg.Builder.Entitlements.NetworkHost = true
			case "security.insecure":
				dcfg.Builder.Entitlements.SecurityInsecure = true
			}
		}
	}

	dcfgdt, err := json.Marshal(dcfg)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to marshal dockerd config")
	}

	deferF := &multiCloser{}
	cl = deferF.F()

	defer func() {
		if err != nil {
			deferF.F()()
			cl = nil
		}
	}()

	var proxyGroup errgroup.Group
	deferF.append(proxyGroup.Wait)

	workDir, err := os.MkdirTemp("", "integration")
	if err != nil {
		return nil, nil, err
	}

	d, err := dockerd.NewDaemon(workDir, dockerd.WithDockerdBinary(c.Dockerd))
	if err != nil {
		return nil, nil, errors.Errorf("new daemon error: %q, %s", err, formatLogs(cfg.Logs))
	}

	dockerdConfigFile := filepath.Join(workDir, "daemon.json")
	if err := os.WriteFile(dockerdConfigFile, dcfgdt, 0644); err != nil {
		return nil, nil, err
	}

	dockerdFlags := []string{
		"--config-file", dockerdConfigFile,
		"--userland-proxy=false",
		"--debug",
	}
	if s := os.Getenv("BUILDKIT_INTEGRATION_DOCKERD_FLAGS"); s != "" {
		dockerdFlags = append(dockerdFlags, strings.Split(strings.TrimSpace(s), "\n")...)
	}

	err = d.StartWithError(cfg.Logs, c.ExtraEnv, dockerdFlags...)
	if err != nil {
		return nil, nil, err
	}
	deferF.append(d.StopWithError)

	if err := waitUnix(d.Sock(), 5*time.Second, nil); err != nil {
		return nil, nil, errors.Errorf("dockerd did not start up: %q, %s", err, formatLogs(cfg.Logs))
	}

	dockerAPI, err := client.NewClientWithOpts(client.WithHost(d.Sock()))
	if err != nil {
		return nil, nil, err
	}
	deferF.append(dockerAPI.Close)

	err = waitForAPI(ctx, dockerAPI, 5*time.Second)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "dockerd client api timed out: %s", formatLogs(cfg.Logs))
	}

	// Create a file descriptor to be used as a Unix domain socket.
	// Remove it immediately (the name will still be valid for the socket) so that
	// we don't leave files all over the users tmp tree.
	f, err := os.CreateTemp("", "buildkit-integration")
	if err != nil {
		return
	}
	localPath := f.Name()
	f.Close()
	os.Remove(localPath)

	listener, err := net.Listen("unix", localPath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "dockerd listener error: %s", formatLogs(cfg.Logs))
	}
	deferF.append(listener.Close)

	proxyGroup.Go(func() error {
		for {
			tmpConn, err := listener.Accept()
			if err != nil {
				// Ignore the error from accept which is always a system error.
				return nil
			}
			conn, err := dockerAPI.DialHijack(ctx, "/grpc", "h2c", nil)
			if err != nil {
				return err
			}

			proxyGroup.Go(func() error {
				_, err := io.Copy(conn, tmpConn)
				if err != nil {
					return err
				}
				return tmpConn.Close()
			})
			proxyGroup.Go(func() error {
				_, err := io.Copy(tmpConn, conn)
				if err != nil {
					return err
				}
				return conn.Close()
			})
		}
	})

	return backend{
		address:             "unix://" + listener.Addr().String(),
		dockerAddress:       d.Sock(),
		rootless:            c.IsRootless,
		isDockerd:           true,
		unsupportedFeatures: c.Unsupported,
	}, cl, nil
}

func waitForAPI(ctx context.Context, apiClient *client.Client, d time.Duration) error {
	step := 50 * time.Millisecond
	i := 0
	for {
		if _, err := apiClient.Ping(ctx); err == nil {
			break
		}
		i++
		if time.Duration(i)*step > d {
			return errors.New("failed to connect to /_ping endpoint")
		}
		time.Sleep(step)
	}
	return nil
}

func IsTestDockerd() bool {
	return os.Getenv("TEST_DOCKERD") == "1"
}

func IsTestDockerdMoby(sb Sandbox) bool {
	b, err := getBackend(sb)
	if err != nil {
		return false
	}
	return b.isDockerd && sb.Name() == "dockerd"
}
