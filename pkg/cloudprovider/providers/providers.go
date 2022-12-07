// +build !providerless

package cloudprovider

import (
	// Cloud providers
	_ "k8s.io/legacy-cloud-providers/aws"
	_ "k8s.io/legacy-cloud-providers/azure"
	_ "k8s.io/legacy-cloud-providers/gce"
	_ "k8s.io/legacy-cloud-providers/openstack"
	_ "k8s.io/legacy-cloud-providers/vsphere"
)
