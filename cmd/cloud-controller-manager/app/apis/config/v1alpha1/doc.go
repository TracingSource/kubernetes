// Note: The referenced generic ComponentConfig packages with conversions
// between the types (e.g. the external package) needs to be given as an
// input to conversion-gen for it to find the native conversation funcs to
// call.

// +k8s:deepcopy-gen=package
// +k8s:conversion-gen=k8s.io/kubernetes/cmd/cloud-controller-manager/app/apis/config
// +k8s:conversion-gen=k8s.io/component-base/config/v1alpha1
// +k8s:conversion-gen=k8s.io/kubernetes/pkg/controller/apis/config/v1alpha1
// +k8s:conversion-gen=k8s.io/kubernetes/pkg/controller/service/config/v1alpha1
// +k8s:openapi-gen=true
// +k8s:defaulter-gen=TypeMeta
// +groupName=cloudcontrollermanager.config.k8s.io

package v1alpha1 // import "k8s.io/kubernetes/cmd/cloud-controller-manager/app/apis/config/v1alpha1"
