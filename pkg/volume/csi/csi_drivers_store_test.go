package csi_test

import (
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/volume/csi"
)

func TestDriversStore(t *testing.T) {
	store := &csi.DriversStore{}
	someDriver := csi.Driver{}

	expectAbsent(t, store, "does-not-exist")

	store.Set("some-driver", someDriver)
	expectPresent(t, store, "some-driver", someDriver)

	store.Delete("some-driver")
	expectAbsent(t, store, "some-driver")

	store.Set("some-driver", someDriver)

	store.Clear()
	expectAbsent(t, store, "some-driver")
}

func expectPresent(t *testing.T, store *csi.DriversStore, name string, expected csi.Driver) {
	t.Helper()

	retrieved, ok := store.Get(name)

	if !ok {
		t.Fatalf("expected driver '%s' to exist", name)
	}

	if !reflect.DeepEqual(retrieved, expected) {
		t.Fatalf("expected driver '%s' to be equal to %v", name, expected)
	}
}

func expectAbsent(t *testing.T, store *csi.DriversStore, name string) {
	t.Helper()

	if _, ok := store.Get(name); ok {
		t.Fatalf("expected driver '%s' not to exist in store", name)
	}
}
