package phases

import (
	"io"

	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// resetData is the interface to use for reset phases.
// The "resetData" type from "cmd/reset.go" must satisfy this interface.
type resetData interface {
	ForceReset() bool
	InputReader() io.Reader
	IgnorePreflightErrors() sets.String
	Cfg() *kubeadmapi.InitConfiguration
	Client() clientset.Interface
	AddDirsToClean(dirs ...string)
	CertificatesDir() string
	CRISocketPath() string
}
