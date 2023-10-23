package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/moby/buildkit/util/testutil/integration"
	"github.com/stretchr/testify/require"
)

var createTests = []func(t *testing.T, sb integration.Sandbox){
	testCreateBootFail,
}

func testCreateBootFail(t *testing.T, sb integration.Sandbox) {
	if sb.Name() != "docker" {
		// create command does not require a specific driver
		t.Skip("only test with docker driver")
	}

	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer func() {
		_ = l.Close()
	}()

	cmd := buildxCmd(sb, withArgs("create", "--name", "foo", "--driver", "docker-container", "--bootstrap", fmt.Sprintf("tcp://localhost:%d", l.Addr().(*net.TCPAddr).Port)))
	out, err := cmd.CombinedOutput()
	require.Error(t, err, string(out))
}
