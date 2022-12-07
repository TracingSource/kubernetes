package remote

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	internalapi "k8s.io/cri-api/pkg/apis"
	apitest "k8s.io/cri-api/pkg/apis/testing"
	fakeremote "k8s.io/kubernetes/pkg/kubelet/remote/fake"
)

const (
	defaultConnectionTimeout = 15 * time.Second
)

// createAndStartFakeRemoteRuntime creates and starts fakeremote.RemoteRuntime.
// It returns the RemoteRuntime, endpoint on success.
// Users should call fakeRuntime.Stop() to cleanup the server.
func createAndStartFakeRemoteRuntime(t *testing.T) (*fakeremote.RemoteRuntime, string) {
	endpoint, err := fakeremote.GenerateEndpoint()
	require.NoError(t, err)

	fakeRuntime := fakeremote.NewFakeRemoteRuntime()
	fakeRuntime.Start(endpoint)

	return fakeRuntime, endpoint
}

func createRemoteRuntimeService(endpoint string, t *testing.T) internalapi.RuntimeService {
	runtimeService, err := NewRemoteRuntimeService(endpoint, defaultConnectionTimeout)
	require.NoError(t, err)

	return runtimeService
}

func TestVersion(t *testing.T) {
	fakeRuntime, endpoint := createAndStartFakeRemoteRuntime(t)
	defer fakeRuntime.Stop()

	r := createRemoteRuntimeService(endpoint, t)
	version, err := r.Version(apitest.FakeVersion)
	assert.NoError(t, err)
	assert.Equal(t, apitest.FakeVersion, version.Version)
	assert.Equal(t, apitest.FakeRuntimeName, version.RuntimeName)
}
