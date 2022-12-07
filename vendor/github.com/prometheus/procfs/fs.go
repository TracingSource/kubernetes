package procfs

import (
	"github.com/prometheus/procfs/internal/fs"
)

// FS represents the pseudo-filesystem sys, which provides an interface to
// kernel data structures.
type FS struct {
	proc fs.FS
}

// DefaultMountPoint is the common mount point of the proc filesystem.
const DefaultMountPoint = fs.DefaultProcMountPoint

// NewDefaultFS returns a new proc FS mounted under the default proc mountPoint.
// It will error if the mount point directory can't be read or is a file.
func NewDefaultFS() (FS, error) {
	return NewFS(DefaultMountPoint)
}

// NewFS returns a new proc FS mounted under the given proc mountPoint. It will error
// if the mount point directory can't be read or is a file.
func NewFS(mountPoint string) (FS, error) {
	fs, err := fs.NewFS(mountPoint)
	if err != nil {
		return FS{}, err
	}
	return FS{fs}, nil
}
