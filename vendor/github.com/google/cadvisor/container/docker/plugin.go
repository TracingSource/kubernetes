package docker

import (
	"time"

	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/fs"
	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/watcher"
	"golang.org/x/net/context"
	"k8s.io/klog"
)

const dockerClientTimeout = 10 * time.Second

// NewPlugin returns an implementation of container.Plugin suitable for passing to container.RegisterPlugin()
func NewPlugin() container.Plugin {
	return &plugin{}
}

type plugin struct{}

func (p *plugin) InitializeFSContext(context *fs.Context) error {
	SetTimeout(dockerClientTimeout)
	// Try to connect to docker indefinitely on startup.
	dockerStatus := retryDockerStatus()
	context.Docker = fs.DockerContext{
		Root:         RootDir(),
		Driver:       dockerStatus.Driver,
		DriverStatus: dockerStatus.DriverStatus,
	}
	return nil
}

func (p *plugin) Register(factory info.MachineInfoFactory, fsInfo fs.FsInfo, includedMetrics container.MetricSet) (watcher.ContainerWatcher, error) {
	err := Register(factory, fsInfo, includedMetrics)
	return nil, err
}

func retryDockerStatus() info.DockerStatus {
	startupTimeout := dockerClientTimeout
	maxTimeout := 4 * startupTimeout
	for {
		ctx, _ := context.WithTimeout(context.Background(), startupTimeout)
		dockerStatus, err := StatusWithContext(ctx)
		if err == nil {
			return dockerStatus
		}

		switch err {
		case context.DeadlineExceeded:
			klog.Warningf("Timeout trying to communicate with docker during initialization, will retry")
		default:
			klog.V(5).Infof("Docker not connected: %v", err)
			return info.DockerStatus{}
		}

		startupTimeout = 2 * startupTimeout
		if startupTimeout > maxTimeout {
			startupTimeout = maxTimeout
		}
	}
}
