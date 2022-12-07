// +build !providerless

package gce

import (
	compute "google.golang.org/api/compute/v1"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
)

func newTargetPoolMetricContext(request, region string) *metricContext {
	return newGenericMetricContext("targetpool", request, region, unusedMetricLabel, computeV1Version)
}

// GetTargetPool returns the TargetPool by name.
func (g *Cloud) GetTargetPool(name, region string) (*compute.TargetPool, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newTargetPoolMetricContext("get", region)
	v, err := g.c.TargetPools().Get(ctx, meta.RegionalKey(name, region))
	return v, mc.Observe(err)
}

// CreateTargetPool creates the passed TargetPool
func (g *Cloud) CreateTargetPool(tp *compute.TargetPool, region string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newTargetPoolMetricContext("create", region)
	return mc.Observe(g.c.TargetPools().Insert(ctx, meta.RegionalKey(tp.Name, region), tp))
}

// DeleteTargetPool deletes TargetPool by name.
func (g *Cloud) DeleteTargetPool(name, region string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newTargetPoolMetricContext("delete", region)
	return mc.Observe(g.c.TargetPools().Delete(ctx, meta.RegionalKey(name, region)))
}

// AddInstancesToTargetPool adds instances by link to the TargetPool
func (g *Cloud) AddInstancesToTargetPool(name, region string, instanceRefs []*compute.InstanceReference) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	req := &compute.TargetPoolsAddInstanceRequest{
		Instances: instanceRefs,
	}
	mc := newTargetPoolMetricContext("add_instances", region)
	return mc.Observe(g.c.TargetPools().AddInstance(ctx, meta.RegionalKey(name, region), req))
}

// RemoveInstancesFromTargetPool removes instances by link to the TargetPool
func (g *Cloud) RemoveInstancesFromTargetPool(name, region string, instanceRefs []*compute.InstanceReference) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	req := &compute.TargetPoolsRemoveInstanceRequest{
		Instances: instanceRefs,
	}
	mc := newTargetPoolMetricContext("remove_instances", region)
	return mc.Observe(g.c.TargetPools().RemoveInstance(ctx, meta.RegionalKey(name, region), req))
}
