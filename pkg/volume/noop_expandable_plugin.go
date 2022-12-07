package volume

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

type noopExpandableVolumePluginInstance struct {
	spec *Spec
}

var _ ExpandableVolumePlugin = &noopExpandableVolumePluginInstance{}

func (n *noopExpandableVolumePluginInstance) ExpandVolumeDevice(spec *Spec, newSize resource.Quantity, oldSize resource.Quantity) (resource.Quantity, error) {
	return newSize, nil
}

func (n *noopExpandableVolumePluginInstance) Init(host VolumeHost) error {
	return nil
}

func (n *noopExpandableVolumePluginInstance) GetPluginName() string {
	return n.spec.KubeletExpandablePluginName()
}

func (n *noopExpandableVolumePluginInstance) GetVolumeName(spec *Spec) (string, error) {
	return n.spec.Name(), nil
}

func (n *noopExpandableVolumePluginInstance) CanSupport(spec *Spec) bool {
	return true
}

func (n *noopExpandableVolumePluginInstance) RequiresRemount() bool {
	return false
}

func (n *noopExpandableVolumePluginInstance) NewMounter(spec *Spec, podRef *v1.Pod, opts VolumeOptions) (Mounter, error) {
	return nil, nil
}

func (n *noopExpandableVolumePluginInstance) NewUnmounter(name string, podUID types.UID) (Unmounter, error) {
	return nil, nil
}

func (n *noopExpandableVolumePluginInstance) ConstructVolumeSpec(volumeName, mountPath string) (*Spec, error) {
	return n.spec, nil
}

func (n *noopExpandableVolumePluginInstance) SupportsMountOption() bool {
	return true
}

func (n *noopExpandableVolumePluginInstance) SupportsBulkVolumeVerification() bool {
	return false
}

func (n *noopExpandableVolumePluginInstance) RequiresFSResize() bool {
	return true
}
