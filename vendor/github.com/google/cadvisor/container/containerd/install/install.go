// The install package registers containerd.NewPlugin() as the "containerd" container provider when imported
package install

import (
	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/container/containerd"
	"k8s.io/klog"
)

func init() {
	err := container.RegisterPlugin("containerd", containerd.NewPlugin())
	if err != nil {
		klog.Fatalf("Failed to register containerd plugin: %v", err)
	}
}
