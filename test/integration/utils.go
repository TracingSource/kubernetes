package integration

import (
	"testing"
	"time"

	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	clientset "k8s.io/client-go/kubernetes"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

// DeletePodOrErrorf deletes a pod or fails with a call to t.Errorf.
func DeletePodOrErrorf(t *testing.T, c clientset.Interface, ns, name string) {
	if err := c.CoreV1().Pods(ns).Delete(name, nil); err != nil {
		t.Errorf("unable to delete pod %v: %v", name, err)
	}
}

// Requests to try.  Each one should be forbidden or not forbidden
// depending on the authentication and authorization setup of the master.
var (
	Code200 = map[int]bool{200: true}
	Code201 = map[int]bool{201: true}
	Code400 = map[int]bool{400: true}
	Code401 = map[int]bool{401: true}
	Code403 = map[int]bool{403: true}
	Code404 = map[int]bool{404: true}
	Code405 = map[int]bool{405: true}
	Code409 = map[int]bool{409: true}
	Code422 = map[int]bool{422: true}
	Code500 = map[int]bool{500: true}
	Code503 = map[int]bool{503: true}
)

// WaitForPodToDisappear polls the API server if the pod has been deleted.
func WaitForPodToDisappear(podClient coreclient.PodInterface, podName string, interval, timeout time.Duration) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		_, err := podClient.Get(podName, metav1.GetOptions{})
		if err == nil {
			return false, nil
		}
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}

// GetEtcdClients returns an initialized  clientv3.Client and clientv3.KV.
func GetEtcdClients(config storagebackend.TransportConfig) (*clientv3.Client, clientv3.KV, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:      config.CertFile,
		KeyFile:       config.KeyFile,
		TrustedCAFile: config.TrustedCAFile,
	}

	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	cfg := clientv3.Config{
		Endpoints:   config.ServerList,
		DialTimeout: 20 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(), // block until the underlying connection is up
		},
		TLS: tlsConfig,
	}

	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, nil, err
	}

	return c, clientv3.NewKV(c), nil
}
