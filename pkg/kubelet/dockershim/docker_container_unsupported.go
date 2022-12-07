// +build !windows

package dockershim

import (
	dockertypes "github.com/docker/docker/api/types"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type containerCleanupInfo struct{}

// applyPlatformSpecificDockerConfig applies platform-specific configurations to a dockertypes.ContainerCreateConfig struct.
// The containerCleanupInfo struct it returns will be passed as is to performPlatformSpecificContainerCleanup
// after either the container creation has failed or the container has been removed.
func (ds *dockerService) applyPlatformSpecificDockerConfig(*runtimeapi.CreateContainerRequest, *dockertypes.ContainerCreateConfig) (*containerCleanupInfo, error) {
	return nil, nil
}

// performPlatformSpecificContainerCleanup is responsible for doing any platform-specific cleanup
// after either the container creation has failed or the container has been removed.
func (ds *dockerService) performPlatformSpecificContainerCleanup(cleanupInfo *containerCleanupInfo) (errors []error) {
	return
}

// platformSpecificContainerInitCleanup is called when dockershim
// is starting, and is meant to clean up any cruft left by previous runs
// creating containers.
// Errors are simply logged, but don't prevent dockershim from starting.
func (ds *dockerService) platformSpecificContainerInitCleanup() (errors []error) {
	return
}
