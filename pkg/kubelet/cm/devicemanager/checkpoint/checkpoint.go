package checkpoint

import (
	"encoding/json"

	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/checksum"
)

// DeviceManagerCheckpoint defines the operations to retrieve pod devices
type DeviceManagerCheckpoint interface {
	checkpointmanager.Checkpoint
	GetData() ([]PodDevicesEntry, map[string][]string)
}

// PodDevicesEntry connects pod information to devices
type PodDevicesEntry struct {
	PodUID        string
	ContainerName string
	ResourceName  string
	DeviceIDs     []string
	AllocResp     []byte
}

// checkpointData struct is used to store pod to device allocation information
// in a checkpoint file.
// TODO: add version control when we need to change checkpoint format.
type checkpointData struct {
	PodDeviceEntries  []PodDevicesEntry
	RegisteredDevices map[string][]string
}

// Data holds checkpoint data and its checksum
type Data struct {
	Data     checkpointData
	Checksum checksum.Checksum
}

// caller:
// 	1. pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl.readCheckpoint()
// 	2. pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl.writeCheckpoint()
// New returns an instance of Checkpoint
func New(devEntries []PodDevicesEntry,
	devices map[string][]string) DeviceManagerCheckpoint {
	return &Data{
		Data: checkpointData{
			PodDeviceEntries:  devEntries,
			RegisteredDevices: devices,
		},
	}
}

// MarshalCheckpoint returns marshalled data
func (cp *Data) MarshalCheckpoint() ([]byte, error) {
	cp.Checksum = checksum.New(cp.Data)
	return json.Marshal(*cp)
}

// UnmarshalCheckpoint returns unmarshalled data
func (cp *Data) UnmarshalCheckpoint(blob []byte) error {
	return json.Unmarshal(blob, cp)
}

// VerifyChecksum verifies that passed checksum is same as calculated checksum
func (cp *Data) VerifyChecksum() error {
	return cp.Checksum.Verify(cp.Data)
}

// GetData returns device entries and registered devices
func (cp *Data) GetData() ([]PodDevicesEntry, map[string][]string) {
	return cp.Data.PodDeviceEntries, cp.Data.RegisteredDevices
}
