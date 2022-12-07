package reconciliation

import (
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// tryEnsureNamespace gets or creates the given namespace while ignoring forbidden errors.
// It is a best effort attempt as the user may not be able to get or create namespaces.
// This allows us to handle flows where the user can only mutate roles and role bindings.
func tryEnsureNamespace(client corev1client.NamespaceInterface, namespace string) error {
	_, getErr := client.Get(namespace, metav1.GetOptions{})
	if getErr == nil {
		return nil
	}

	if fatalGetErr := utilerrors.FilterOut(getErr, apierrors.IsNotFound, apierrors.IsForbidden); fatalGetErr != nil {
		return fatalGetErr
	}

	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}
	_, createErr := client.Create(ns)

	return utilerrors.FilterOut(createErr, apierrors.IsAlreadyExists, apierrors.IsForbidden)
}
