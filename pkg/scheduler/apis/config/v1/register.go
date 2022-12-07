package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeschedulerconfigv1 "k8s.io/kube-scheduler/config/v1"
)

// GroupName is the group name used in this package
const GroupName = "kubescheduler.config.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var (
	// localSchemeBuilder extends the SchemeBuilder instance with the external types. In this package,
	// defaulting and conversion init funcs are registered as well.
	localSchemeBuilder = &kubeschedulerconfigv1.SchemeBuilder
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(RegisterDefaults)
}
