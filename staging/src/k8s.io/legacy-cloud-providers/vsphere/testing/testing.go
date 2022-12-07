package testing

import (
	// test dependencies for k8s.io/legacy-cloud-providers/vsphere
	// import this package to vendor test dependencies since go modules does not
	// vendor transitive test dependencies
	_ "github.com/vmware/govmomi/lookup/simulator"
	_ "github.com/vmware/govmomi/simulator"
	_ "github.com/vmware/govmomi/sts/simulator"
	_ "github.com/vmware/govmomi/vapi/simulator"
)
