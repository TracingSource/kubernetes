package devicemapper

import (
	"fmt"
	"os"
	"path/filepath"
)

// ThinLsBinaryPresent returns the location of the thin_ls binary in the mount
// namespace cadvisor is running in or an error.  The locations checked are:
//
// - /sbin/
// - /bin/
// - /usr/sbin/
// - /usr/bin/
//
// The thin_ls binary is provided by the device-mapper-persistent-data
// package.
func ThinLsBinaryPresent() (string, error) {
	var (
		thinLsPath string
		err        error
	)

	for _, path := range []string{"/sbin", "/bin", "/usr/sbin/", "/usr/bin"} {
		// try paths for non-containerized operation
		// note: thin_ls is most likely a symlink to pdata_tools
		thinLsPath = filepath.Join(path, "thin_ls")
		_, err = os.Stat(thinLsPath)
		if err == nil {
			return thinLsPath, nil
		}
	}

	return "", fmt.Errorf("unable to find thin_ls binary")
}
