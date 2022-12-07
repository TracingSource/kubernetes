package bandwidth

import (
	"errors"

	"k8s.io/apimachinery/pkg/api/resource"
)

// FakeShaper provides an implementation of the bandwith.Shaper.
// Beware this is implementation has no features besides Reset and GetCIDRs.
type FakeShaper struct {
	CIDRs      []string
	ResetCIDRs []string
}

// Limit is not implemented
func (f *FakeShaper) Limit(cidr string, egress, ingress *resource.Quantity) error {
	return errors.New("unimplemented")
}

// Reset appends a particular CIDR to the set of ResetCIDRs being managed by this shaper
func (f *FakeShaper) Reset(cidr string) error {
	f.ResetCIDRs = append(f.ResetCIDRs, cidr)
	return nil
}

// ReconcileInterface is not implemented
func (f *FakeShaper) ReconcileInterface() error {
	return errors.New("unimplemented")
}

// ReconcileCIDR is not implemented
func (f *FakeShaper) ReconcileCIDR(cidr string, egress, ingress *resource.Quantity) error {
	return errors.New("unimplemented")
}

// GetCIDRs returns the set of CIDRs that are being managed by this shaper
func (f *FakeShaper) GetCIDRs() ([]string, error) {
	return f.CIDRs, nil
}
