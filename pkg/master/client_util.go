package master

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

func createNamespaceIfNeeded(c corev1client.NamespacesGetter, ns string) error {
	if _, err := c.Namespaces().Get(ns, metav1.GetOptions{}); err == nil {
		// the namespace already exists
		return nil
	}
	newNs := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ns,
			Namespace: "",
		},
	}
	_, err := c.Namespaces().Create(newNs)
	if err != nil && errors.IsAlreadyExists(err) {
		err = nil
	}
	return err
}
