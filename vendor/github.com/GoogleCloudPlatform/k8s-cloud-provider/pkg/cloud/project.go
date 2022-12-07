package cloud

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
)

// ProjectRouter routes service calls to the appropriate GCE project.
type ProjectRouter interface {
	// ProjectID returns the project ID (non-numeric) to be used for a call
	// to an API (version,service). Example tuples: ("ga", "ForwardingRules"),
	// ("alpha", "GlobalAddresses").
	//
	// This allows for plumbing different service calls to the appropriate
	// project, for instance, networking services to a separate project
	// than instance management.
	ProjectID(ctx context.Context, version meta.Version, service string) string
}

// SingleProjectRouter routes all service calls to the same project ID.
type SingleProjectRouter struct {
	ID string
}

// ProjectID returns the project ID to be used for a call to the API.
func (r *SingleProjectRouter) ProjectID(ctx context.Context, version meta.Version, service string) string {
	return r.ID
}
