package cloud

import (
	// ensure these packages are scanned by ginkgo for e2e tests
	_ "k8s.io/kubernetes/test/e2e/cloud/gcp"
)
