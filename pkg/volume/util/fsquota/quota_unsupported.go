// +build !linux

package fsquota

import (
	"errors"

	"k8s.io/utils/mount"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

// Dummy quota implementation for systems that do not implement support
// for volume quotas

var errNotImplemented = errors.New("not implemented")

// SupportsQuotas -- dummy implementation
func SupportsQuotas(_ mount.Interface, _ string) (bool, error) {
	return false, errNotImplemented
}

// AssignQuota -- dummy implementation
func AssignQuota(_ mount.Interface, _ string, _ types.UID, _ *resource.Quantity) error {
	return errNotImplemented
}

// GetConsumption -- dummy implementation
func GetConsumption(_ string) (*resource.Quantity, error) {
	return nil, errNotImplemented
}

// GetInodes -- dummy implementation
func GetInodes(_ string) (*resource.Quantity, error) {
	return nil, errNotImplemented
}

// ClearQuota -- dummy implementation
func ClearQuota(_ mount.Interface, _ string) error {
	return errNotImplemented
}
