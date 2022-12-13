package dockershim

import (
	"context"
	"fmt"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ReopenContainerLog reopens the container log file.
func (ds *dockerService) ReopenContainerLog(
	_ context.Context, _ *runtimeapi.ReopenContainerLogRequest,
) (*runtimeapi.ReopenContainerLogResponse, error) {
	return nil, fmt.Errorf("docker does not support reopening container log files")
}
