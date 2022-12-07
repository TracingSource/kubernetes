package testsuites

import (
	"fmt"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/storage/testpatterns"
)

// GetDriverNameWithFeatureTags returns driver name with feature tags
// For example)
//  - [Driver: nfs]
//  - [Driver: rbd][Feature:Volumes]
func GetDriverNameWithFeatureTags(driver TestDriver) string {
	dInfo := driver.GetDriverInfo()

	return fmt.Sprintf("[Driver: %s]%s", dInfo.Name, dInfo.FeatureTag)
}

// CreateVolume creates volume for test unless dynamicPV test
func CreateVolume(driver TestDriver, config *PerTestConfig, volType testpatterns.TestVolType) TestVolume {
	switch volType {
	case testpatterns.InlineVolume:
		fallthrough
	case testpatterns.PreprovisionedPV:
		if pDriver, ok := driver.(PreprovisionedVolumeTestDriver); ok {
			return pDriver.CreateVolume(config, volType)
		}
	case testpatterns.DynamicPV:
		// No need to create volume
	default:
		framework.Failf("Invalid volType specified: %v", volType)
	}
	return nil
}

// GetStorageClass constructs a new StorageClass instance
// with a unique name that is based on namespace + suffix.
func GetStorageClass(
	provisioner string,
	parameters map[string]string,
	bindingMode *storagev1.VolumeBindingMode,
	ns string,
	suffix string,
) *storagev1.StorageClass {
	if bindingMode == nil {
		defaultBindingMode := storagev1.VolumeBindingImmediate
		bindingMode = &defaultBindingMode
	}
	return &storagev1.StorageClass{
		TypeMeta: metav1.TypeMeta{
			Kind: "StorageClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			// Name must be unique, so let's base it on namespace name and use GenerateName
			Name: names.SimpleNameGenerator.GenerateName(ns + "-" + suffix),
		},
		Provisioner:       provisioner,
		Parameters:        parameters,
		VolumeBindingMode: bindingMode,
	}
}

// GetSnapshotClass constructs a new SnapshotClass instance
// with a unique name that is based on namespace + suffix.
func GetSnapshotClass(
	snapshotter string,
	parameters map[string]string,
	ns string,
	suffix string,
) *unstructured.Unstructured {
	snapshotClass := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "VolumeSnapshotClass",
			"apiVersion": snapshotAPIVersion,
			"metadata": map[string]interface{}{
				// Name must be unique, so let's base it on namespace name
				"name": ns + "-" + suffix,
			},
			"driver":         snapshotter,
			"parameters":     parameters,
			"deletionPolicy": "Delete",
		},
	}

	return snapshotClass
}
