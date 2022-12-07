// +build windows

package fs

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	modkernel32            = windows.NewLazySystemDLL("kernel32.dll")
	procGetDiskFreeSpaceEx = modkernel32.NewProc("GetDiskFreeSpaceExW")
)

// FSInfo returns (available bytes, byte capacity, byte usage, total inodes, inodes free, inode usage, error)
// for the filesystem that path resides upon.
func FsInfo(path string) (int64, int64, int64, int64, int64, int64, error) {
	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64
	var err error

	ret, _, err := syscall.Syscall6(
		procGetDiskFreeSpaceEx.Addr(),
		4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
		0,
		0,
	)
	if ret == 0 {
		return 0, 0, 0, 0, 0, 0, err
	}

	return freeBytesAvailable, totalNumberOfBytes, totalNumberOfBytes - freeBytesAvailable, 0, 0, 0, nil
}

// DiskUsage gets disk usage of specified path.
func DiskUsage(path string) (*resource.Quantity, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	usage, err := diskUsage(path, info)
	if err != nil {
		return nil, err
	}

	used, err := resource.ParseQuantity(fmt.Sprintf("%d", usage))
	if err != nil {
		return nil, fmt.Errorf("failed to parse fs usage %d due to %v", usage, err)
	}
	used.Format = resource.BinarySI
	return &used, nil
}

// Always return zero since inodes is not supported on Windows.
func Find(path string) (int64, error) {
	return 0, nil
}

func diskUsage(currPath string, info os.FileInfo) (int64, error) {
	var size int64

	if info.Mode()&os.ModeSymlink != 0 {
		return size, nil
	}

	size += info.Size()

	if !info.IsDir() {
		return size, nil
	}

	dir, err := os.Open(currPath)
	if err != nil {
		return size, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return size, err
	}

	for _, file := range files {
		if file.IsDir() {
			s, err := diskUsage(fmt.Sprintf("%s/%s", currPath, file.Name()), file)
			if err != nil {
				return size, err
			}
			size += s
		} else {
			size += file.Size()
		}
	}
	return size, nil
}
