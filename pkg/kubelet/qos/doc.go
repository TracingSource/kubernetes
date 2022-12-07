// Package qos contains helper functions for quality of service.
// For each resource (memory, CPU) Kubelet supports three classes of containers.
// Memory guaranteed containers will receive the highest priority and will get all the resources
// they need.
// Burstable containers will be guaranteed their request and can “burst” and use more resources
// when available.
// Best-Effort containers, which don’t specify a request, can use resources only if not being used
// by other pods.
package qos // import "k8s.io/kubernetes/pkg/kubelet/qos"
