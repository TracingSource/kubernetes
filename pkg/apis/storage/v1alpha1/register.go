package v1alpha1

import (
	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "storage.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	localSchemeBuilder = &storagev1alpha1.SchemeBuilder
	// AddToScheme is a common registration function for mapping packaged scoped group & version keys to a scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)
