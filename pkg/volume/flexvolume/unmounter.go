package flexvolume

import (
	"fmt"
	"os"

	"k8s.io/klog"
	"k8s.io/utils/exec"
	"k8s.io/utils/mount"

	"k8s.io/kubernetes/pkg/volume"
)

// FlexVolumeUnmounter is the disk that will be cleaned by this plugin.
type flexVolumeUnmounter struct {
	*flexVolume
	// Runner used to teardown the volume.
	runner exec.Interface
}

var _ volume.Unmounter = &flexVolumeUnmounter{}

// Unmounter interface
func (f *flexVolumeUnmounter) TearDown() error {
	path := f.GetPath()
	return f.TearDownAt(path)
}

func (f *flexVolumeUnmounter) TearDownAt(dir string) error {
	pathExists, pathErr := mount.PathExists(dir)
	if pathErr != nil {
		// only log warning here since plugins should anyways have to deal with errors
		klog.Warningf("Error checking path: %v", pathErr)
	} else {
		if !pathExists {
			klog.Warningf("Warning: Unmount skipped because path does not exist: %v", dir)
			return nil
		}
	}

	call := f.plugin.NewDriverCall(unmountCmd)
	call.Append(dir)
	_, err := call.Run()
	if isCmdNotSupportedErr(err) {
		err = (*unmounterDefaults)(f).TearDownAt(dir)
	}
	if err != nil {
		return err
	}

	// Flexvolume driver may remove the directory. Ignore if it does.
	if pathExists, pathErr := mount.PathExists(dir); pathErr != nil {
		return fmt.Errorf("Error checking if path exists: %v", pathErr)
	} else if !pathExists {
		return nil
	}
	return os.Remove(dir)
}
