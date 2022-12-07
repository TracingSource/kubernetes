package extensions

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	"k8s.io/kubernetes/pkg/apis/networking"
	"k8s.io/kubernetes/pkg/apis/policy"
)

// GroupName is the group name use in this package
const GroupName = "extensions"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Builds new Scheme of known types
var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	// TODO this gets cleaned up when the types are fixed
	scheme.AddKnownTypes(SchemeGroupVersion,
		&apps.Deployment{},
		&apps.DeploymentList{},
		&apps.DeploymentRollback{},
		&ReplicationControllerDummy{},
		&apps.DaemonSetList{},
		&apps.DaemonSet{},
		&networking.Ingress{},
		&networking.IngressList{},
		&apps.ReplicaSet{},
		&apps.ReplicaSetList{},
		&policy.PodSecurityPolicy{},
		&policy.PodSecurityPolicyList{},
		&autoscaling.Scale{},
		&networking.NetworkPolicy{},
		&networking.NetworkPolicyList{},
	)
	return nil
}
