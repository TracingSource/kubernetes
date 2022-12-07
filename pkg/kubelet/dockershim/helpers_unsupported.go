// +build !linux,!windows

package dockershim

import (
	"fmt"

	"github.com/blang/semver"
	dockertypes "github.com/docker/docker/api/types"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/klog"
)

func DefaultMemorySwap() int64 {
	return -1
}

func (ds *dockerService) getSecurityOpts(seccompProfile string, separator rune) ([]string, error) {
	klog.Warningf("getSecurityOpts is unsupported in this build")
	return nil, nil
}

func (ds *dockerService) updateCreateConfig(
	createConfig *dockertypes.ContainerCreateConfig,
	config *runtimeapi.ContainerConfig,
	sandboxConfig *runtimeapi.PodSandboxConfig,
	podSandboxID string, securityOptSep rune, apiVersion *semver.Version) error {
	klog.Warningf("updateCreateConfig is unsupported in this build")
	return nil
}

func (ds *dockerService) determinePodIPBySandboxID(uid string) []string {
	klog.Warningf("determinePodIPBySandboxID is unsupported in this build")
	return nil
}

func getNetworkNamespace(c *dockertypes.ContainerJSON) (string, error) {
	return "", fmt.Errorf("unsupported platform")
}

// applyExperimentalCreateConfig applies experimental configures from sandbox annotations.
func applyExperimentalCreateConfig(createConfig *dockertypes.ContainerCreateConfig, annotations map[string]string) {
}
