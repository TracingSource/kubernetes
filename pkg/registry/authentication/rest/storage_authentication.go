package rest

import (
	authenticationv1 "k8s.io/api/authentication/v1"
	authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/authentication"
	"k8s.io/kubernetes/pkg/registry/authentication/tokenreview"
)

type RESTStorageProvider struct {
	Authenticator authenticator.Request
	APIAudiences  authenticator.Audiences
}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	// TODO figure out how to make the swagger generation stable, while allowing this endpoint to be disabled.
	// if p.Authenticator == nil {
	// 	return genericapiserver.APIGroupInfo{}, false
	// }

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(authentication.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(authenticationv1beta1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[authenticationv1beta1.SchemeGroupVersion.Version] = p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
	}
	if apiResourceConfigSource.VersionEnabled(authenticationv1.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[authenticationv1.SchemeGroupVersion.Version] = p.v1Storage(apiResourceConfigSource, restOptionsGetter)
	}

	return apiGroupInfo, true, nil
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {
	storage := map[string]rest.Storage{}
	// tokenreviews
	tokenReviewStorage := tokenreview.NewREST(p.Authenticator, p.APIAudiences)
	storage["tokenreviews"] = tokenReviewStorage

	return storage
}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {
	storage := map[string]rest.Storage{}
	// tokenreviews
	tokenReviewStorage := tokenreview.NewREST(p.Authenticator, p.APIAudiences)
	storage["tokenreviews"] = tokenReviewStorage

	return storage
}

func (p RESTStorageProvider) GroupName() string {
	return authentication.GroupName
}
