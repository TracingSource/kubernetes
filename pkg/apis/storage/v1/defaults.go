package v1

import (
	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_StorageClass(obj *storagev1.StorageClass) {
	if obj.ReclaimPolicy == nil {
		obj.ReclaimPolicy = new(v1.PersistentVolumeReclaimPolicy)
		*obj.ReclaimPolicy = v1.PersistentVolumeReclaimDelete
	}

	if obj.VolumeBindingMode == nil {
		obj.VolumeBindingMode = new(storagev1.VolumeBindingMode)
		*obj.VolumeBindingMode = storagev1.VolumeBindingImmediate
	}
}
