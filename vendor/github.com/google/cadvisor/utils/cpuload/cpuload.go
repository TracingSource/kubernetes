package cpuload

import (
	"fmt"

	info "github.com/google/cadvisor/info/v1"

	"github.com/google/cadvisor/utils/cpuload/netlink"
	"k8s.io/klog"
)

type CpuLoadReader interface {
	// Start the reader.
	Start() error

	// Stop the reader and clean up internal state.
	Stop()

	// Retrieve Cpu load for a given group.
	// name is the full hierarchical name of the container.
	// Path is an absolute filesystem path for a container under CPU cgroup hierarchy.
	GetCpuLoad(name string, path string) (info.LoadStats, error)
}

func New() (CpuLoadReader, error) {
	reader, err := netlink.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create a netlink based cpuload reader: %v", err)
	}
	klog.V(4).Info("Using a netlink-based load reader")
	return reader, nil
}
