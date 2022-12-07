package v1

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/checksum"
)

type protocol string

// portMapping is the port mapping configurations of a sandbox.
type PortMapping struct {
	// protocol of the port mapping.
	Protocol *protocol
	// Port number within the container.
	ContainerPort *int32
	// Port number on the host.
	HostPort *int32
}

// CheckpointData contains all types of data that can be stored in the checkpoint.
type Data struct {
	PortMappings []*PortMapping `json:"port_mappings,omitempty"`
	HostNetwork  bool           `json:"host_network,omitempty"`
}

// CheckpointData is a sample example structure to be used in test cases for checkpointing
type CheckpointData struct {
	Version  string
	Name     string
	Data     *Data
	Checksum checksum.Checksum
}

func (cp *CheckpointData) MarshalCheckpoint() ([]byte, error) {
	cp.Checksum = checksum.New(*cp.Data)
	return json.Marshal(*cp)
}

func (cp *CheckpointData) UnmarshalCheckpoint(blob []byte) error {
	return json.Unmarshal(blob, cp)
}

func (cp *CheckpointData) VerifyChecksum() error {
	return cp.Checksum.Verify(*cp.Data)
}
