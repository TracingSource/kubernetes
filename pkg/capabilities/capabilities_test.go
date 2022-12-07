package capabilities

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	defaultCap := Capabilities{
		AllowPrivileged: false,
		PrivilegedSources: PrivilegedSources{
			HostNetworkSources: []string{},
			HostPIDSources:     []string{},
			HostIPCSources:     []string{},
		},
	}

	res := Get()
	if !reflect.DeepEqual(defaultCap, res) {
		t.Fatalf("expected Capabilities: %#v, got a non-default: %#v", defaultCap, res)
	}

	cap := Capabilities{
		PrivilegedSources: PrivilegedSources{
			HostNetworkSources: []string{"A", "B"},
		},
	}
	SetForTests(cap)

	res = Get()
	if !reflect.DeepEqual(cap, res) {
		t.Fatalf("expected Capabilities: %#v , got a different: %#v", cap, res)
	}
}
