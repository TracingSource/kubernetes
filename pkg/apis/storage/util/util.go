package util

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/apis/storage"
	"k8s.io/kubernetes/pkg/features"
)

// DropDisabledFields removes disabled fields from the StorageClass object.
func DropDisabledFields(class, oldClass *storage.StorageClass) {
	if !utilfeature.DefaultFeatureGate.Enabled(features.ExpandPersistentVolumes) && !allowVolumeExpansionInUse(oldClass) {
		class.AllowVolumeExpansion = nil
	}
}

func allowVolumeExpansionInUse(oldClass *storage.StorageClass) bool {
	if oldClass == nil {
		return false
	}
	if oldClass.AllowVolumeExpansion != nil {
		return true
	}
	return false
}
