package flexvolume

import (
	"k8s.io/utils/mount"
	utilstrings "k8s.io/utils/strings"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/volume"
)

type flexVolume struct {
	// driverName is the name of the plugin driverName.
	driverName string
	// Driver executable used to setup the volume.
	execPath string
	// mounter provides the interface that is used to mount the actual
	// block device.
	mounter mount.Interface
	// podName is the name of the pod, if available.
	podName string
	// podUID is the UID of the pod.
	podUID types.UID
	// podNamespace is the namespace of the pod, if available.
	podNamespace string
	// podServiceAccountName is the service account name of the pod, if available.
	podServiceAccountName string
	// volName is the name of the pod's volume.
	volName string
	// the underlying plugin
	plugin *flexVolumePlugin
	// the metric plugin
	volume.MetricsProvider
}

// volume.Volume interface

func (f *flexVolume) GetPath() string {
	name := f.driverName
	return f.plugin.host.GetPodVolumeDir(f.podUID, utilstrings.EscapeQualifiedName(name), f.volName)
}
