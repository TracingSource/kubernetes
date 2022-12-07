// +build !linux,!windows

package dockershim

import (
	"fmt"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func (ds *dockerService) getContainerStats(containerID string) (*runtimeapi.ContainerStats, error) {
	return nil, fmt.Errorf("not implemented")
}
