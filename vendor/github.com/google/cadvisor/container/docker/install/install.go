// Copyright 2019 Google Inc. All Rights Reserved.
//


// The install package registers docker.NewPlugin() as the "docker" container provider when imported
package install

import (
	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/container/docker"
	"k8s.io/klog"
)

func init() {
	err := container.RegisterPlugin("docker", docker.NewPlugin())
	if err != nil {
		klog.Fatalf("Failed to register docker plugin: %v", err)
	}
}
