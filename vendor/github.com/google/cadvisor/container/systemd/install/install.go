// The install package registers systemd.NewPlugin() as the "systemd" container provider when imported
package install

import (
	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/container/systemd"
	"k8s.io/klog"
)

func init() {
	err := container.RegisterPlugin("systemd", systemd.NewPlugin())
	if err != nil {
		klog.Fatalf("Failed to register systemd plugin: %v", err)
	}
}
