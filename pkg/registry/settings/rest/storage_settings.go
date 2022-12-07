package rest

import (
	settingsapiv1alpha1 "k8s.io/api/settings/v1alpha1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/settings"
	podpresetstore "k8s.io/kubernetes/pkg/registry/settings/podpreset/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(settings.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(settingsapiv1alpha1.SchemeGroupVersion) {
		if storageMap, err := p.v1alpha1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, nil
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[settingsapiv1alpha1.SchemeGroupVersion.Version] = storageMap
		}
	}

	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1alpha1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// podpresets
	podPresetStorage, err := podpresetstore.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["podpresets"] = podPresetStorage

	return storage, err
}

func (p RESTStorageProvider) GroupName() string {
	return settings.GroupName
}
