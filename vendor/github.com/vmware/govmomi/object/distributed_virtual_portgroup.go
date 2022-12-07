package object

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type DistributedVirtualPortgroup struct {
	Common
}

func NewDistributedVirtualPortgroup(c *vim25.Client, ref types.ManagedObjectReference) *DistributedVirtualPortgroup {
	return &DistributedVirtualPortgroup{
		Common: NewCommon(c, ref),
	}
}

// EthernetCardBackingInfo returns the VirtualDeviceBackingInfo for this DistributedVirtualPortgroup
func (p DistributedVirtualPortgroup) EthernetCardBackingInfo(ctx context.Context) (types.BaseVirtualDeviceBackingInfo, error) {
	var dvp mo.DistributedVirtualPortgroup
	var dvs mo.DistributedVirtualSwitch
	prop := "config.distributedVirtualSwitch"

	if err := p.Properties(ctx, p.Reference(), []string{"key", prop}, &dvp); err != nil {
		return nil, err
	}

	// "This property should always be set unless the user's setting does not have System.Read privilege on the object referred to by this property."
	if dvp.Config.DistributedVirtualSwitch == nil {
		return nil, fmt.Errorf("no System.Read privilege on: %s.%s", p.Reference(), prop)
	}

	if err := p.Properties(ctx, *dvp.Config.DistributedVirtualSwitch, []string{"uuid"}, &dvs); err != nil {
		return nil, err
	}

	backing := &types.VirtualEthernetCardDistributedVirtualPortBackingInfo{
		Port: types.DistributedVirtualSwitchPortConnection{
			PortgroupKey: dvp.Key,
			SwitchUuid:   dvs.Uuid,
		},
	}

	return backing, nil
}

func (p DistributedVirtualPortgroup) Reconfigure(ctx context.Context, spec types.DVPortgroupConfigSpec) (*Task, error) {
	req := types.ReconfigureDVPortgroup_Task{
		This: p.Reference(),
		Spec: spec,
	}

	res, err := methods.ReconfigureDVPortgroup_Task(ctx, p.Client(), &req)
	if err != nil {
		return nil, err
	}

	return NewTask(p.Client(), res.Returnval), nil
}
