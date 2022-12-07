package simulator

import (
	"time"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type Datastore struct {
	mo.Datastore
}

func parseDatastorePath(dsPath string) (*object.DatastorePath, types.BaseMethodFault) {
	var p object.DatastorePath

	if p.FromString(dsPath) {
		return &p, nil
	}

	return nil, &types.InvalidDatastorePath{DatastorePath: dsPath}
}

func (ds *Datastore) RefreshDatastore(*types.RefreshDatastore) soap.HasFault {
	r := &methods.RefreshDatastoreBody{}

	err := ds.stat()
	if err != nil {
		r.Fault_ = Fault(err.Error(), &types.HostConfigFault{})
		return r
	}

	info := ds.Info.GetDatastoreInfo()

	now := time.Now()

	info.Timestamp = &now

	return r
}
