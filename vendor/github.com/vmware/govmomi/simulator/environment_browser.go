package simulator

import (
	"github.com/vmware/govmomi/simulator/esx"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type EnvironmentBrowser struct {
	mo.EnvironmentBrowser
}

func newEnvironmentBrowser() *types.ManagedObjectReference {
	env := new(EnvironmentBrowser)
	Map.Put(env)
	return &env.Self
}

func (b *EnvironmentBrowser) QueryConfigOption(req *types.QueryConfigOption) soap.HasFault {
	body := new(methods.QueryConfigOptionBody)

	opt := &types.VirtualMachineConfigOption{
		Version:       esx.HardwareVersion,
		DefaultDevice: esx.VirtualDevice,
	}

	body.Res = &types.QueryConfigOptionResponse{
		Returnval: opt,
	}

	return body
}

func (b *EnvironmentBrowser) QueryConfigOptionEx(req *types.QueryConfigOptionEx) soap.HasFault {
	body := new(methods.QueryConfigOptionExBody)

	opt := &types.VirtualMachineConfigOption{
		Version:       esx.HardwareVersion,
		DefaultDevice: esx.VirtualDevice,
	}

	body.Res = &types.QueryConfigOptionExResponse{
		Returnval: opt,
	}

	return body
}
