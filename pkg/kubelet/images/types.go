package images

import (
	"errors"

	"k8s.io/api/core/v1"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

var (
	// ErrImagePullBackOff - Container image pull failed, kubelet is backing off image pull
	ErrImagePullBackOff = errors.New("ImagePullBackOff")

	// ErrImageInspect - Unable to inspect image
	ErrImageInspect = errors.New("ImageInspectError")

	// ErrImagePull - General image pull error
	ErrImagePull = errors.New("ErrImagePull")

	// ErrImageNeverPull - Required Image is absent on host and PullPolicy is NeverPullImage
	ErrImageNeverPull = errors.New("ErrImageNeverPull")

	// ErrRegistryUnavailable - Get http error when pulling image from registry
	ErrRegistryUnavailable = errors.New("RegistryUnavailable")

	// ErrInvalidImageName - Unable to parse the image name.
	ErrInvalidImageName = errors.New("InvalidImageName")
)

// ImageManager provides an interface to manage the lifecycle of images.
// Implementations of this interface are expected to deal with pulling (downloading),
// managing, and deleting container images.
// Implementations are expected to abstract the underlying runtimes.
// Implementations are expected to be thread safe.
type ImageManager interface {
	// EnsureImageExists ensures that image specified in `container` exists.
	EnsureImageExists(pod *v1.Pod, container *v1.Container, pullSecrets []v1.Secret, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, string, error)

	// TODO(ronl): consolidating image managing and deleting operation in this interface
}
