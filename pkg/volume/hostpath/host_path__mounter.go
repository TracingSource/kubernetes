package hostpath

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/mount"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util/hostutil"
	"k8s.io/kubernetes/pkg/volume/validation"
)

type hostPathMounter struct {
	*hostPath
	readOnly bool
	mounter  mount.Interface
	hu       hostutil.HostUtils
}

var _ volume.Mounter = &hostPathMounter{}

func (b *hostPathMounter) GetAttributes() volume.Attributes {
	return volume.Attributes{
		ReadOnly:        b.readOnly,
		Managed:         false,
		SupportsSELinux: false,
	}
}

// Checks prior to mount operations to verify that the required components (binaries, etc.)
// to mount the volume are available on the underlying node.
// If not, it returns an error
func (b *hostPathMounter) CanMount() error {
	return nil
}

// SetUp does nothing.
func (b *hostPathMounter) SetUp(mounterArgs volume.MounterArgs) error {
	err := validation.ValidatePathNoBacksteps(b.GetPath())
	if err != nil {
		return fmt.Errorf("invalid HostPath `%s`: %v", b.GetPath(), err)
	}

	if *b.pathType == v1.HostPathUnset {
		return nil
	}
	return checkType(b.GetPath(), b.pathType, b.hu)
}

// SetUpAt does not make sense for host paths - probably programmer error.
func (b *hostPathMounter) SetUpAt(dir string, mounterArgs volume.MounterArgs) error {
	return fmt.Errorf("SetUpAt() does not make sense for host paths")
}

func (b *hostPathMounter) GetPath() string {
	return b.path
}

// checkType checks whether the given path is the exact pathType
func checkType(path string, pathType *v1.HostPathType, hu hostutil.HostUtils) error {
	return checkTypeInternal(newFileTypeChecker(path, hu), pathType)
}

func checkTypeInternal(ftc hostPathTypeChecker, pathType *v1.HostPathType) error {
	switch *pathType {
	case v1.HostPathDirectoryOrCreate:
		if !ftc.Exists() {
			return ftc.MakeDir()
		}
		fallthrough
	case v1.HostPathDirectory:
		if !ftc.IsDir() {
			return fmt.Errorf("hostPath type check failed: %s is not a directory", ftc.GetPath())
		}
	case v1.HostPathFileOrCreate:
		if !ftc.Exists() {
			return ftc.MakeFile()
		}
		fallthrough
	case v1.HostPathFile:
		if !ftc.IsFile() {
			return fmt.Errorf("hostPath type check failed: %s is not a file", ftc.GetPath())
		}
	case v1.HostPathSocket:
		if !ftc.IsSocket() {
			return fmt.Errorf("hostPath type check failed: %s is not a socket file", ftc.GetPath())
		}
	case v1.HostPathCharDev:
		if !ftc.IsChar() {
			return fmt.Errorf("hostPath type check failed: %s is not a character device", ftc.GetPath())
		}
	case v1.HostPathBlockDev:
		if !ftc.IsBlock() {
			return fmt.Errorf("hostPath type check failed: %s is not a block device", ftc.GetPath())
		}
	default:
		return fmt.Errorf("%s is an invalid volume type", *pathType)
	}

	return nil
}
