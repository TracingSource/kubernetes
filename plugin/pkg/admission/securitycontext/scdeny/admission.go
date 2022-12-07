package scdeny

import (
	"context"
	"fmt"
	"io"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// PluginName indicates name of admission plugin.
const PluginName = "SecurityContextDeny"

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewSecurityContextDeny(), nil
	})
}

// Plugin implements admission.Interface.
type Plugin struct {
	*admission.Handler
}

var _ admission.ValidationInterface = &Plugin{}

// NewSecurityContextDeny creates a new instance of the SecurityContextDeny admission controller
func NewSecurityContextDeny() *Plugin {
	return &Plugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

// Validate will deny any pod that defines SupplementalGroups, SELinuxOptions, RunAsUser or FSGroup
func (p *Plugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) (err error) {
	if a.GetSubresource() != "" || a.GetResource().GroupResource() != api.Resource("pods") {
		return nil
	}

	pod, ok := a.GetObject().(*api.Pod)
	if !ok {
		return apierrors.NewBadRequest("Resource was marked with kind Pod but was unable to be converted")
	}

	if pod.Spec.SecurityContext != nil {
		if pod.Spec.SecurityContext.SupplementalGroups != nil {
			return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("pod.Spec.SecurityContext.SupplementalGroups is forbidden"))
		}
		if pod.Spec.SecurityContext.SELinuxOptions != nil {
			return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("pod.Spec.SecurityContext.SELinuxOptions is forbidden"))
		}
		if pod.Spec.SecurityContext.RunAsUser != nil {
			return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("pod.Spec.SecurityContext.RunAsUser is forbidden"))
		}
		if pod.Spec.SecurityContext.FSGroup != nil {
			return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("pod.Spec.SecurityContext.FSGroup is forbidden"))
		}
	}

	for _, v := range pod.Spec.InitContainers {
		if v.SecurityContext != nil {
			if v.SecurityContext.SELinuxOptions != nil {
				return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("SecurityContext.SELinuxOptions is forbidden"))
			}
			if v.SecurityContext.RunAsUser != nil {
				return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("SecurityContext.RunAsUser is forbidden"))
			}
		}
	}

	for _, v := range pod.Spec.Containers {
		if v.SecurityContext != nil {
			if v.SecurityContext.SELinuxOptions != nil {
				return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("SecurityContext.SELinuxOptions is forbidden"))
			}
			if v.SecurityContext.RunAsUser != nil {
				return apierrors.NewForbidden(a.GetResource().GroupResource(), pod.Name, fmt.Errorf("SecurityContext.RunAsUser is forbidden"))
			}
		}
	}
	return nil
}
