package container

import (
	"fmt"
	"time"

	"k8s.io/klog"
)

// Specified a policy for garbage collecting containers.
type ContainerGCPolicy struct {
	// Minimum age at which a container can be garbage collected, zero for no limit.
	MinAge time.Duration

	// Max number of dead containers any single pod (UID, container name) pair is
	// allowed to have, less than zero for no limit.
	MaxPerPodContainer int

	// Max number of total dead containers, less than zero for no limit.
	MaxContainers int
}

// ContainerGC 由 realContainerGC{} 结构体实现
//
// Manages garbage collection of dead containers.
//
// Implementation is thread-compatible.
type ContainerGC interface {
	// Garbage collect containers.
	GarbageCollect() error
	// Deletes all unused containers, including containers belonging to pods
	// that are terminated but not deleted
	DeleteAllUnusedContainers() error
}

// SourcesReadyProvider knows how to determine if configuration sources are ready
type SourcesReadyProvider interface {
	// AllReady returns true if the currently configured sources have all been seen.
	AllReady() bool
}

// TODO(vmarmol): Preferentially remove pod infra containers.
type realContainerGC struct {
	// pkg/kubelet/kuberuntime/kuberuntime_manager.go -> kubeGenericRuntimeManager{}
	//
	// Container runtime
	runtime Runtime

	// Policy for garbage collection.
	policy ContainerGCPolicy

	// sourcesReadyProvider provides the readiness of kubelet configuration sources.
	sourcesReadyProvider SourcesReadyProvider
}

// ContainerGC ...
//
// caller: 
// 	1. pkg/kubelet/kubelet.go -> NewMainKubelet()
//
// New ContainerGC instance with the specified policy.
func NewContainerGC(
	runtime Runtime, policy ContainerGCPolicy, sourcesReadyProvider SourcesReadyProvider,
) (ContainerGC, error) {
	if policy.MinAge < 0 {
		return nil, fmt.Errorf("invalid minimum garbage collection age: %v", policy.MinAge)
	}

	return &realContainerGC{
		runtime:              runtime,
		policy:               policy,
		sourcesReadyProvider: sourcesReadyProvider,
	}, nil
}

// GarbageCollect ...
//
// caller:
// 	1. pkg/kubelet/kubelet.go -> Kubelet.StartGarbageCollection()
func (cgc *realContainerGC) GarbageCollect() error {
	return cgc.runtime.GarbageCollect(cgc.policy, cgc.sourcesReadyProvider.AllReady(), false)
}

func (cgc *realContainerGC) DeleteAllUnusedContainers() error {
	klog.Infof("attempting to delete unused containers")
	return cgc.runtime.GarbageCollect(
		cgc.policy, cgc.sourcesReadyProvider.AllReady(), true,
	)
}
