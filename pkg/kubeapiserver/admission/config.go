package admission

import (
	"io/ioutil"
	"net/http"
	"time"

	"k8s.io/klog"

	utilwait "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/admission"
	webhookinit "k8s.io/apiserver/pkg/admission/plugin/webhook/initializer"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/util/webhook"
	cacheddiscovery "k8s.io/client-go/discovery/cached/memory"
	externalinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	quotainstall "k8s.io/kubernetes/pkg/quota/v1/install"
)

// Config holds the configuration needed to for initialize the admission plugins
type Config struct {
	CloudConfigFile      string
	LoopbackClientConfig *rest.Config
	ExternalInformers    externalinformers.SharedInformerFactory
}

// New sets up the plugins and admission start hooks needed for admission
func (c *Config) New(proxyTransport *http.Transport, serviceResolver webhook.ServiceResolver) ([]admission.PluginInitializer, genericapiserver.PostStartHookFunc, error) {
	webhookAuthResolverWrapper := webhook.NewDefaultAuthenticationInfoResolverWrapper(proxyTransport, c.LoopbackClientConfig)
	webhookPluginInitializer := webhookinit.NewPluginInitializer(webhookAuthResolverWrapper, serviceResolver)

	var cloudConfig []byte
	if c.CloudConfigFile != "" {
		var err error
		cloudConfig, err = ioutil.ReadFile(c.CloudConfigFile)
		if err != nil {
			klog.Fatalf("Error reading from cloud configuration file %s: %#v", c.CloudConfigFile, err)
		}
	}
	clientset, err := kubernetes.NewForConfig(c.LoopbackClientConfig)
	if err != nil {
		return nil, nil, err
	}

	discoveryClient := cacheddiscovery.NewMemCacheClient(clientset.Discovery())
	discoveryRESTMapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	kubePluginInitializer := NewPluginInitializer(
		cloudConfig,
		discoveryRESTMapper,
		quotainstall.NewQuotaConfigurationForAdmission(),
	)

	admissionPostStartHook := func(context genericapiserver.PostStartHookContext) error {
		discoveryRESTMapper.Reset()
		go utilwait.Until(discoveryRESTMapper.Reset, 30*time.Second, context.StopCh)
		return nil
	}

	return []admission.PluginInitializer{webhookPluginInitializer, kubePluginInitializer}, admissionPostStartHook, nil
}
