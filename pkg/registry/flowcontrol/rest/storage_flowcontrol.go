package rest

import (
	flowcontrolv1alpha1 "k8s.io/api/flowcontrol/v1alpha1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/flowcontrol"
	flowschemastore "k8s.io/kubernetes/pkg/registry/flowcontrol/flowschema/storage"
	prioritylevelconfigurationstore "k8s.io/kubernetes/pkg/registry/flowcontrol/prioritylevelconfiguration/storage"
)

// RESTStorageProvider implements
type RESTStorageProvider struct{}

// NewRESTStorage creates a new rest storage for flow-control api models.
func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(flowcontrol.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	if apiResourceConfigSource.VersionEnabled(flowcontrolv1alpha1.SchemeGroupVersion) {
		flowControlStorage, err := p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter)
		if err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		}
		apiGroupInfo.VersionedResourcesStorageMap[flowcontrolv1alpha1.SchemeGroupVersion.Version] = flowControlStorage
	}
	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// flow-schema
	flowSchemaStorage, flowSchemaStatusStorage, err := flowschemastore.NewREST(restOptionsGetter)
	if err != nil {
		return nil, err
	}
	storage["flowschemas"] = flowSchemaStorage
	storage["flowschemas/status"] = flowSchemaStatusStorage

	// priority-level-configuration
	priorityLevelConfigurationStorage, priorityLevelConfigurationStatusStorage, err := prioritylevelconfigurationstore.NewREST(restOptionsGetter)
	if err != nil {
		return nil, err
	}
	storage["prioritylevelconfigurations"] = priorityLevelConfigurationStorage
	storage["prioritylevelconfigurations/status"] = priorityLevelConfigurationStatusStorage

	return storage, nil
}

// GroupName returns group name of the storage
func (p RESTStorageProvider) GroupName() string {
	return flowcontrol.GroupName
}
