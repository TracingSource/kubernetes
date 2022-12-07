package flexvolume

import (
	"testing"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/test/utils/harness"
)

func TestAttach(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	spec := fakeVolumeSpec()

	plugin, _ := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), attachCmd,
			specJSON(plugin, spec, nil), "localhost"),
	)

	a, _ := plugin.NewAttacher()
	a.Attach(spec, "localhost")
}

func TestWaitForAttach(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	spec := fakeVolumeSpec()
	var pod *v1.Pod
	plugin, _ := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), waitForAttachCmd, "/dev/sdx",
			specJSON(plugin, spec, nil)),
	)

	a, _ := plugin.NewAttacher()
	a.WaitForAttach(spec, "/dev/sdx", pod, 1*time.Second)
}

func TestMountDevice(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	spec := fakeVolumeSpec()

	plugin, rootDir := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), mountDeviceCmd, rootDir+"/mount-dir", "/dev/sdx",
			specJSON(plugin, spec, nil)),
	)

	a, _ := plugin.NewAttacher()
	a.MountDevice(spec, "/dev/sdx", rootDir+"/mount-dir")
}

func TestIsVolumeAttached(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	spec := fakeVolumeSpec()

	plugin, _ := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), isAttached, specJSON(plugin, spec, nil), "localhost"),
	)
	a, _ := plugin.NewAttacher()
	specs := []*volume.Spec{spec}
	a.VolumesAreAttached(specs, "localhost")
}
