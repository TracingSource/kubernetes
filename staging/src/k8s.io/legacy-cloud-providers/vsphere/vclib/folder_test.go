package vclib

import (
	"context"
	"path"
	"testing"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/simulator"
)

func TestFolder(t *testing.T) {
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

	const folderName = "F0"
	vmFolder := path.Join("/", folderName, dc.Name(), "vm")

	tests := []struct {
		folderPath string
		expect     int
	}{
		{vmFolder, 0},
		{path.Join(vmFolder, folderName), (model.Host + model.Cluster) * model.Machine},
	}

	for i, test := range tests {
		folder, cerr := dc.GetFolderByPath(ctx, test.folderPath)
		if cerr != nil {
			t.Fatal(cerr)
		}

		vms, cerr := folder.GetVirtualMachines(ctx)
		if cerr != nil {
			t.Fatalf("%d: %s", i, cerr)
		}

		if len(vms) != test.expect {
			t.Errorf("%d: expected %d VMs, got: %d", i, test.expect, len(vms))
		}
	}
}
