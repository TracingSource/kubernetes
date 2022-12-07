package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/apis/abac"
)

// GroupName is the group name use in this package
const GroupName = "abac.authorization.kubernetes.io"

// SchemeGroupVersion is the API group and version for abac v1beta1
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

func init() {
	// TODO: Delete this init function, abac should not have its own scheme.
	utilruntime.Must(addKnownTypes(abac.Scheme))

	utilruntime.Must(RegisterConversions(abac.Scheme))
}

var (
	// SchemeBuilder is the scheme builder with scheme init functions to run for this API package
	// TODO: move SchemeBuilder with zz_generated.deepcopy.go to k8s.io/api.
	SchemeBuilder runtime.SchemeBuilder
	// localSchemeBuilder ïs a pointer to SchemeBuilder instance. Using localSchemeBuilder
	// defaulting and conversion init funcs are registered as well.
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.
	localSchemeBuilder = &SchemeBuilder
	// AddToScheme is a common registration function for mapping packaged scoped group & version keys to a scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes, RegisterDefaults)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Policy{},
	)
	return nil
}
