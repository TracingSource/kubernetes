package services

import (
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/metadata"
	restclient "k8s.io/client-go/rest"
	namespacecontroller "k8s.io/kubernetes/pkg/controller/namespace"
)

const (
	// ncName is the name of namespace controller
	ncName = "namespace-controller"
	// ncResyncPeriod is resync period of the namespace controller
	ncResyncPeriod = 5 * time.Minute
	// ncConcurrency is concurrency of the namespace controller
	ncConcurrency = 2
)

// NamespaceController is a server which manages namespace controller.
type NamespaceController struct {
	host   string
	stopCh chan struct{}
}

// NewNamespaceController creates a new namespace controller.
func NewNamespaceController(host string) *NamespaceController {
	return &NamespaceController{host: host, stopCh: make(chan struct{})}
}

// Start starts the namespace controller.
func (n *NamespaceController) Start() error {
	config := restclient.AddUserAgent(&restclient.Config{Host: n.host}, ncName)

	// the namespace cleanup controller is very chatty.  It makes lots of discovery calls and then it makes lots of delete calls.
	config.QPS = 50
	config.Burst = 200

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}
	metadataClient, err := metadata.NewForConfig(config)
	if err != nil {
		return err
	}
	discoverResourcesFn := client.Discovery().ServerPreferredNamespacedResources
	informerFactory := informers.NewSharedInformerFactory(client, ncResyncPeriod)
	nc := namespacecontroller.NewNamespaceController(
		client,
		metadataClient,
		discoverResourcesFn,
		informerFactory.Core().V1().Namespaces(),
		ncResyncPeriod, v1.FinalizerKubernetes,
	)
	informerFactory.Start(n.stopCh)
	go nc.Run(ncConcurrency, n.stopCh)
	return nil
}

// Stop stops the namespace controller.
func (n *NamespaceController) Stop() error {
	close(n.stopCh)
	return nil
}

// Name returns the name of namespace controller.
func (n *NamespaceController) Name() string {
	return ncName
}
