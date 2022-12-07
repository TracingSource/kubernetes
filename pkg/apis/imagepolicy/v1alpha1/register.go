package v1alpha1

import (
	imagepolicyv1alpha1 "k8s.io/api/imagepolicy/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name for this API.
const GroupName = "imagepolicy.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	localSchemeBuilder = &imagepolicyv1alpha1.SchemeBuilder
	// AddToScheme is a common registration function for mapping packaged scoped group & version keys to a scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(RegisterDefaults)
}
