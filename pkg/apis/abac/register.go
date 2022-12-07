package abac

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// GroupName is the API group for abac
const GroupName = "abac.authorization.kubernetes.io"

// SchemeGroupVersion is the API group version used to register abac internal
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Scheme is the default instance of runtime.Scheme to which types in the abac API group are api.Registry.
// TODO: remove this, abac should not have its own scheme.
var Scheme = runtime.NewScheme()

// Codecs provides access to encoding and decoding for the scheme
var Codecs = serializer.NewCodecFactory(Scheme)

func init() {
	// TODO: delete this, abac should not have its own scheme.
	addKnownTypes(Scheme)
}

var (
	// SchemeBuilder is the scheme builder with scheme init functions to run for this API package
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a common registration function for mapping packaged scoped group & version keys to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Policy{},
	)
	return nil
}
