package capabilities

import (
	"sync"
)

// Capabilities defines the set of capabilities available within the system.
// For now these are global.  Eventually they may be per-user
type Capabilities struct {
	AllowPrivileged bool

	// Pod sources from which to allow privileged capabilities like host networking, sharing the host
	// IPC namespace, and sharing the host PID namespace.
	PrivilegedSources PrivilegedSources

	// PerConnectionBandwidthLimitBytesPerSec limits the throughput of each connection (currently only used for proxy, exec, attach)
	PerConnectionBandwidthLimitBytesPerSec int64
}

// PrivilegedSources defines the pod sources allowed to make privileged requests for certain types
// of capabilities like host networking, sharing the host IPC namespace, and sharing the host PID namespace.
type PrivilegedSources struct {
	// List of pod sources for which using host network is allowed.
	HostNetworkSources []string

	// List of pod sources for which using host pid namespace is allowed.
	HostPIDSources []string

	// List of pod sources for which using host ipc is allowed.
	HostIPCSources []string
}

var capInstance struct {
	once         sync.Once
	lock         sync.Mutex
	capabilities *Capabilities
}

// caller: 
// 	1. cmd/kubelet/app/server.go -> RunKubelet()
//
// Initialize the capability set. 
// This can only be done once per binary, subsequent calls are ignored.
func Initialize(c Capabilities) {
	// Only do this once
	capInstance.once.Do(func() {
		capInstance.capabilities = &c
	})
}

// Setup the capability set.  It wraps Initialize for improving usability.
func Setup(allowPrivileged bool, perConnectionBytesPerSec int64) {
	Initialize(Capabilities{
		AllowPrivileged:                        allowPrivileged,
		PerConnectionBandwidthLimitBytesPerSec: perConnectionBytesPerSec,
	})
}

// SetForTests sets capabilities for tests. Convenience method for testing. 
// This should only be called from tests.
func SetForTests(c Capabilities) {
	capInstance.lock.Lock()
	defer capInstance.lock.Unlock()
	capInstance.capabilities = &c
}

// Get returns a read-only copy of the system capabilities.
func Get() Capabilities {
	capInstance.lock.Lock()
	defer capInstance.lock.Unlock()
	// This check prevents clobbering of capabilities that might've been set via SetForTests
	if capInstance.capabilities == nil {
		Initialize(Capabilities{
			AllowPrivileged: false,
			PrivilegedSources: PrivilegedSources{
				HostNetworkSources: []string{},
				HostPIDSources:     []string{},
				HostIPCSources:     []string{},
			},
		})
	}
	return *capInstance.capabilities
}
