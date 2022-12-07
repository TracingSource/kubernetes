package nodevolumelimits

import (
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/informers"
	storagelisters "k8s.io/client-go/listers/storage/v1"
	kubefeatures "k8s.io/kubernetes/pkg/features"
)

// getCSINodeListerIfEnabled returns the CSINode lister or nil if the feature is disabled
func getCSINodeListerIfEnabled(factory informers.SharedInformerFactory) storagelisters.CSINodeLister {
	if !utilfeature.DefaultFeatureGate.Enabled(kubefeatures.CSINodeInfo) {
		return nil
	}
	return factory.Storage().V1().CSINodes().Lister()
}
