package https

import (
	"io/ioutil"
	"net/http"
	"time"

	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubernetes/cmd/kubeadm/app/discovery/file"
)

// RetrieveValidatedConfigInfo connects to the API Server and makes sure it can talk
// securely to the API Server using the provided CA cert and
// optionally refreshes the cluster-info information from the cluster-info ConfigMap
func RetrieveValidatedConfigInfo(httpsURL, clustername string, discoveryTimeout time.Duration) (*clientcmdapi.Config, error) {
	client := &http.Client{Transport: netutil.SetOldTransportDefaults(&http.Transport{})}
	response, err := client.Get(httpsURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	kubeconfig, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return nil, err
	}
	return file.ValidateConfigInfo(config, clustername, discoveryTimeout)
}
