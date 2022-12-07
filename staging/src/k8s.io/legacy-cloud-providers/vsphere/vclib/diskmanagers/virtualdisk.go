package diskmanagers

import (
	"context"
	"fmt"

	"k8s.io/klog"
	"k8s.io/legacy-cloud-providers/vsphere/vclib"
)

// VirtualDisk is for the Disk Management
type VirtualDisk struct {
	DiskPath      string
	VolumeOptions *vclib.VolumeOptions
	VMOptions     *vclib.VMOptions
}

// VirtualDisk Operations Const
const (
	VirtualDiskCreateOperation = "Create"
	VirtualDiskDeleteOperation = "Delete"
)

// VirtualDiskProvider defines interfaces for creating disk
type VirtualDiskProvider interface {
	Create(ctx context.Context, datastore *vclib.Datastore) (string, error)
	Delete(ctx context.Context, datacenter *vclib.Datacenter) error
}

// getDiskManager returns vmDiskManager or vdmDiskManager based on given volumeoptions
func getDiskManager(disk *VirtualDisk, diskOperation string) VirtualDiskProvider {
	var diskProvider VirtualDiskProvider
	switch diskOperation {
	case VirtualDiskDeleteOperation:
		diskProvider = virtualDiskManager{disk.DiskPath, disk.VolumeOptions}
	case VirtualDiskCreateOperation:
		if disk.VolumeOptions.StoragePolicyName != "" || disk.VolumeOptions.VSANStorageProfileData != "" || disk.VolumeOptions.StoragePolicyID != "" {
			diskProvider = vmDiskManager{disk.DiskPath, disk.VolumeOptions, disk.VMOptions}
		} else {
			diskProvider = virtualDiskManager{disk.DiskPath, disk.VolumeOptions}
		}
	}
	return diskProvider
}

// Create gets appropriate disk manager and calls respective create method
func (virtualDisk *VirtualDisk) Create(ctx context.Context, datastore *vclib.Datastore) (string, error) {
	if virtualDisk.VolumeOptions.DiskFormat == "" {
		virtualDisk.VolumeOptions.DiskFormat = vclib.ThinDiskType
	}
	if !virtualDisk.VolumeOptions.VerifyVolumeOptions() {
		klog.Error("VolumeOptions verification failed. volumeOptions: ", virtualDisk.VolumeOptions)
		return "", vclib.ErrInvalidVolumeOptions
	}
	if virtualDisk.VolumeOptions.StoragePolicyID != "" && virtualDisk.VolumeOptions.StoragePolicyName != "" {
		return "", fmt.Errorf("Storage Policy ID and Storage Policy Name both set, Please set only one parameter")
	}
	return getDiskManager(virtualDisk, VirtualDiskCreateOperation).Create(ctx, datastore)
}

// Delete gets appropriate disk manager and calls respective delete method
func (virtualDisk *VirtualDisk) Delete(ctx context.Context, datacenter *vclib.Datacenter) error {
	return getDiskManager(virtualDisk, VirtualDiskDeleteOperation).Delete(ctx, datacenter)
}
