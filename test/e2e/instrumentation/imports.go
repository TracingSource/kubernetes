package instrumentation

import (
	// ensure libs have a chance to perform initialization
	_ "k8s.io/kubernetes/test/e2e/instrumentation/logging"
	_ "k8s.io/kubernetes/test/e2e/instrumentation/monitoring"
)
