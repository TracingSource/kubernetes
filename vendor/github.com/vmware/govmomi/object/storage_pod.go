package object

import (
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

type StoragePod struct {
	*Folder
}

func NewStoragePod(c *vim25.Client, ref types.ManagedObjectReference) *StoragePod {
	return &StoragePod{
		Folder: &Folder{
			Common: NewCommon(c, ref),
		},
	}
}
