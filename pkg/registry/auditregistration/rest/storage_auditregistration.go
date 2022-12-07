package rest

import (
	auditv1alpha1 "k8s.io/api/auditregistration/v1alpha1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	auditstorage "k8s.io/kubernetes/pkg/registry/auditregistration/auditsink/storage"
)

// RESTStorageProvider is a REST storage provider for audit.k8s.io
type RESTStorageProvider struct{}

// NewRESTStorage returns a RESTStorageProvider
func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(auditv1alpha1.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)

	if apiResourceConfigSource.VersionEnabled(auditv1alpha1.SchemeGroupVersion) {
		if storageMap, err := p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[auditv1alpha1.SchemeGroupVersion.Version] = storageMap
		}
	}
	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	s, err := auditstorage.NewREST(restOptionsGetter)
	storage["auditsinks"] = s

	return storage, err
}

// GroupName is the group name for the storage provider
func (p RESTStorageProvider) GroupName() string {
	return auditv1alpha1.GroupName
}
