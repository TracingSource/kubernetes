package openstack

import (
	"k8s.io/kubernetes/test/e2e/framework"
)

func init() {
	framework.RegisterProvider("openstack", newProvider)
}

func newProvider() (framework.ProviderInterface, error) {
	return &Provider{}, nil
}

// Provider is a structure to handle OpenStack clouds for e2e testing
type Provider struct {
	framework.NullProvider
}
