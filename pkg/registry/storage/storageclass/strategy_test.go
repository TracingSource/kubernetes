package storageclass

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/storage"
)

func TestStorageClassStrategy(t *testing.T) {
	ctx := genericapirequest.NewDefaultContext()
	if Strategy.NamespaceScoped() {
		t.Errorf("StorageClass must not be namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("StorageClass should not allow create on update")
	}

	deleteReclaimPolicy := api.PersistentVolumeReclaimDelete
	bindingMode := storage.VolumeBindingWaitForFirstConsumer
	storageClass := &storage.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "valid-class",
		},
		Provisioner: "kubernetes.io/aws-ebs",
		Parameters: map[string]string{
			"foo": "bar",
		},
		ReclaimPolicy:     &deleteReclaimPolicy,
		VolumeBindingMode: &bindingMode,
	}

	Strategy.PrepareForCreate(ctx, storageClass)

	errs := Strategy.Validate(ctx, storageClass)
	if len(errs) != 0 {
		t.Errorf("unexpected error validating %v", errs)
	}

	newStorageClass := &storage.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "valid-class-2",
			ResourceVersion: "4",
		},
		Provisioner: "kubernetes.io/aws-ebs",
		Parameters: map[string]string{
			"foo": "bar",
		},
		ReclaimPolicy:     &deleteReclaimPolicy,
		VolumeBindingMode: &bindingMode,
	}

	Strategy.PrepareForUpdate(ctx, newStorageClass, storageClass)

	errs = Strategy.ValidateUpdate(ctx, newStorageClass, storageClass)
	if len(errs) == 0 {
		t.Errorf("Expected a validation error")
	}
}
