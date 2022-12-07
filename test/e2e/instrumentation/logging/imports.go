package logging

import (
	_ "k8s.io/kubernetes/test/e2e/instrumentation/logging/elasticsearch" // for elasticsearch provider
	_ "k8s.io/kubernetes/test/e2e/instrumentation/logging/stackdriver"   // for stackdriver provider
)
