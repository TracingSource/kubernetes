package services

import (
	"fmt"
	"net"

	"k8s.io/apiserver/pkg/storage/storagebackend"
	apiserver "k8s.io/kubernetes/cmd/kube-apiserver/app"
	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"
)

const (
	clusterIPRange          = "10.0.0.1/24"
	apiserverClientURL      = "http://localhost:8080"
	apiserverHealthCheckURL = apiserverClientURL + "/healthz"
)

// APIServer is a server which manages apiserver.
type APIServer struct {
	storageConfig storagebackend.Config
	stopCh        chan struct{}
}

// NewAPIServer creates an apiserver.
func NewAPIServer(storageConfig storagebackend.Config) *APIServer {
	return &APIServer{
		storageConfig: storageConfig,
		stopCh:        make(chan struct{}),
	}
}

// Start starts the apiserver, returns when apiserver is ready.
func (a *APIServer) Start() error {
	o := options.NewServerRunOptions()
	o.Etcd.StorageConfig = a.storageConfig
	_, ipnet, err := net.ParseCIDR(clusterIPRange)
	if err != nil {
		return err
	}
	o.ServiceClusterIPRanges = ipnet.String()
	o.AllowPrivileged = true
	o.Admission.GenericAdmission.DisablePlugins = []string{"ServiceAccount", "TaintNodesByCondition"}
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		completedOptions, err := apiserver.Complete(o)
		if err != nil {
			errCh <- fmt.Errorf("set apiserver default options error: %v", err)
			return
		}
		err = apiserver.Run(completedOptions, a.stopCh)
		if err != nil {
			errCh <- fmt.Errorf("run apiserver error: %v", err)
			return
		}
	}()

	err = readinessCheck("apiserver", []string{apiserverHealthCheckURL}, errCh)
	if err != nil {
		return err
	}
	return nil
}

// Stop stops the apiserver. Currently, there is no way to stop the apiserver.
// The function is here only for completion.
func (a *APIServer) Stop() error {
	if a.stopCh != nil {
		close(a.stopCh)
		a.stopCh = nil
	}
	return nil
}

const apiserverName = "apiserver"

// Name returns the name of APIServer.
func (a *APIServer) Name() string {
	return apiserverName
}

func getAPIServerClientURL() string {
	return apiserverClientURL
}

func getAPIServerHealthCheckURL() string {
	return apiserverHealthCheckURL
}
