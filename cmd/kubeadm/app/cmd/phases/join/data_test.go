package phases

import (
	"io"

	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// a package local type for testing purposes.
type testJoinData struct{}

// testJoinData must satisfy JoinData.
var _ JoinData = &testJoinData{}

func (j *testJoinData) CertificateKey() string                          { return "" }
func (j *testJoinData) Cfg() *kubeadmapi.JoinConfiguration              { return nil }
func (j *testJoinData) TLSBootstrapCfg() (*clientcmdapi.Config, error)  { return nil, nil }
func (j *testJoinData) InitCfg() (*kubeadmapi.InitConfiguration, error) { return nil, nil }
func (j *testJoinData) ClientSet() (*clientset.Clientset, error)        { return nil, nil }
func (j *testJoinData) IgnorePreflightErrors() sets.String              { return nil }
func (j *testJoinData) OutputWriter() io.Writer                         { return nil }
func (j *testJoinData) KustomizeDir() string                            { return "" }
