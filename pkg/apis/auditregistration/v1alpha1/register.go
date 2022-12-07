package v1alpha1

import (
	auditregistrationv1alpha1 "k8s.io/api/auditregistration/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName for audit registration
const GroupName = "auditregistration.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	localSchemeBuilder = &auditregistrationv1alpha1.SchemeBuilder
	// AddToScheme audit registration
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addDefaultingFuncs)
}
