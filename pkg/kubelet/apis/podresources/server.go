package podresources

import (
	"context"

	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
)

// 	@implementBy: pkg/kubelet/cm/container_manager_linux.go -> containerManagerImpl{}
//
// DevicesProvider knows how to provide the devices used by the given container
type DevicesProvider interface {
	GetDevices(podUID, containerName string) []*v1alpha1.ContainerDevices
}

// PodsProvider knows how to provide the pods admitted by the node
type PodsProvider interface {
	GetPods() []*v1.Pod
}

// podResourcesServer implements PodResourcesListerServer
type podResourcesServer struct {
	podsProvider    PodsProvider
	devicesProvider DevicesProvider
}

// NewPodResourcesServer returns a PodResourcesListerServer which lists pods
// provided by the PodsProvider with device information provided by the DevicesProvider
func NewPodResourcesServer(
	podsProvider PodsProvider, devicesProvider DevicesProvider,
) v1alpha1.PodResourcesListerServer {
	return &podResourcesServer{
		podsProvider:    podsProvider,
		devicesProvider: devicesProvider,
	}
}

// List returns information about the resources assigned to pods on the node
func (p *podResourcesServer) List(
	ctx context.Context, req *v1alpha1.ListPodResourcesRequest,
) (*v1alpha1.ListPodResourcesResponse, error) {
	pods := p.podsProvider.GetPods()
	podResources := make([]*v1alpha1.PodResources, len(pods))

	for i, pod := range pods {
		pRes := v1alpha1.PodResources{
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			Containers: make([]*v1alpha1.ContainerResources, len(pod.Spec.Containers)),
		}

		for j, container := range pod.Spec.Containers {
			pRes.Containers[j] = &v1alpha1.ContainerResources{
				Name:    container.Name,
				Devices: p.devicesProvider.GetDevices(string(pod.UID), container.Name),
			}
		}
		podResources[i] = &pRes
	}

	return &v1alpha1.ListPodResourcesResponse{
		PodResources: podResources,
	}, nil
}
