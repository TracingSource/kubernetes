package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "cloudcontrollermanager.config.k8s.io"

var (
	// GroupName is the group name use in this package
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	// SchemeBuilder is the scheme builder with scheme init functions to run for this API package
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// localSchemeBuilder ïs a pointer to SchemeBuilder instance. Using localSchemeBuilder
	// defaulting and conversion init funcs are registered as well.
	localSchemeBuilder = &SchemeBuilder
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addDefaultingFuncs)
}

// addKnownTypes registers known types to the given scheme
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CloudControllerManagerConfiguration{},
	)
	return nil
}
