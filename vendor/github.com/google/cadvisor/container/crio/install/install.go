// The install package registers crio.NewPlugin() as the "crio" container provider when imported
package install

import (
	"github.com/google/cadvisor/container"
	"github.com/google/cadvisor/container/crio"
	"k8s.io/klog"
)

func init() {
	err := container.RegisterPlugin("crio", crio.NewPlugin())
	if err != nil {
		klog.Fatalf("Failed to register crio plugin: %v", err)
	}
}
