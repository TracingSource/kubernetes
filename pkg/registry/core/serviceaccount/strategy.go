package serviceaccount

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/validation"
)

// strategy implements behavior for ServiceAccount objects
type strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating ServiceAccount
// objects via the REST API.
var Strategy = strategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (strategy) NamespaceScoped() bool {
	return true
}

func (strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	cleanSecretReferences(obj.(*api.ServiceAccount))
}

func (strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return validation.ValidateServiceAccount(obj.(*api.ServiceAccount))
}

// Canonicalize normalizes the object after validation.
func (strategy) Canonicalize(obj runtime.Object) {
}

func (strategy) AllowCreateOnUpdate() bool {
	return false
}

func (strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	cleanSecretReferences(obj.(*api.ServiceAccount))
}

func cleanSecretReferences(serviceAccount *api.ServiceAccount) {
	for i, secret := range serviceAccount.Secrets {
		serviceAccount.Secrets[i] = api.ObjectReference{Name: secret.Name}
	}
}

func (strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateServiceAccountUpdate(obj.(*api.ServiceAccount), old.(*api.ServiceAccount))
}

func (strategy) AllowUnconditionalUpdate() bool {
	return true
}
