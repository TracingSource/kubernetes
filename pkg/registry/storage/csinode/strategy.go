package csinode

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/storage"
	"k8s.io/kubernetes/pkg/apis/storage/validation"
)

// csiNodeStrategy implements behavior for CSINode objects
type csiNodeStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating
// CSINode objects via the REST API.
var Strategy = csiNodeStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (csiNodeStrategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate clears fields that are not allowed to be set on creation.
func (csiNodeStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (csiNodeStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	csiNode := obj.(*storage.CSINode)

	errs := validation.ValidateCSINode(csiNode)
	errs = append(errs, validation.ValidateCSINode(csiNode)...)

	return errs
}

// Canonicalize normalizes the object after validation.
func (csiNodeStrategy) Canonicalize(obj runtime.Object) {
}

func (csiNodeStrategy) AllowCreateOnUpdate() bool {
	return false
}

// PrepareForUpdate sets the driver's Allocatable fields that are not allowed to be set by an end user updating a CSINode.
func (csiNodeStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (csiNodeStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newCSINodeObj := obj.(*storage.CSINode)
	oldCSINodeObj := old.(*storage.CSINode)
	errorList := validation.ValidateCSINode(newCSINodeObj)
	return append(errorList, validation.ValidateCSINodeUpdate(newCSINodeObj, oldCSINodeObj)...)
}

func (csiNodeStrategy) AllowUnconditionalUpdate() bool {
	return false
}
