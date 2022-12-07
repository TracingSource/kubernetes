package storageclass

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/storage"
	storageutil "k8s.io/kubernetes/pkg/apis/storage/util"
	"k8s.io/kubernetes/pkg/apis/storage/validation"
)

// storageClassStrategy implements behavior for StorageClass objects
type storageClassStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// StorageClass objects via the REST API.
var Strategy = storageClassStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (storageClassStrategy) NamespaceScoped() bool {
	return false
}

// ResetBeforeCreate clears the Status field which is not allowed to be set by end users on creation.
func (storageClassStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	class := obj.(*storage.StorageClass)

	storageutil.DropDisabledFields(class, nil)
}

func (storageClassStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	storageClass := obj.(*storage.StorageClass)
	return validation.ValidateStorageClass(storageClass)
}

// Canonicalize normalizes the object after validation.
func (storageClassStrategy) Canonicalize(obj runtime.Object) {
}

func (storageClassStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the Status fields which is not allowed to be set by an end user updating a PV
func (storageClassStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newClass := obj.(*storage.StorageClass)
	oldClass := old.(*storage.StorageClass)

	storageutil.DropDisabledFields(oldClass, newClass)
}

func (storageClassStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	errorList := validation.ValidateStorageClass(obj.(*storage.StorageClass))
	return append(errorList, validation.ValidateStorageClassUpdate(obj.(*storage.StorageClass), old.(*storage.StorageClass))...)
}

func (storageClassStrategy) AllowUnconditionalUpdate() bool {
	return true
}
