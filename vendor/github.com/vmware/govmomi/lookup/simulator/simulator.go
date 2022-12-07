package simulator

import (
	"github.com/vmware/govmomi/lookup"
	"github.com/vmware/govmomi/lookup/methods"
	"github.com/vmware/govmomi/lookup/types"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25/soap"
	vim "github.com/vmware/govmomi/vim25/types"
)

var content = types.LookupServiceContent{
	LookupService:                vim.ManagedObjectReference{Type: "LookupLookupService", Value: "lookupService"},
	ServiceRegistration:          &vim.ManagedObjectReference{Type: "LookupServiceRegistration", Value: "ServiceRegistration"},
	DeploymentInformationService: vim.ManagedObjectReference{Type: "LookupDeploymentInformationService", Value: "deploymentInformationService"},
	L10n:                         vim.ManagedObjectReference{Type: "LookupL10n", Value: "l10n"},
}

func New() *simulator.Registry {
	r := simulator.NewRegistry()
	r.Namespace = lookup.Namespace
	r.Path = lookup.Path

	r.Put(&ServiceInstance{
		ManagedObjectReference: lookup.ServiceInstance,
		Content:                content,
	})
	r.Put(&ServiceRegistration{
		ManagedObjectReference: *content.ServiceRegistration,
		Info:                   registrationInfo(),
	})

	return r
}

type ServiceInstance struct {
	vim.ManagedObjectReference

	Content types.LookupServiceContent
}

func (s *ServiceInstance) RetrieveServiceContent(_ *types.RetrieveServiceContent) soap.HasFault {
	return &methods.RetrieveServiceContentBody{
		Res: &types.RetrieveServiceContentResponse{
			Returnval: s.Content,
		},
	}
}

type ServiceRegistration struct {
	vim.ManagedObjectReference

	Info []types.LookupServiceRegistrationInfo
}

func (s *ServiceRegistration) GetSiteId(_ *types.GetSiteId) soap.HasFault {
	return &methods.GetSiteIdBody{
		Res: &types.GetSiteIdResponse{
			Returnval: siteID,
		},
	}
}

func matchServiceType(filter, info *types.LookupServiceRegistrationServiceType) bool {
	if filter.Product != "" {
		if filter.Product != info.Product {
			return false
		}
	}

	if filter.Type != "" {
		if filter.Type != info.Type {
			return false
		}
	}

	return true
}

func matchEndpointType(filter, info *types.LookupServiceRegistrationEndpointType) bool {
	if filter.Protocol != "" {
		if filter.Protocol != info.Protocol {
			return false
		}
	}

	if filter.Type != "" {
		if filter.Type != info.Type {
			return false
		}
	}

	return true
}

func (s *ServiceRegistration) List(req *types.List) soap.HasFault {
	body := new(methods.ListBody)
	filter := req.FilterCriteria

	if filter == nil {
		// This is what a real PSC returns if FilterCriteria is nil.
		body.Fault_ = simulator.Fault("LookupFaultServiceFault", &vim.SystemError{
			Reason: "Invalid fault",
		})
		return body
	}
	body.Res = new(types.ListResponse)

	for _, info := range s.Info {
		if filter.SiteId != "" {
			if filter.SiteId != info.SiteId {
				continue
			}
		}
		if filter.NodeId != "" {
			if filter.NodeId != info.NodeId {
				continue
			}
		}
		if filter.ServiceType != nil {
			if !matchServiceType(filter.ServiceType, &info.ServiceType) {
				continue
			}
		}
		if filter.EndpointType != nil {
			services := info.ServiceEndpoints
			info.ServiceEndpoints = nil
			for _, service := range services {
				if !matchEndpointType(filter.EndpointType, &service.EndpointType) {
					continue
				}
				info.ServiceEndpoints = append(info.ServiceEndpoints, service)
			}
			if len(info.ServiceEndpoints) == 0 {
				continue
			}
		}
		body.Res.Returnval = append(body.Res.Returnval, info)
	}

	return body
}
