package vclib

// Volume Constants
const (
	ThinDiskType             = "thin"
	PreallocatedDiskType     = "preallocated"
	EagerZeroedThickDiskType = "eagerZeroedThick"
	ZeroedThickDiskType      = "zeroedThick"
)

// Controller Constants
const (
	SCSIControllerLimit       = 4
	SCSIControllerDeviceLimit = 15
	SCSIDeviceSlots           = 16
	SCSIReservedSlot          = 7

	SCSIControllerType        = "scsi"
	LSILogicControllerType    = "lsiLogic"
	BusLogicControllerType    = "busLogic"
	LSILogicSASControllerType = "lsiLogic-sas"
	PVSCSIControllerType      = "pvscsi"
)

// Other Constants
const (
	LogLevel                   = 4
	DatastoreProperty          = "datastore"
	ResourcePoolProperty       = "resourcePool"
	DatastoreInfoProperty      = "info"
	VirtualMachineType         = "VirtualMachine"
	RoundTripperDefaultCount   = 3
	VSANDatastoreType          = "vsan"
	DummyVMPrefixName          = "vsphere-k8s"
	ActivePowerState           = "poweredOn"
	DatacenterType             = "Datacenter"
	ClusterComputeResourceType = "ClusterComputeResource"
	HostSystemType             = "HostSystem"
	NameProperty               = "name"
)

// Test Constants
const (
	TestDefaultDatacenter = "DC0"
	TestDefaultDatastore  = "LocalDS_0"
	TestDefaultNetwork    = "VM Network"
	testNameNotFound      = "enoent"
)
