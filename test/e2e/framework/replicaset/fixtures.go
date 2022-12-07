package replicaset

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewReplicaSet returns a new ReplicaSet.
func NewReplicaSet(name, namespace string, replicas int32, podLabels map[string]string, imageName, image string) *appsv1.ReplicaSet {
	return &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: appsv1.ReplicaSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: podLabels,
			},
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            imageName,
							Image:           image,
							SecurityContext: &v1.SecurityContext{},
						},
					},
				},
			},
		},
	}
}
