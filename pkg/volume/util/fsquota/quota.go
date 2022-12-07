package fsquota

import (
	"k8s.io/utils/mount"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/features"
)

// Interface -- quota interface
type Interface interface {
	// Does the path provided support quotas, and if so, what types
	SupportsQuotas(m mount.Interface, path string) (bool, error)
	// Assign a quota (picked by the quota mechanism) to a path,
	// and return it.
	AssignQuota(m mount.Interface, path string, poduid types.UID, bytes *resource.Quantity) error

	// Get the quota-based storage consumption for the path
	GetConsumption(path string) (*resource.Quantity, error)

	// Get the quota-based inode consumption for the path
	GetInodes(path string) (*resource.Quantity, error)

	// Remove the quota from a path
	// Implementations may assume that any data covered by the
	// quota has already been removed.
	ClearQuota(m mount.Interface, path string) error
}

func enabledQuotasForMonitoring() bool {
	return utilfeature.DefaultFeatureGate.Enabled(features.LocalStorageCapacityIsolationFSQuotaMonitoring)
}
