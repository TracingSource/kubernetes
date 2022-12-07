package rbac

import (
	"reflect"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

// IsOnlyMutatingGCFields checks finalizers and ownerrefs which GC manipulates
// and indicates that only those fields are changing
func IsOnlyMutatingGCFields(obj, old runtime.Object, equalities conversion.Equalities) bool {
	if old == nil || reflect.ValueOf(old).IsNil() {
		return false
	}

	// make a copy of the newObj so that we can stomp for comparison
	copied := obj.DeepCopyObject()
	copiedMeta, err := meta.Accessor(copied)
	if err != nil {
		return false
	}
	oldMeta, err := meta.Accessor(old)
	if err != nil {
		return false
	}
	copiedMeta.SetOwnerReferences(oldMeta.GetOwnerReferences())
	copiedMeta.SetFinalizers(oldMeta.GetFinalizers())
	copiedMeta.SetSelfLink(oldMeta.GetSelfLink())

	return equalities.DeepEqual(copied, old)
}
