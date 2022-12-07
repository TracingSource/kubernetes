package crio

import (
	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/fs"
	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/watcher"
	"k8s.io/klog"
)

// NewPlugin returns an implementation of container.Plugin suitable for passing to container.RegisterPlugin()
func NewPlugin() container.Plugin {
	return &plugin{}
}

type plugin struct{}

func (p *plugin) InitializeFSContext(context *fs.Context) error {
	crioClient, err := Client()
	if err != nil {
		return err
	}

	crioInfo, err := crioClient.Info()
	if err != nil {
		klog.V(5).Infof("CRI-O not connected: %v", err)
	} else {
		context.Crio = fs.CrioContext{Root: crioInfo.StorageRoot}
	}
	return nil
}

func (p *plugin) Register(factory info.MachineInfoFactory, fsInfo fs.FsInfo, includedMetrics container.MetricSet) (watcher.ContainerWatcher, error) {
	err := Register(factory, fsInfo, includedMetrics)
	return nil, err
}
