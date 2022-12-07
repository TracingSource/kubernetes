package resourcequota

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeBuilder is a pointer used to call AddToScheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is used to register the types to API encoding/decoding machinery
	AddToScheme = SchemeBuilder.AddToScheme
)

// LegacyGroupName is the group name use in this package
const LegacyGroupName = "resourcequota.admission.k8s.io"

// LegacySchemeGroupVersion is group version used to register these objects
var LegacySchemeGroupVersion = schema.GroupVersion{Group: LegacyGroupName, Version: runtime.APIVersionInternal}

// GroupName is the group name use in this package
const GroupName = "apiserver.config.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(LegacySchemeGroupVersion,
		&Configuration{},
	)
	scheme.AddKnownTypeWithName(SchemeGroupVersion.WithKind("ResourceQuotaConfiguration"),
		&Configuration{},
	)
	return nil
}
