// +build !linux,!darwin,!windows

package fs

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

// FSInfo unsupported returns 0 values for available and capacity and an error.
func FsInfo(path string) (int64, int64, int64, int64, int64, int64, error) {
	return 0, 0, 0, 0, 0, 0, fmt.Errorf("FsInfo not supported for this build.")
}

// DiskUsage gets disk usage of specified path.
func DiskUsage(path string) (*resource.Quantity, error) {
	return nil, fmt.Errorf("Du not supported for this build.")
}

func Find(path string) (int64, error) {
	return 0, fmt.Errorf("Find not supported for this build.")
}
