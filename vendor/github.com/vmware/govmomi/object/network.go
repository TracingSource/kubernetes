package object

import (
	"context"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type Network struct {
	Common
}

func NewNetwork(c *vim25.Client, ref types.ManagedObjectReference) *Network {
	return &Network{
		Common: NewCommon(c, ref),
	}
}

// EthernetCardBackingInfo returns the VirtualDeviceBackingInfo for this Network
func (n Network) EthernetCardBackingInfo(ctx context.Context) (types.BaseVirtualDeviceBackingInfo, error) {
	var e mo.Network

	// Use Network.Name rather than Common.Name as the latter does not return the complete name if it contains a '/'
	// We can't use Common.ObjectName here either as we need the ManagedEntity.Name field is not set since mo.Network
	// has its own Name field.
	err := n.Properties(ctx, n.Reference(), []string{"name"}, &e)
	if err != nil {
		return nil, err
	}

	backing := &types.VirtualEthernetCardNetworkBackingInfo{
		VirtualDeviceDeviceBackingInfo: types.VirtualDeviceDeviceBackingInfo{
			DeviceName: e.Name,
		},
	}

	return backing, nil
}
