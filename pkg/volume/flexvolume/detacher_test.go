package flexvolume

import (
	"testing"

	"k8s.io/kubernetes/test/utils/harness"
)

func TestDetach(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	plugin, _ := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), detachCmd,
			"sdx", "localhost"),
	)

	d, _ := plugin.NewDetacher()
	d.Detach("sdx", "localhost")
}

func TestUnmountDevice(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	plugin, rootDir := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), unmountDeviceCmd,
			rootDir+"/mount-dir"),
	)

	d, _ := plugin.NewDetacher()
	d.UnmountDevice(rootDir + "/mount-dir")
}
