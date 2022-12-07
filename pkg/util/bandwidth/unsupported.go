// +build !linux

package bandwidth

import (
	"errors"

	"k8s.io/apimachinery/pkg/api/resource"
)

type unsupportedShaper struct {
}

// NewTCShaper makes a new unsupportedShapper for the given interface
func NewTCShaper(iface string) Shaper {
	return &unsupportedShaper{}
}

func (f *unsupportedShaper) Limit(cidr string, egress, ingress *resource.Quantity) error {
	return errors.New("unimplemented")
}

func (f *unsupportedShaper) Reset(cidr string) error {
	return nil
}

func (f *unsupportedShaper) ReconcileInterface() error {
	return errors.New("unimplemented")
}

func (f *unsupportedShaper) ReconcileCIDR(cidr string, egress, ingress *resource.Quantity) error {
	return errors.New("unimplemented")
}

func (f *unsupportedShaper) GetCIDRs() ([]string, error) {
	return []string{}, nil
}
