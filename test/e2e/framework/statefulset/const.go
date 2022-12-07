package statefulset

import (
	"time"
)

const (
	// StatefulSetPoll is a poll interval for StatefulSet tests
	StatefulSetPoll = 10 * time.Second
	// StatefulSetTimeout is a timeout interval for StatefulSet operations
	StatefulSetTimeout = 10 * time.Minute
	// StatefulPodTimeout is a timeout for stateful pods to change state
	StatefulPodTimeout = 5 * time.Minute
)
