package volumebinder

import (
	"time"

	v1 "k8s.io/api/core/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	storageinformers "k8s.io/client-go/informers/storage/v1"
	clientset "k8s.io/client-go/kubernetes"
	volumescheduling "k8s.io/kubernetes/pkg/controller/volume/scheduling"
)

// VolumeBinder sets up the volume binding library
type VolumeBinder struct {
	Binder volumescheduling.SchedulerVolumeBinder
}

// NewVolumeBinder sets up the volume binding library and binding queue
func NewVolumeBinder(
	client clientset.Interface,
	nodeInformer coreinformers.NodeInformer,
	csiNodeInformer storageinformers.CSINodeInformer,
	pvcInformer coreinformers.PersistentVolumeClaimInformer,
	pvInformer coreinformers.PersistentVolumeInformer,
	storageClassInformer storageinformers.StorageClassInformer,
	bindTimeout time.Duration) *VolumeBinder {

	return &VolumeBinder{
		Binder: volumescheduling.NewVolumeBinder(client, nodeInformer, csiNodeInformer, pvcInformer, pvInformer, storageClassInformer, bindTimeout),
	}
}

// NewFakeVolumeBinder sets up a fake volume binder and binding queue
func NewFakeVolumeBinder(config *volumescheduling.FakeVolumeBinderConfig) *VolumeBinder {
	return &VolumeBinder{
		Binder: volumescheduling.NewFakeVolumeBinder(config),
	}
}

// DeletePodBindings will delete the cached volume bindings for the given pod.
func (b *VolumeBinder) DeletePodBindings(pod *v1.Pod) {
	cache := b.Binder.GetBindingsCache()
	if cache != nil && pod != nil {
		cache.DeleteBindings(pod)
	}
}
