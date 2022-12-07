package testing

import (
	nodev1beta1 "k8s.io/api/node/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubernetes/pkg/kubelet/runtimeclass"
)

const (
	// SandboxRuntimeClass is a valid RuntimeClass pre-populated in the populated dynamic client.
	SandboxRuntimeClass = "sandbox"
	// SandboxRuntimeHandler is the handler associated with the SandboxRuntimeClass.
	SandboxRuntimeHandler = "kata-containers"

	// EmptyRuntimeClass is a valid RuntimeClass without a handler pre-populated in the populated dynamic client.
	EmptyRuntimeClass = "native"
)

// NewPopulatedClient creates a fake client for use with the runtimeclass.Manager,
// and populates it with a few test RuntimeClass objects.
func NewPopulatedClient() clientset.Interface {
	return fake.NewSimpleClientset(
		NewRuntimeClass(EmptyRuntimeClass, ""),
		NewRuntimeClass(SandboxRuntimeClass, SandboxRuntimeHandler),
	)
}

// StartManagerSync starts the manager, and waits for the informer cache to sync.
// Returns a function to stop the manager, which should be called with a defer:
//     defer StartManagerSync(t, m)()
func StartManagerSync(m *runtimeclass.Manager) func() {
	stopCh := make(chan struct{})
	m.Start(stopCh)
	m.WaitForCacheSync(stopCh)
	return func() {
		close(stopCh)
	}
}

// NewRuntimeClass is a helper to generate a RuntimeClass resource with
// the given name & handler.
func NewRuntimeClass(name, handler string) *nodev1beta1.RuntimeClass {
	return &nodev1beta1.RuntimeClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Handler: handler,
	}
}
