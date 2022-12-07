// +build windows

package dockershim

import (
	"context"
	"time"

	"k8s.io/klog"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/winstats"
)

// ImageFsInfo returns information of the filesystem that is used to store images.
func (ds *dockerService) ImageFsInfo(_ context.Context, _ *runtimeapi.ImageFsInfoRequest) (*runtimeapi.ImageFsInfoResponse, error) {
	info, err := ds.client.Info()
	if err != nil {
		klog.Errorf("Failed to get docker info: %v", err)
		return nil, err
	}

	statsClient := &winstats.StatsClient{}
	fsinfo, err := statsClient.GetDirFsInfo(info.DockerRootDir)
	if err != nil {
		klog.Errorf("Failed to get dir fsInfo for %q: %v", info.DockerRootDir, err)
		return nil, err
	}

	filesystems := []*runtimeapi.FilesystemUsage{
		{
			Timestamp: time.Now().UnixNano(),
			UsedBytes: &runtimeapi.UInt64Value{Value: fsinfo.Usage},
			FsId: &runtimeapi.FilesystemIdentifier{
				Mountpoint: info.DockerRootDir,
			},
		},
	}

	return &runtimeapi.ImageFsInfoResponse{ImageFilesystems: filesystems}, nil
}
