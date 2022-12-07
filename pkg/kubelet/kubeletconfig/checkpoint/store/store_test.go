package store

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/checkpoint"
)

func TestReset(t *testing.T) {
	source, _, err := checkpoint.NewRemoteConfigSource(&apiv1.NodeConfigSource{ConfigMap: &apiv1.ConfigMapNodeConfigSource{
		Name:             "name",
		Namespace:        "namespace",
		UID:              "uid",
		KubeletConfigKey: "kubelet",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	otherSource, _, err := checkpoint.NewRemoteConfigSource(&apiv1.NodeConfigSource{ConfigMap: &apiv1.ConfigMapNodeConfigSource{
		Name:             "other-name",
		Namespace:        "namespace",
		UID:              "other-uid",
		KubeletConfigKey: "kubelet",
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := []struct {
		s       *fakeStore
		updated bool
	}{
		{&fakeStore{assigned: nil, lastKnownGood: nil}, false},
		{&fakeStore{assigned: source, lastKnownGood: nil}, true},
		{&fakeStore{assigned: nil, lastKnownGood: source}, false},
		{&fakeStore{assigned: source, lastKnownGood: source}, true},
		{&fakeStore{assigned: source, lastKnownGood: otherSource}, true},
		{&fakeStore{assigned: otherSource, lastKnownGood: source}, true},
	}
	for _, c := range cases {
		updated, err := reset(c.s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.s.assigned != nil || c.s.lastKnownGood != nil {
			t.Errorf("case %q, expect nil for assigned and last-known-good checkpoints, but still have %q and %q, respectively",
				spew.Sdump(c.s), c.s.assigned, c.s.lastKnownGood)
		}
		if c.updated != updated {
			t.Errorf("case %q, expect reset to return %t, but got %t", spew.Sdump(c.s), c.updated, updated)
		}
	}
}
