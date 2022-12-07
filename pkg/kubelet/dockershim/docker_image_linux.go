// +build linux

package dockershim

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"k8s.io/klog"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ImageFsInfo returns information of the filesystem that is used to store images.
func (ds *dockerService) ImageFsInfo(_ context.Context, _ *runtimeapi.ImageFsInfoRequest) (*runtimeapi.ImageFsInfoResponse, error) {
	info, err := ds.client.Info()
	if err != nil {
		klog.Errorf("Failed to get docker info: %v", err)
		return nil, err
	}

	bytes, inodes, err := dirSize(filepath.Join(info.DockerRootDir, "image"))
	if err != nil {
		return nil, err
	}

	return &runtimeapi.ImageFsInfoResponse{
		ImageFilesystems: []*runtimeapi.FilesystemUsage{
			{
				Timestamp: time.Now().Unix(),
				FsId: &runtimeapi.FilesystemIdentifier{
					Mountpoint: info.DockerRootDir,
				},
				UsedBytes: &runtimeapi.UInt64Value{
					Value: uint64(bytes),
				},
				InodesUsed: &runtimeapi.UInt64Value{
					Value: uint64(inodes),
				},
			},
		},
	}, nil
}

func dirSize(path string) (int64, int64, error) {
	bytes := int64(0)
	inodes := int64(0)
	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		inodes += 1
		if !info.IsDir() {
			bytes += info.Size()
		}
		return nil
	})
	return bytes, inodes, err
}
