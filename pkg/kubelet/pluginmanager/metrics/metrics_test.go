package metrics

import (
	"fmt"
	"testing"

	"k8s.io/kubernetes/pkg/kubelet/pluginmanager/cache"
)

func TestMetricCollection(t *testing.T) {
	dsw := cache.NewDesiredStateOfWorld()
	asw := cache.NewActualStateOfWorld()
	fakePlugin := cache.PluginInfo{
		SocketPath: fmt.Sprintf("fake/path/plugin.sock"),
	}
	// Add one plugin to DesiredStateOfWorld
	err := dsw.AddOrUpdatePlugin(fakePlugin.SocketPath)
	if err != nil {
		t.Fatalf("AddOrUpdatePlugin failed. Expected: <no error> Actual: <%v>", err)
	}

	// Add one plugin to ActualStateOfWorld
	err = asw.AddPlugin(fakePlugin)
	if err != nil {
		t.Fatalf("AddOrUpdatePlugin failed. Expected: <no error> Actual: <%v>", err)
	}

	metricCollector := &totalPluginsCollector{asw: asw, dsw: dsw}

	// Check if getPluginCount returns correct data
	count := metricCollector.getPluginCount()
	if len(count) != 2 {
		t.Errorf("getPluginCount failed. Expected <2> states, got <%d>", len(count))
	}

	dswCount, ok := count["desired_state_of_world"]
	if !ok {
		t.Errorf("getPluginCount failed. Expected <desired_state_of_world>, got nothing")
	}

	fakePluginCount := dswCount["fake/path/plugin.sock"]
	if fakePluginCount != 1 {
		t.Errorf("getPluginCount failed. Expected <1> fake/path/plugin.sock in DesiredStateOfWorld, got <%d>",
			fakePluginCount)
	}

	aswCount, ok := count["actual_state_of_world"]
	if !ok {
		t.Errorf("getPluginCount failed. Expected <actual_state_of_world>, got nothing")
	}

	fakePluginCount = aswCount["fake/path/plugin.sock"]
	if fakePluginCount != 1 {
		t.Errorf("getPluginCount failed. Expected <1> fake/path/plugin.sock in ActualStateOfWorld, got <%d>",
			fakePluginCount)
	}
}
