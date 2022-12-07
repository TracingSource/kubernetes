package rest

import (
	networkingapiv1 "k8s.io/api/networking/v1"
	networkingapiv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/networking"
	ingressstore "k8s.io/kubernetes/pkg/registry/networking/ingress/storage"
	networkpolicystore "k8s.io/kubernetes/pkg/registry/networking/networkpolicy/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(networking.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(networkingapiv1.SchemeGroupVersion) {
		if storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[networkingapiv1.SchemeGroupVersion.Version] = storageMap
		}
	}

	if apiResourceConfigSource.VersionEnabled(networkingapiv1beta1.SchemeGroupVersion) {
		if storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
			return genericapiserver.APIGroupInfo{}, false, err
		} else {
			apiGroupInfo.VersionedResourcesStorageMap[networkingapiv1beta1.SchemeGroupVersion.Version] = storageMap
		}
	}

	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// networkpolicies
	networkPolicyStorage, err := networkpolicystore.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["networkpolicies"] = networkPolicyStorage

	return storage, nil
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}
	// ingresses
	ingressStorage, ingressStatusStorage, err := ingressstore.NewREST(restOptionsGetter)
	if err != nil {
		return storage, err
	}
	storage["ingresses"] = ingressStorage
	storage["ingresses/status"] = ingressStatusStorage
	return storage, nil
}

func (p RESTStorageProvider) GroupName() string {
	return networking.GroupName
}
