// +build !providerless

package gce

import (
	compute "google.golang.org/api/compute/v1"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
)

func newFirewallMetricContext(request string) *metricContext {
	return newGenericMetricContext("firewall", request, unusedMetricLabel, unusedMetricLabel, computeV1Version)
}

// GetFirewall returns the Firewall by name.
func (g *Cloud) GetFirewall(name string) (*compute.Firewall, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newFirewallMetricContext("get")
	v, err := g.c.Firewalls().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// CreateFirewall creates the passed firewall
func (g *Cloud) CreateFirewall(f *compute.Firewall) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newFirewallMetricContext("create")
	return mc.Observe(g.c.Firewalls().Insert(ctx, meta.GlobalKey(f.Name), f))
}

// DeleteFirewall deletes the given firewall rule.
func (g *Cloud) DeleteFirewall(name string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newFirewallMetricContext("delete")
	return mc.Observe(g.c.Firewalls().Delete(ctx, meta.GlobalKey(name)))
}

// UpdateFirewall applies the given firewall as an update to an existing service.
func (g *Cloud) UpdateFirewall(f *compute.Firewall) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newFirewallMetricContext("update")
	return mc.Observe(g.c.Firewalls().Update(ctx, meta.GlobalKey(f.Name), f))
}
