package resource

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

// TODO require negotiatedSerializer.  leaving it optional lets us plumb current behavior and deal with the difference after major plumbing is complete
func (clientConfigFn ClientConfigFunc) clientForGroupVersion(gv schema.GroupVersion, negotiatedSerializer runtime.NegotiatedSerializer) (RESTClient, error) {
	cfg, err := clientConfigFn()
	if err != nil {
		return nil, err
	}
	if negotiatedSerializer != nil {
		cfg.ContentConfig.NegotiatedSerializer = negotiatedSerializer
	}
	cfg.GroupVersion = &gv
	if len(gv.Group) == 0 {
		cfg.APIPath = "/api"
	} else {
		cfg.APIPath = "/apis"
	}

	return rest.RESTClientFor(cfg)
}

func (clientConfigFn ClientConfigFunc) unstructuredClientForGroupVersion(gv schema.GroupVersion) (RESTClient, error) {
	cfg, err := clientConfigFn()
	if err != nil {
		return nil, err
	}
	cfg.ContentConfig = UnstructuredPlusDefaultContentConfig()
	cfg.GroupVersion = &gv
	if len(gv.Group) == 0 {
		cfg.APIPath = "/api"
	} else {
		cfg.APIPath = "/apis"
	}

	return rest.RESTClientFor(cfg)
}
