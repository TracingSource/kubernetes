package endpointslice

import (
	"context"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/discovery"
	"k8s.io/kubernetes/pkg/apis/discovery/validation"
)

// endpointSliceStrategy implements verification logic for Replication.
type endpointSliceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Replication EndpointSlice objects.
var Strategy = endpointSliceStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped returns true because all EndpointSlices need to be within a namespace.
func (endpointSliceStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of an EndpointSlice before creation.
func (endpointSliceStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	endpointSlice := obj.(*discovery.EndpointSlice)
	endpointSlice.Generation = 1
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (endpointSliceStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newEPS := obj.(*discovery.EndpointSlice)
	oldEPS := old.(*discovery.EndpointSlice)

	// Increment generation if anything other than meta changed
	// This needs to be changed if a status attribute is added to EndpointSlice
	ogNewMeta := newEPS.ObjectMeta
	ogOldMeta := oldEPS.ObjectMeta
	newEPS.ObjectMeta = v1.ObjectMeta{}
	oldEPS.ObjectMeta = v1.ObjectMeta{}

	if !apiequality.Semantic.DeepEqual(newEPS, oldEPS) {
		ogNewMeta.Generation = ogOldMeta.Generation + 1
	}

	newEPS.ObjectMeta = ogNewMeta
	oldEPS.ObjectMeta = ogOldMeta
}

// Validate validates a new EndpointSlice.
func (endpointSliceStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	endpointSlice := obj.(*discovery.EndpointSlice)
	err := validation.ValidateEndpointSliceCreate(endpointSlice)
	return err
}

// Canonicalize normalizes the object after validation.
func (endpointSliceStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for EndpointSlice; this means POST is needed to create one.
func (endpointSliceStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (endpointSliceStrategy) ValidateUpdate(ctx context.Context, new, old runtime.Object) field.ErrorList {
	newEPS := new.(*discovery.EndpointSlice)
	oldEPS := old.(*discovery.EndpointSlice)
	return validation.ValidateEndpointSliceUpdate(newEPS, oldEPS)
}

// AllowUnconditionalUpdate is the default update policy for EndpointSlice objects.
func (endpointSliceStrategy) AllowUnconditionalUpdate() bool {
	return true
}
