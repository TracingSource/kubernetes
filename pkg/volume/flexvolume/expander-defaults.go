package flexvolume

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util"
)

type expanderDefaults struct {
	plugin *flexVolumePlugin
}

func newExpanderDefaults(plugin *flexVolumePlugin) *expanderDefaults {
	return &expanderDefaults{plugin}
}

func (e *expanderDefaults) ExpandVolumeDevice(spec *volume.Spec, newSize resource.Quantity, oldSize resource.Quantity) (resource.Quantity, error) {
	klog.Warning(logPrefix(e.plugin), "using default expand for volume ", spec.Name(), ", to size ", newSize, " from ", oldSize)
	return newSize, nil
}

// the defaults for NodeExpand return a generic resize indicator that will trigger the operation executor to go ahead with
// generic filesystem resize
func (e *expanderDefaults) NodeExpand(rsOpt volume.NodeResizeOptions) (bool, error) {
	klog.Warning(logPrefix(e.plugin), "using default filesystem resize for volume ", rsOpt.VolumeSpec.Name(), ", at ", rsOpt.DevicePath)
	fsVolume, err := util.CheckVolumeModeFilesystem(rsOpt.VolumeSpec)
	if err != nil {
		return false, fmt.Errorf("error checking VolumeMode: %v", err)
	}
	// if volume is not a fs file system, there is nothing for us to do here.
	if !fsVolume {
		return true, nil
	}

	_, err = util.GenericResizeFS(e.plugin.host, e.plugin.GetPluginName(), rsOpt.DevicePath, rsOpt.DeviceMountPath)
	if err != nil {
		return false, err
	}
	return true, nil
}
