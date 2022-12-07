package multipoint

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// CommunicatingPlugin is an example of a plugin that implements two
// extension points. It communicates through state with another function.
type CommunicatingPlugin struct{}

var _ framework.ReservePlugin = CommunicatingPlugin{}
var _ framework.PreBindPlugin = CommunicatingPlugin{}

// Name is the name of the plugin used in Registry and configurations.
const Name = "multipoint-communicating-plugin"

// Name returns name of the plugin. It is used in logs, etc.
func (mc CommunicatingPlugin) Name() string {
	return Name
}

type stateData struct {
	data string
}

func (s *stateData) Clone() framework.StateData {
	copy := &stateData{
		data: s.data,
	}
	return copy
}

// Reserve is the functions invoked by the framework at "reserve" extension point.
func (mc CommunicatingPlugin) Reserve(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	if pod == nil {
		return framework.NewStatus(framework.Error, "pod cannot be nil")
	}
	if pod.Name == "my-test-pod" {
		state.Lock()
		state.Write(framework.StateKey(pod.Name), &stateData{data: "never bind"})
		state.Unlock()
	}
	return nil
}

// PreBind is the functions invoked by the framework at "prebind" extension point.
func (mc CommunicatingPlugin) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	if pod == nil {
		return framework.NewStatus(framework.Error, "pod cannot be nil")
	}
	state.RLock()
	defer state.RUnlock()
	if v, e := state.Read(framework.StateKey(pod.Name)); e == nil {
		if value, ok := v.(*stateData); ok && value.data == "never bind" {
			return framework.NewStatus(framework.Unschedulable, "pod is not permitted")
		}
	}
	return nil
}

// New initializes a new plugin and returns it.
func New(_ *runtime.Unknown, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &CommunicatingPlugin{}, nil
}
