package openshift

import (
	"testing"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/kubernetes/pkg/master"
)

// This test references methods that OpenShift uses to customize the master on startup, that
// are not referenced directly by a master.
func TestMasterExportsSymbols(t *testing.T) {
	_ = &master.Config{
		GenericConfig: &genericapiserver.Config{
			EnableMetrics: true,
		},
		ExtraConfig: master.ExtraConfig{
			EnableLogsSupport: false,
		},
	}
	_ = &master.Master{
		GenericAPIServer: &genericapiserver.GenericAPIServer{},
	}
}
