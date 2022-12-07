package diskmanagers

import (
	"context"
	"time"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"

	"k8s.io/klog"
	"k8s.io/legacy-cloud-providers/vsphere/vclib"
)

// virtualDiskManager implements VirtualDiskProvider Interface for creating and deleting volume using VirtualDiskManager
type virtualDiskManager struct {
	diskPath      string
	volumeOptions *vclib.VolumeOptions
}

// Create implements Disk's Create interface
// Contains implementation of virtualDiskManager based Provisioning
func (diskManager virtualDiskManager) Create(ctx context.Context, datastore *vclib.Datastore) (canonicalDiskPath string, err error) {
	if diskManager.volumeOptions.SCSIControllerType == "" {
		diskManager.volumeOptions.SCSIControllerType = vclib.LSILogicControllerType
	}

	// Check for existing VMDK before attempting create. Because a name collision
	// is unlikely, "VMDK already exists" is likely from a previous attempt to
	// create this volume.
	if dsPath := vclib.GetPathFromVMDiskPath(diskManager.diskPath); datastore.Exists(ctx, dsPath) {
		klog.V(2).Infof("Create: VirtualDisk already exists, returning success. Name=%q", diskManager.diskPath)
		return diskManager.diskPath, nil
	}

	// Create specification for new virtual disk
	diskFormat := vclib.DiskFormatValidType[diskManager.volumeOptions.DiskFormat]
	vmDiskSpec := &types.FileBackedVirtualDiskSpec{
		VirtualDiskSpec: types.VirtualDiskSpec{
			AdapterType: diskManager.volumeOptions.SCSIControllerType,
			DiskType:    diskFormat,
		},
		CapacityKb: int64(diskManager.volumeOptions.CapacityKB),
	}

	vdm := object.NewVirtualDiskManager(datastore.Client())
	requestTime := time.Now()
	// Create virtual disk
	task, err := vdm.CreateVirtualDisk(ctx, diskManager.diskPath, datastore.Datacenter.Datacenter, vmDiskSpec)
	if err != nil {
		vclib.RecordvSphereMetric(vclib.APICreateVolume, requestTime, err)
		klog.Errorf("Failed to create virtual disk: %s. err: %+v", diskManager.diskPath, err)
		return "", err
	}
	taskInfo, err := task.WaitForResult(ctx, nil)
	vclib.RecordvSphereMetric(vclib.APICreateVolume, requestTime, err)
	if err != nil {
		klog.Errorf("Failed to complete virtual disk creation: %s. err: %+v", diskManager.diskPath, err)
		return "", err
	}
	canonicalDiskPath = taskInfo.Result.(string)
	return canonicalDiskPath, nil
}

// Delete implements Disk's Delete interface
func (diskManager virtualDiskManager) Delete(ctx context.Context, datacenter *vclib.Datacenter) error {
	// Create a virtual disk manager
	virtualDiskManager := object.NewVirtualDiskManager(datacenter.Client())
	diskPath := vclib.RemoveStorageClusterORFolderNameFromVDiskPath(diskManager.diskPath)
	requestTime := time.Now()
	// Delete virtual disk
	task, err := virtualDiskManager.DeleteVirtualDisk(ctx, diskPath, datacenter.Datacenter)
	if err != nil {
		klog.Errorf("Failed to delete virtual disk. err: %v", err)
		vclib.RecordvSphereMetric(vclib.APIDeleteVolume, requestTime, err)
		return err
	}
	err = task.Wait(ctx)
	vclib.RecordvSphereMetric(vclib.APIDeleteVolume, requestTime, err)
	if err != nil {
		klog.Errorf("Failed to delete virtual disk. err: %v", err)
		return err
	}
	return nil
}
