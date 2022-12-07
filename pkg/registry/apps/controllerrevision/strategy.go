package controllerrevision

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/apps/validation"
)

// strategy implements behavior for ConfigMap objects
type strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating ControllerRevision
// objects via the REST API.
var Strategy = strategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// Strategy should implement rest.RESTCreateStrategy
var _ rest.RESTCreateStrategy = Strategy

// Strategy should implement rest.RESTUpdateStrategy
var _ rest.RESTUpdateStrategy = Strategy

func (strategy) NamespaceScoped() bool {
	return true
}

func (strategy) Canonicalize(obj runtime.Object) {
}

func (strategy) AllowCreateOnUpdate() bool {
	return false
}

func (strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	_ = obj.(*apps.ControllerRevision)
}

func (strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	revision := obj.(*apps.ControllerRevision)

	return validation.ValidateControllerRevision(revision)
}

func (strategy) PrepareForUpdate(ctx context.Context, newObj, oldObj runtime.Object) {
	_ = oldObj.(*apps.ControllerRevision)
	_ = newObj.(*apps.ControllerRevision)
}

func (strategy) AllowUnconditionalUpdate() bool {
	return true
}

func (strategy) ValidateUpdate(ctx context.Context, newObj, oldObj runtime.Object) field.ErrorList {
	oldRevision, newRevision := oldObj.(*apps.ControllerRevision), newObj.(*apps.ControllerRevision)
	return validation.ValidateControllerRevisionUpdate(newRevision, oldRevision)
}
