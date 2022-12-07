// +build !providerless

package gce

import (
	compute "google.golang.org/api/compute/v1"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
)

func newURLMapMetricContext(request string) *metricContext {
	return newGenericMetricContext("urlmap", request, unusedMetricLabel, unusedMetricLabel, computeV1Version)
}

// GetURLMap returns the UrlMap by name.
func (g *Cloud) GetURLMap(name string) (*compute.UrlMap, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newURLMapMetricContext("get")
	v, err := g.c.UrlMaps().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// CreateURLMap creates a url map
func (g *Cloud) CreateURLMap(urlMap *compute.UrlMap) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newURLMapMetricContext("create")
	return mc.Observe(g.c.UrlMaps().Insert(ctx, meta.GlobalKey(urlMap.Name), urlMap))
}

// UpdateURLMap applies the given UrlMap as an update
func (g *Cloud) UpdateURLMap(urlMap *compute.UrlMap) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newURLMapMetricContext("update")
	return mc.Observe(g.c.UrlMaps().Update(ctx, meta.GlobalKey(urlMap.Name), urlMap))
}

// DeleteURLMap deletes a url map by name.
func (g *Cloud) DeleteURLMap(name string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newURLMapMetricContext("delete")
	return mc.Observe(g.c.UrlMaps().Delete(ctx, meta.GlobalKey(name)))
}

// ListURLMaps lists all UrlMaps in the project.
func (g *Cloud) ListURLMaps() ([]*compute.UrlMap, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newURLMapMetricContext("list")
	v, err := g.c.UrlMaps().List(ctx, filter.None)
	return v, mc.Observe(err)
}
