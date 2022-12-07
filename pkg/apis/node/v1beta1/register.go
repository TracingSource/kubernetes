package v1beta1

import (
	nodev1beta1 "k8s.io/api/node/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName for node API
const GroupName = "node.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	localSchemeBuilder = &nodev1beta1.SchemeBuilder
	// AddToScheme node API registration
	AddToScheme = localSchemeBuilder.AddToScheme
)
