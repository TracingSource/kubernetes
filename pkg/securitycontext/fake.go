package securitycontext

import (
	"k8s.io/api/core/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// ValidSecurityContextWithContainerDefaults creates a valid security context provider based on
// empty container defaults.  Used for testing.
func ValidSecurityContextWithContainerDefaults() *v1.SecurityContext {
	priv := false
	defProcMount := v1.DefaultProcMount
	return &v1.SecurityContext{
		Capabilities: &v1.Capabilities{},
		Privileged:   &priv,
		ProcMount:    &defProcMount,
	}
}

// ValidInternalSecurityContextWithContainerDefaults creates a valid security context provider based on
// empty container defaults.  Used for testing.
func ValidInternalSecurityContextWithContainerDefaults() *api.SecurityContext {
	priv := false
	dpm := api.DefaultProcMount
	return &api.SecurityContext{
		Capabilities: &api.Capabilities{},
		Privileged:   &priv,
		ProcMount:    &dpm,
	}
}
