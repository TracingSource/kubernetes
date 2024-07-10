package phases

import (
	"io"

	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// 	@implementBy: cmd/kubeadm/app/cmd/init.go -> initData{}
//
// InitData is the interface to use for init phases.
// The "initData" type from "cmd/init.go" must satisfy this interface.
type InitData interface {
	UploadCerts() bool
	CertificateKey() string
	SetCertificateKey(key string)
	SkipCertificateKeyPrint() bool
	// 返回结构体中的 cfg 字段, 该字段的内容即是 kubeadm-config.yaml 的配置对象
	Cfg() *kubeadmapi.InitConfiguration
	DryRun() bool
	SkipTokenPrint() bool
	IgnorePreflightErrors() sets.String
	CertificateWriteDir() string
	CertificateDir() string
	KubeConfigDir() string
	KubeConfigPath() string
	ManifestDir() string
	KubeletDir() string
	ExternalCA() bool
	OutputWriter() io.Writer
	Client() (clientset.Interface, error)
	Tokens() []string
	KustomizeDir() string
}
