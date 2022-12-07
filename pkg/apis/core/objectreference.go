//TODO: consider making these methods functions, because we don't want helper
//functions in the k8s.io/api repo.

package core

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SetGroupVersionKind sets the API version and kind of the object reference
func (obj *ObjectReference) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	obj.APIVersion, obj.Kind = gvk.ToAPIVersionAndKind()
}

// GroupVersionKind returns the API version and kind of the object reference
func (obj *ObjectReference) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(obj.APIVersion, obj.Kind)
}

// GetObjectKind returns the kind of object reference
func (obj *ObjectReference) GetObjectKind() schema.ObjectKind { return obj }
