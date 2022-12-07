// +build linux

package dockershim

import (
	"context"
	"time"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func (ds *dockerService) getContainerStats(containerID string) (*runtimeapi.ContainerStats, error) {
	info, err := ds.client.Info()
	if err != nil {
		return nil, err
	}

	statsJSON, err := ds.client.GetContainerStats(containerID)
	if err != nil {
		return nil, err
	}

	containerJSON, err := ds.client.InspectContainerWithSize(containerID)
	if err != nil {
		return nil, err
	}

	statusResp, err := ds.ContainerStatus(context.Background(), &runtimeapi.ContainerStatusRequest{ContainerId: containerID})
	if err != nil {
		return nil, err
	}
	status := statusResp.GetStatus()

	dockerStats := statsJSON.Stats
	timestamp := time.Now().UnixNano()
	containerStats := &runtimeapi.ContainerStats{
		Attributes: &runtimeapi.ContainerAttributes{
			Id:          containerID,
			Metadata:    status.Metadata,
			Labels:      status.Labels,
			Annotations: status.Annotations,
		},
		Cpu: &runtimeapi.CpuUsage{
			Timestamp:            timestamp,
			UsageCoreNanoSeconds: &runtimeapi.UInt64Value{Value: dockerStats.CPUStats.CPUUsage.TotalUsage},
		},
		Memory: &runtimeapi.MemoryUsage{
			Timestamp:       timestamp,
			WorkingSetBytes: &runtimeapi.UInt64Value{Value: dockerStats.MemoryStats.PrivateWorkingSet},
		},
		WritableLayer: &runtimeapi.FilesystemUsage{
			Timestamp: timestamp,
			FsId:      &runtimeapi.FilesystemIdentifier{Mountpoint: info.DockerRootDir},
			UsedBytes: &runtimeapi.UInt64Value{Value: uint64(*containerJSON.SizeRw)},
		},
	}
	return containerStats, nil
}
