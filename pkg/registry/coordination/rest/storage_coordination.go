package rest

import (
	coordinationv1 "k8s.io/api/coordination/v1"
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/coordination"
	leasestorage "k8s.io/kubernetes/pkg/registry/coordination/lease/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(coordination.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(coordinationv1beta1.SchemeGroupVersion) {
		if storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[coordinationv1beta1.SchemeGroupVersion.Version] = storageMap
		}
	}
	if apiResourceConfigSource.VersionEnabled(coordinationv1.SchemeGroupVersion) {
		if storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[coordinationv1.SchemeGroupVersion.Version] = storageMap
		}
	}
	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// leases
	leaseStorage, err := leasestorage.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["leases"] = leaseStorage

	return storage, err
}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// leases
	leaseStorage, err := leasestorage.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["leases"] = leaseStorage

	return storage, err
}

func (p RESTStorageProvider) GroupName() string {
	return coordination.GroupName
}
