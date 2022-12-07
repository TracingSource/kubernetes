package object

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type OpaqueNetwork struct {
	Common
}

func NewOpaqueNetwork(c *vim25.Client, ref types.ManagedObjectReference) *OpaqueNetwork {
	return &OpaqueNetwork{
		Common: NewCommon(c, ref),
	}
}

// EthernetCardBackingInfo returns the VirtualDeviceBackingInfo for this Network
func (n OpaqueNetwork) EthernetCardBackingInfo(ctx context.Context) (types.BaseVirtualDeviceBackingInfo, error) {
	var net mo.OpaqueNetwork

	if err := n.Properties(ctx, n.Reference(), []string{"summary"}, &net); err != nil {
		return nil, err
	}

	summary, ok := net.Summary.(*types.OpaqueNetworkSummary)
	if !ok {
		return nil, fmt.Errorf("%s unsupported network type: %T", n, net.Summary)
	}

	backing := &types.VirtualEthernetCardOpaqueNetworkBackingInfo{
		OpaqueNetworkId:   summary.OpaqueNetworkId,
		OpaqueNetworkType: summary.OpaqueNetworkType,
	}

	return backing, nil
}
