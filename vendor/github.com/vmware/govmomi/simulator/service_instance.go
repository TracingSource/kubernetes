package simulator

import (
	"time"

	"github.com/google/uuid"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/simulator/vpx"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type ServiceInstance struct {
	mo.ServiceInstance
}

func NewServiceInstance(content types.ServiceContent, folder mo.Folder) *ServiceInstance {
	Map = NewRegistry()

	s := &ServiceInstance{}

	s.Self = vim25.ServiceInstance
	s.Content = content

	Map.Put(s)

	f := &Folder{Folder: folder}
	Map.Put(f)

	var setting []types.BaseOptionValue

	if content.About.ApiType == "HostAgent" {
		CreateDefaultESX(f)
	} else {
		content.About.InstanceUuid = uuid.New().String()
		setting = vpx.Setting
	}

	objects := []object.Reference{
		NewSessionManager(*s.Content.SessionManager),
		NewAuthorizationManager(*s.Content.AuthorizationManager),
		NewPerformanceManager(*s.Content.PerfManager),
		NewPropertyCollector(s.Content.PropertyCollector),
		NewFileManager(*s.Content.FileManager),
		NewVirtualDiskManager(*s.Content.VirtualDiskManager),
		NewLicenseManager(*s.Content.LicenseManager),
		NewSearchIndex(*s.Content.SearchIndex),
		NewViewManager(*s.Content.ViewManager),
		NewEventManager(*s.Content.EventManager),
		NewTaskManager(*s.Content.TaskManager),
		NewUserDirectory(*s.Content.UserDirectory),
		NewOptionManager(s.Content.Setting, setting),
		NewStorageResourceManager(*s.Content.StorageResourceManager),
	}

	switch content.VStorageObjectManager.Type {
	case "HostVStorageObjectManager":
		// TODO: NewHostVStorageObjectManager(*content.VStorageObjectManager)
	case "VcenterVStorageObjectManager":
		objects = append(objects, NewVcenterVStorageObjectManager(*content.VStorageObjectManager))
	}

	if s.Content.CustomFieldsManager != nil {
		objects = append(objects, NewCustomFieldsManager(*s.Content.CustomFieldsManager))
	}

	if s.Content.IpPoolManager != nil {
		objects = append(objects, NewIpPoolManager(*s.Content.IpPoolManager))
	}

	if s.Content.AccountManager != nil {
		objects = append(objects, NewHostLocalAccountManager(*s.Content.AccountManager))
	}

	for _, o := range objects {
		Map.Put(o)
	}

	return s
}

func (s *ServiceInstance) RetrieveServiceContent(*types.RetrieveServiceContent) soap.HasFault {
	return &methods.RetrieveServiceContentBody{
		Res: &types.RetrieveServiceContentResponse{
			Returnval: s.Content,
		},
	}
}

func (*ServiceInstance) CurrentTime(*types.CurrentTime) soap.HasFault {
	return &methods.CurrentTimeBody{
		Res: &types.CurrentTimeResponse{
			Returnval: time.Now(),
		},
	}
}
