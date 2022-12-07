package vclib

import "errors"

// Error Messages
const (
	FileAlreadyExistErrMsg     = "File requested already exist"
	NoDiskUUIDFoundErrMsg      = "No disk UUID found"
	NoDevicesFoundErrMsg       = "No devices found"
	DiskNotFoundErrMsg         = "No vSphere disk ID found"
	InvalidVolumeOptionsErrMsg = "VolumeOptions verification failed"
	NoVMFoundErrMsg            = "No VM found"
)

// Error constants
var (
	ErrFileAlreadyExist     = errors.New(FileAlreadyExistErrMsg)
	ErrNoDiskUUIDFound      = errors.New(NoDiskUUIDFoundErrMsg)
	ErrNoDevicesFound       = errors.New(NoDevicesFoundErrMsg)
	ErrNoDiskIDFound        = errors.New(DiskNotFoundErrMsg)
	ErrInvalidVolumeOptions = errors.New(InvalidVolumeOptionsErrMsg)
	ErrNoVMFound            = errors.New(NoVMFoundErrMsg)
)
