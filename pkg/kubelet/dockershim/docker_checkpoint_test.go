package dockershim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodSandboxCheckpoint(t *testing.T) {
	data := &CheckpointData{HostNetwork: true}
	checkpoint := NewPodSandboxCheckpoint("ns1", "sandbox1", data)
	version, name, namespace, _, hostNetwork := checkpoint.GetData()
	assert.Equal(t, schemaVersion, version)
	assert.Equal(t, "ns1", namespace)
	assert.Equal(t, "sandbox1", name)
	assert.Equal(t, true, hostNetwork)
}
