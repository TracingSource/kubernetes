package flexvolume

import (
	"testing"

	"k8s.io/kubernetes/test/utils/harness"
	"k8s.io/utils/exec/testing"
)

func TestInit(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	plugin, _ := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, successOutput(), "init"),
	)
	plugin.Init(plugin.host)
}

func fakeVolumeNameOutput(name string) testingexec.FakeCombinedOutputAction {
	return fakeResultOutput(&DriverStatus{
		Status:     StatusSuccess,
		VolumeName: name,
	})
}

func TestGetVolumeName(tt *testing.T) {
	t := harness.For(tt)
	defer t.Close()

	spec := fakeVolumeSpec()
	plugin, _ := testPlugin(t)
	plugin.runner = fakeRunner(
		assertDriverCall(t, fakeVolumeNameOutput(spec.Name()), getVolumeNameCmd,
			specJSON(plugin, spec, nil)),
	)

	name, err := plugin.GetVolumeName(spec)
	if err != nil {
		t.Errorf("GetVolumeName() failed: %v", err)
	}
	expectedName := spec.Name()
	if name != expectedName {
		t.Errorf("GetVolumeName() returned %v instead of %v", name, expectedName)
	}
}
