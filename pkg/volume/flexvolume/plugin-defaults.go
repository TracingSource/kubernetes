package flexvolume

import (
	"k8s.io/klog"

	"k8s.io/kubernetes/pkg/volume"
)

type pluginDefaults flexVolumePlugin

func logPrefix(plugin *flexVolumePlugin) string {
	return "flexVolume driver " + plugin.driverName + ": "
}

func (plugin *pluginDefaults) GetVolumeName(spec *volume.Spec) (string, error) {
	klog.V(4).Infof(logPrefix((*flexVolumePlugin)(plugin)), "using default GetVolumeName for volume ", spec.Name())
	return spec.Name(), nil
}
