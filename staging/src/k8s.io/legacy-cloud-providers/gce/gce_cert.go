// +build !providerless

package gce

import (
	compute "google.golang.org/api/compute/v1"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
)

func newCertMetricContext(request string) *metricContext {
	return newGenericMetricContext("cert", request, unusedMetricLabel, unusedMetricLabel, computeV1Version)
}

// GetSslCertificate returns the SslCertificate by name.
func (g *Cloud) GetSslCertificate(name string) (*compute.SslCertificate, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newCertMetricContext("get")
	v, err := g.c.SslCertificates().Get(ctx, meta.GlobalKey(name))
	return v, mc.Observe(err)
}

// CreateSslCertificate creates and returns a SslCertificate.
func (g *Cloud) CreateSslCertificate(sslCerts *compute.SslCertificate) (*compute.SslCertificate, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newCertMetricContext("create")
	err := g.c.SslCertificates().Insert(ctx, meta.GlobalKey(sslCerts.Name), sslCerts)
	if err != nil {
		return nil, mc.Observe(err)
	}
	return g.GetSslCertificate(sslCerts.Name)
}

// DeleteSslCertificate deletes the SslCertificate by name.
func (g *Cloud) DeleteSslCertificate(name string) error {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newCertMetricContext("delete")
	return mc.Observe(g.c.SslCertificates().Delete(ctx, meta.GlobalKey(name)))
}

// ListSslCertificates lists all SslCertificates in the project.
func (g *Cloud) ListSslCertificates() ([]*compute.SslCertificate, error) {
	ctx, cancel := cloud.ContextWithCallTimeout()
	defer cancel()

	mc := newCertMetricContext("list")
	v, err := g.c.SslCertificates().List(ctx, filter.None)
	return v, mc.Observe(err)
}
