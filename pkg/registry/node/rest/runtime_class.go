package rest

import (
	nodev1alpha1 "k8s.io/api/node/v1alpha1"
	nodev1beta1 "k8s.io/api/node/v1beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	nodeinternal "k8s.io/kubernetes/pkg/apis/node"
	runtimeclassstorage "k8s.io/kubernetes/pkg/registry/node/runtimeclass/storage"
)

// RESTStorageProvider is a REST storage provider for node.k8s.io
type RESTStorageProvider struct{}

// NewRESTStorage returns a RESTStorageProvider
func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(nodeinternal.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	if apiResourceConfigSource.VersionEnabled(nodev1alpha1.SchemeGroupVersion) {
		if storageMap, err := p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[nodev1alpha1.SchemeGroupVersion.Version] = storageMap
		}
	}

	if apiResourceConfigSource.VersionEnabled(nodev1beta1.SchemeGroupVersion) {
		if storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[nodev1beta1.SchemeGroupVersion.Version] = storageMap
		}
	}

	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	s, err := runtimeclassstorage.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["runtimeclasses"] = s

	return storage, err
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	s, err := runtimeclassstorage.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["runtimeclasses"] = s

	return storage, err
}

// GroupName is the group name for the storage provider
func (p RESTStorageProvider) GroupName() string {
	return nodeinternal.GroupName
}
