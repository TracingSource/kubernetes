package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
)

// GroupName is the group name used in this package
const GroupName = "kubelet.config.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

var (
	// localSchemeBuilder extends the SchemeBuilder instance with the external types. In this package,
	// defaulting and conversion init funcs are registered as well.
	localSchemeBuilder = &kubeletconfigv1beta1.SchemeBuilder
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addDefaultingFuncs)
}
