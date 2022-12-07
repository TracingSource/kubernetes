package rest

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/registry/registrytest"
)

func TestPodLogValidates(t *testing.T) {
	config, server := registrytest.NewEtcdStorage(t, "")
	defer server.Terminate(t)
	s, destroyFunc, err := generic.NewRawStorage(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer destroyFunc()
	store := &genericregistry.Store{
		Storage: genericregistry.DryRunnableStorage{Storage: s},
	}
	logRest := &LogREST{Store: store, KubeletConn: nil}

	negativeOne := int64(-1)
	testCases := []*api.PodLogOptions{
		{SinceSeconds: &negativeOne},
		{TailLines: &negativeOne},
	}

	for _, tc := range testCases {
		_, err := logRest.Get(genericapirequest.NewDefaultContext(), "test", tc)
		if !errors.IsInvalid(err) {
			t.Fatalf("Unexpected error: %v", err)
		}
	}
}
