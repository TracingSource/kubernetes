package phases

import (
	"io"

	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// 	@implementBy: cmd/kubeadm/app/cmd/join.go -> joinData{}
//
// JoinData is the interface to use for join phases.
// The "joinData" type from "cmd/join.go" must satisfy this interface.
type JoinData interface {
	CertificateKey() string
	Cfg() *kubeadmapi.JoinConfiguration
	TLSBootstrapCfg() (*clientcmdapi.Config, error)
	InitCfg() (*kubeadmapi.InitConfiguration, error)
	ClientSet() (*clientset.Clientset, error)
	IgnorePreflightErrors() sets.String
	OutputWriter() io.Writer
	KustomizeDir() string
}
