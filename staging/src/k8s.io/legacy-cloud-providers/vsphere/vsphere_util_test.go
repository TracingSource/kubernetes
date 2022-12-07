// +build !providerless

package vsphere

import (
	"context"
	"fmt"
	"testing"

	"k8s.io/legacy-cloud-providers/vsphere/vclib"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/simulator"
)

func TestGetPathFromFileNotFound(t *testing.T) {
	ctx := context.Background()

	// vCenter model + initial set of objects (cluster, hosts, VMs, network, datastore, etc)
	model := simulator.VPX()

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

	vc := &vclib.VSphereConnection{Client: c.Client}

	dc, err := vclib.GetDatacenter(ctx, vc, vclib.TestDefaultDatacenter)
	if err != nil {
		t.Errorf("failed to get datacenter: %v", err)
	}

	requestDiskPath := fmt.Sprintf("[%s] %s", vclib.TestDefaultDatastore, DummyDiskName)
	_, err = dc.GetVirtualDiskPage83Data(ctx, requestDiskPath)
	if err == nil {
		t.Error("expected error when calling GetVirtualDiskPage83Data")
	}

	_, err = getPathFromFileNotFound(err)
	if err != nil {
		t.Errorf("expected err to be nil but was %v", err)
	}

	_, err = getPathFromFileNotFound(nil)
	if err == nil {
		t.Errorf("expected err when calling getPathFromFileNotFound with nil err")
	}
}
