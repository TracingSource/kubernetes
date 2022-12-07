package cloudprovider

import (
	// transitive test dependencies are not vendored by go modules
	// so we have to explicitly import them here
	_ "k8s.io/legacy-cloud-providers/vsphere/testing"
)
