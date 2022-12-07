package flexvolume

import (
	"testing"

	"k8s.io/utils/mount"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/test/utils/harness"
)

func TestTearDownAt(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	mounter := mount.NewFakeMounter(nil)

	plugin, rootDir := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, notSupportedOutput(), unmountCmd,
			rootDir+"/mount-dir"),
	)

	u, _ := plugin.newUnmounterInternal("volName", types.UID("poduid"), mounter, plugin.runner)
	u.TearDownAt(rootDir + "/mount-dir")
}
