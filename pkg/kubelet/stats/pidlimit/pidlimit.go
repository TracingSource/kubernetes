package pidlimit

import (
	"k8s.io/api/core/v1"
)

const (
	// PIDs is the (internal) name for this resource
	PIDs v1.ResourceName = "pid"
)
