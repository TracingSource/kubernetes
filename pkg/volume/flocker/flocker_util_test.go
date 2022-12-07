package flocker

import (
	"fmt"
	"os"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/volume"
	volumetest "k8s.io/kubernetes/pkg/volume/testing"

	"github.com/stretchr/testify/assert"
)

func TestFlockerUtil_CreateVolume(t *testing.T) {
	assert := assert.New(t)

	// test CreateVolume happy path
	pvc := volumetest.CreateTestPVC("3Gi", []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce})
	options := volume.VolumeOptions{
		PVC:                           pvc,
		PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimDelete,
	}

	fakeFlockerClient := newFakeFlockerClient()
	dir, p := newTestableProvisioner(assert, options)
	provisioner := p.(*flockerVolumeProvisioner)
	defer os.RemoveAll(dir)
	provisioner.flockerClient = fakeFlockerClient

	flockerUtil := &flockerUtil{}

	datasetID, size, _, err := flockerUtil.CreateVolume(provisioner)
	assert.NoError(err)
	assert.Equal(datasetOneID, datasetID)
	assert.Equal(3, size)

	// test error during CreateVolume
	fakeFlockerClient.Error = fmt.Errorf("Do not feel like provisioning")
	_, _, _, err = flockerUtil.CreateVolume(provisioner)
	assert.Equal(fakeFlockerClient.Error.Error(), err.Error())
}
