package container

import (
	"fmt"

	"k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	ref "k8s.io/client-go/tools/reference"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/features"
)

var ImplicitContainerPrefix string = "implicitly required container "

// GenerateContainerRef 其实就是根据传入的 Pod 对象, 选取其中的 container 部分,
// 组成一个 container Ref 对象.
// 其实这个 Ref 的 Kind 都是 Pod, 不过其中还包括了该 Pod 的 uid, name, ns 等信息.
//
// GenerateContainerRef returns an *v1.ObjectReference which references
// the given container within the given pod.
// Returns an error if the reference can't be constructed or 
// the container doesn't actually belong to the pod.
//
// This function will return an error if the provided Pod does not have a selfLink,
// but we expect selfLink to be populated at all call sites for the function.
func GenerateContainerRef(pod *v1.Pod, container *v1.Container) (*v1.ObjectReference, error) {
	fieldPath, err := fieldPath(pod, container)
	if err != nil {
		// TODO: figure out intelligent way to refer to containers that we implicitly
		// start (like the pod infra container). This is not a good way, ugh.
		fieldPath = ImplicitContainerPrefix + container.Name
	}
	ref, err := ref.GetPartialReference(legacyscheme.Scheme, pod, fieldPath)
	if err != nil {
		return nil, err
	}
	return ref, nil
}

// fieldPath returns a fieldPath locating container within pod.
// Returns an error if the container isn't part of the pod.
func fieldPath(pod *v1.Pod, container *v1.Container) (string, error) {
	for i := range pod.Spec.Containers {
		here := &pod.Spec.Containers[i]
		if here.Name == container.Name {
			if here.Name == "" {
				return fmt.Sprintf("spec.containers[%d]", i), nil
			}
			return fmt.Sprintf("spec.containers{%s}", here.Name), nil
		}
	}
	for i := range pod.Spec.InitContainers {
		here := &pod.Spec.InitContainers[i]
		if here.Name == container.Name {
			if here.Name == "" {
				return fmt.Sprintf("spec.initContainers[%d]", i), nil
			}
			return fmt.Sprintf("spec.initContainers{%s}", here.Name), nil
		}
	}
	if utilfeature.DefaultFeatureGate.Enabled(features.EphemeralContainers) {
		for i := range pod.Spec.EphemeralContainers {
			here := &pod.Spec.EphemeralContainers[i]
			if here.Name == container.Name {
				if here.Name == "" {
					return fmt.Sprintf("spec.ephemeralContainers[%d]", i), nil
				}
				return fmt.Sprintf("spec.ephemeralContainers{%s}", here.Name), nil
			}
		}
	}
	return "", fmt.Errorf("container %q not found in pod %s/%s", container.Name, pod.Namespace, pod.Name)
}
