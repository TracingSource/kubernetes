package rest

import (
	autoscalingapiv1 "k8s.io/api/autoscaling/v1"
	autoscalingapiv2beta1 "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	autoscalingapiv2beta2 "k8s.io/kubernetes/pkg/apis/autoscaling/v2beta2"
	horizontalpodautoscalerstore "k8s.io/kubernetes/pkg/registry/autoscaling/horizontalpodautoscaler/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(autoscaling.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(autoscalingapiv2beta2.SchemeGroupVersion) {
		if storageMap, err := p.v2beta2Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[autoscalingapiv2beta2.SchemeGroupVersion.Version] = storageMap
		}
	}
	if apiResourceConfigSource.VersionEnabled(autoscalingapiv2beta1.SchemeGroupVersion) {
		if storageMap, err := p.v2beta1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[autoscalingapiv2beta1.SchemeGroupVersion.Version] = storageMap
		}
	}
	if apiResourceConfigSource.VersionEnabled(autoscalingapiv1.SchemeGroupVersion) {
		if storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[autoscalingapiv1.SchemeGroupVersion.Version] = storageMap
		}
	}

	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// horizontalpodautoscalers
	hpaStorage, hpaStatusStorage, err := horizontalpodautoscalerstore.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["horizontalpodautoscalers"] = hpaStorage
	storage["horizontalpodautoscalers/status"] = hpaStatusStorage

	return storage, err
}

func (p RESTStorageProvider) v2beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// horizontalpodautoscalers
	hpaStorage, hpaStatusStorage, err := horizontalpodautoscalerstore.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["horizontalpodautoscalers"] = hpaStorage
	storage["horizontalpodautoscalers/status"] = hpaStatusStorage

	return storage, err
}

func (p RESTStorageProvider) v2beta2Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// horizontalpodautoscalers
	hpaStorage, hpaStatusStorage, err := horizontalpodautoscalerstore.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["horizontalpodautoscalers"] = hpaStorage
	storage["horizontalpodautoscalers/status"] = hpaStatusStorage

	return storage, err
}

func (p RESTStorageProvider) GroupName() string {
	return autoscaling.GroupName
}
