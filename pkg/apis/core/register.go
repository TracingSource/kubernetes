package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = ""

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

var (
	// SchemeBuilder object to register various known types
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme represents a func that can be used to apply all the registered
	// funcs in a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	if err := scheme.AddIgnoredConversionType(&metav1.TypeMeta{}, &metav1.TypeMeta{}); err != nil {
		return err
	}
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Pod{},
		&PodList{},
		&PodStatusResult{},
		&PodTemplate{},
		&PodTemplateList{},
		&ReplicationControllerList{},
		&ReplicationController{},
		&ServiceList{},
		&Service{},
		&ServiceProxyOptions{},
		&NodeList{},
		&Node{},
		&NodeProxyOptions{},
		&Endpoints{},
		&EndpointsList{},
		&Binding{},
		&Event{},
		&EventList{},
		&List{},
		&LimitRange{},
		&LimitRangeList{},
		&ResourceQuota{},
		&ResourceQuotaList{},
		&Namespace{},
		&NamespaceList{},
		&ServiceAccount{},
		&ServiceAccountList{},
		&Secret{},
		&SecretList{},
		&PersistentVolume{},
		&PersistentVolumeList{},
		&PersistentVolumeClaim{},
		&PersistentVolumeClaimList{},
		&PodAttachOptions{},
		&PodLogOptions{},
		&PodExecOptions{},
		&PodPortForwardOptions{},
		&PodProxyOptions{},
		&ComponentStatus{},
		&ComponentStatusList{},
		&SerializedReference{},
		&RangeAllocation{},
		&ConfigMap{},
		&ConfigMapList{},
		&EphemeralContainers{},
	)

	return nil
}
