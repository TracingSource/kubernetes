package rest

import (
	discoveryv1alpha1 "k8s.io/api/discovery/v1alpha1"
	discoveryv1beta1 "k8s.io/api/discovery/v1beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/discovery"
	endpointslicestorage "k8s.io/kubernetes/pkg/registry/discovery/endpointslice/storage"
)

// StorageProvider is a REST storage provider for discovery.k8s.io.
type StorageProvider struct{}

// NewRESTStorage returns a new storage provider.
func (p StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(discovery.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(discoveryv1alpha1.SchemeGroupVersion) {
		storageMap, err := p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter)
		if err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		}
		apiGroupInfo.VersionedResourcesStorageMap[discoveryv1alpha1.SchemeGroupVersion.Version] = storageMap
	}

	if apiResourceConfigSource.VersionEnabled(discoveryv1beta1.SchemeGroupVersion) {
		storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
		if err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		}
		apiGroupInfo.VersionedResourcesStorageMap[discoveryv1beta1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, true, nil
}

func (p StorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	endpointSliceStorage, err := endpointslicestorage.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}

	storage["endpointslices"] = endpointSliceStorage
	return storage, err
}

func (p StorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	endpointSliceStorage, err := endpointslicestorage.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}

	storage["endpointslices"] = endpointSliceStorage
	return storage, err
}

// GroupName is the group name for the storage provider.
func (p StorageProvider) GroupName() string {
	return discovery.GroupName
}
