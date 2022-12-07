// +build linux

package cadvisor

import (
	"fmt"

	cadvisorfs "github.com/google/cadvisor/fs"
	"k8s.io/kubernetes/pkg/kubelet/types"
)

// imageFsInfoProvider knows how to translate the configured runtime
// to its file system label for images.
type imageFsInfoProvider struct {
	runtime         string
	runtimeEndpoint string
}

// ImageFsInfoLabel returns the image fs label for the configured runtime.
// For remote runtimes, it handles additional runtimes natively understood by cAdvisor.
func (i *imageFsInfoProvider) ImageFsInfoLabel() (string, error) {
	switch i.runtime {
	case types.DockerContainerRuntime:
		return cadvisorfs.LabelDockerImages, nil
	case types.RemoteContainerRuntime:
		// This is a temporary workaround to get stats for cri-o from cadvisor
		// and should be removed.
		// Related to https://github.com/kubernetes/kubernetes/issues/51798
		if i.runtimeEndpoint == CrioSocket || i.runtimeEndpoint == "unix://"+CrioSocket {
			return cadvisorfs.LabelCrioImages, nil
		}
	}
	return "", fmt.Errorf("no imagefs label for configured runtime")
}

// NewImageFsInfoProvider 直接返回 imageFsInfoProvider{} 结构体
//
// caller: 
// 	1. cmd/kubelet/app/server.go -> run()
//
// NewImageFsInfoProvider returns a provider for the specified runtime configuration.
func NewImageFsInfoProvider(runtime, runtimeEndpoint string) ImageFsInfoProvider {
	return &imageFsInfoProvider{runtime: runtime, runtimeEndpoint: runtimeEndpoint}
}
