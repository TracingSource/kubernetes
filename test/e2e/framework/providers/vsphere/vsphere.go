package vsphere

import (
	"k8s.io/kubernetes/test/e2e/framework"
)

func init() {
	framework.RegisterProvider("vsphere", newProvider)
}

func newProvider() (framework.ProviderInterface, error) {
	return &framework.NullProvider{}, nil
}
