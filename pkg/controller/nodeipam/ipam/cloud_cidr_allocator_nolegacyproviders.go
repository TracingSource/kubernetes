// +build providerless

package ipam

import (
	"errors"

	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
)

// NewCloudCIDRAllocator creates a new cloud CIDR allocator.
func NewCloudCIDRAllocator(client clientset.Interface, cloud cloudprovider.Interface, nodeInformer informers.NodeInformer) (CIDRAllocator, error) {
	return nil, errors.New("legacy cloud provider support not built")
}
