package errors

import (
	k8stypes "k8s.io/apimachinery/pkg/types"
)

// NewDeletedVolumeInUseError returns a new instance of DeletedVolumeInUseError
// error.
func NewDeletedVolumeInUseError(message string) error {
	return deletedVolumeInUseError(message)
}

type deletedVolumeInUseError string

var _ error = deletedVolumeInUseError("")

// IsDeletedVolumeInUse returns true if an error returned from Delete() is
// deletedVolumeInUseError
func IsDeletedVolumeInUse(err error) bool {
	switch err.(type) {
	case deletedVolumeInUseError:
		return true
	default:
		return false
	}
}

func (err deletedVolumeInUseError) Error() string {
	return string(err)
}

// DanglingAttachError indicates volume is attached to a different node
// than we expected.
type DanglingAttachError struct {
	msg         string
	CurrentNode k8stypes.NodeName
	DevicePath  string
}

func (err *DanglingAttachError) Error() string {
	return err.msg
}

// NewDanglingError create a new dangling error
func NewDanglingError(msg string, node k8stypes.NodeName, devicePath string) error {
	return &DanglingAttachError{
		msg:         msg,
		CurrentNode: node,
		DevicePath:  devicePath,
	}
}
