// Copyright 2017 Google Inc. All Rights Reserved.
//

package accelerators

import info "github.com/google/cadvisor/info/v1"

// This is supposed to store global state about an accelerator metrics collector.
// cadvisor manager will call Setup() when it starts and Destroy() when it stops.
// For each container detected by the cadvisor manager, it will call
// GetCollector() with the devices cgroup path for that container.
// GetCollector() is supposed to return an object that can update
// accelerator stats for that container.
type AcceleratorManager interface {
	Setup()
	Destroy()
	GetCollector(deviceCgroup string) (AcceleratorCollector, error)
}

type AcceleratorCollector interface {
	UpdateStats(*info.ContainerStats) error
}
