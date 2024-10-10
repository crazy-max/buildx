package docker

import (
	"bytes"
	"context"

	"github.com/docker/buildx/util/progress"
	"github.com/pkg/errors"
)

func (d *Driver) createDevice(ctx context.Context, dev string, l progress.SubLogger) error {
	switch dev {
	case "docker.com/gpu":
		return d.createVenusGPU(ctx, l)
	case "nvidia.com/gpu":
		return d.createNvidiaGPU(ctx, l)
	default:
		return errors.Errorf("unsupported device %q, loading custom devices not implemented yet", dev)
	}
}

func (d *Driver) createVenusGPU(ctx context.Context, l progress.SubLogger) error {
	script := `#!/bin/sh
set -e
if [ ! -d /dev/dri ]; then
  echo >&2 "No Venus GPU detected. Requires Docker Desktop with Docker VMM virtualization enabled."
  exit 1
fi
mkdir -p /etc/cdi
cat <<EOF > /etc/cdi/venus-gpu.json
cdiVersion: "0.6.0"
kind: "docker.com/gpu"
annotations:
  cdi.device.name: "Virtio-GPU Venus (Docker Desktop)"
devices:
- name: venus
  containerEdits:
    deviceNodes:
    # make this dynamic
    - path: /dev/dri/card0
    - path: /dev/dri/renderD128
EOF
`

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	err := d.run(ctx, []string{"/bin/ash"}, bytes.NewReader([]byte(script)), stdout, stderr)
	if err != nil {
		l.Log(1, stdout.Bytes())
		l.Log(2, stderr.Bytes())
		return err
	}
	return nil
}

func (d *Driver) createNvidiaGPU(ctx context.Context, l progress.SubLogger) error {
	script := `#!/bin/sh
set -e
mkdir -p /etc/cdi
cat <<EOF > /etc/cdi/nvidia-gpu.yaml
---
cdiVersion: 0.3.0
containerEdits:
  env:
  - NVIDIA_VISIBLE_DEVICES=void
  hooks:
  - args:
    - nvidia-cdi-hook
    - create-symlinks
    - --link
    - /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/nvidia-smi::/usr/bin/nvidia-smi
    hookName: createContainer
    path: /usr/bin/nvidia-cdi-hook
  - args:
    - nvidia-cdi-hook
    - update-ldcache
    - --folder
    - /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9
    - --folder
    - /usr/lib/wsl/lib
    hookName: createContainer
    path: /usr/bin/nvidia-cdi-hook
  mounts:
  - containerPath: /usr/lib/wsl/lib/libdxcore.so
    hostPath: /usr/lib/wsl/lib/libdxcore.so
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libcuda.so.1.1
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libcuda.so.1.1
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libcuda_loader.so
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libcuda_loader.so
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvdxgdmal.so.1
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvdxgdmal.so.1
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvidia-ml.so.1
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvidia-ml.so.1
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvidia-ml_loader.so
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvidia-ml_loader.so
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvidia-ptxjitcompiler.so.1
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/libnvidia-ptxjitcompiler.so.1
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/nvcubins.bin
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/nvcubins.bin
    options:
    - ro
    - nosuid
    - nodev
    - bind
  - containerPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/nvidia-smi
    hostPath: /usr/lib/wsl/drivers/nvmdi.inf_amd64_978a3b585e321cd9/nvidia-smi
    options:
    - ro
    - nosuid
    - nodev
    - bind
devices:
- containerEdits:
    deviceNodes:
    - path: /dev/dxg
  name: all
kind: nvidia.com/gpu
EOF
`

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	err := d.run(ctx, []string{"/bin/ash"}, bytes.NewReader([]byte(script)), stdout, stderr)
	if err != nil {
		l.Log(1, stdout.Bytes())
		l.Log(2, stderr.Bytes())
		return err
	}
	return nil
}
