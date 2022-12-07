package protectionutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/util/slice"
)

// IsDeletionCandidate checks if object is candidate to be deleted
func IsDeletionCandidate(obj metav1.Object, finalizer string) bool {
	return obj.GetDeletionTimestamp() != nil && slice.ContainsString(obj.GetFinalizers(),
		finalizer, nil)
}

// NeedToAddFinalizer checks if need to add finalizer to object
func NeedToAddFinalizer(obj metav1.Object, finalizer string) bool {
	return obj.GetDeletionTimestamp() == nil && !slice.ContainsString(obj.GetFinalizers(),
		finalizer, nil)
}
