package vclib

import (
	"context"
	"testing"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/simulator"
)

func TestUtils(t *testing.T) {
	ctx := context.Background()

	model := simulator.VPX()
	// Child folder "F0" will be created under the root folder and datacenter folders,
	// and all resources are created within the "F0" child folders.
	model.Folder = 1

	defer model.Remove()
	err := model.Create()
	if err != nil {
		t.Fatal(err)
	}

	s := model.Service.NewServer()
	defer s.Close()

	c, err := govmomi.NewClient(ctx, s.URL, true)
	if err != nil {
		t.Fatal(err)
	}

	vc := &VSphereConnection{Client: c.Client}

	dc, err := GetDatacenter(ctx, vc, TestDefaultDatacenter)
	if err != nil {
		t.Error(err)
	}

	finder := getFinder(dc)
	datastores, err := finder.DatastoreList(ctx, "*")
	if err != nil {
		t.Fatal(err)
	}

	count := model.Count()
	if count.Datastore != len(datastores) {
		t.Errorf("got %d Datastores, expected: %d", len(datastores), count.Datastore)
	}

	_, err = finder.Datastore(ctx, testNameNotFound)
	if !IsNotFound(err) {
		t.Errorf("unexpected error: %s", err)
	}
}
