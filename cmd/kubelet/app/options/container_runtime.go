package options

import (
	"runtime"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/kubelet/config"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	// When these values are updated, also update test/e2e/framework/util.go
	defaultPodSandboxImageName    = "k8s.gcr.io/pause"
	defaultPodSandboxImageVersion = "3.1"
)

var (
	defaultPodSandboxImage = defaultPodSandboxImageName +
		":" + defaultPodSandboxImageVersion
)

// caller: 
// 	1. cmd/kubelet/app/options/options.go -> NewKubeletFlags()
//
// NewContainerRuntimeOptions will create a new ContainerRuntimeOptions with
// default values.
func NewContainerRuntimeOptions() *config.ContainerRuntimeOptions {
	dockerEndpoint := ""
	if runtime.GOOS != "windows" {
		dockerEndpoint = "unix:///var/run/docker.sock"
	}

	return &config.ContainerRuntimeOptions{
		ContainerRuntime:           kubetypes.DockerContainerRuntime,
		RedirectContainerStreaming: false,
		DockerEndpoint:             dockerEndpoint,
		DockershimRootDirectory:    "/var/lib/dockershim",
		PodSandboxImage:            defaultPodSandboxImage,
		ImagePullProgressDeadline:  metav1.Duration{Duration: 1 * time.Minute},
		ExperimentalDockershim:     false,

		//Alpha feature
		CNIBinDir:   "/opt/cni/bin",
		CNIConfDir:  "/etc/cni/net.d",
		CNICacheDir: "/var/lib/cni/cache",
	}
}
