package phases

import (
	"io"

	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// a package local type for testing purposes.
type testInitData struct{}

// testInitData must satisfy InitData.
var _ InitData = &testInitData{}

func (t *testInitData) UploadCerts() bool                    { return false }
func (t *testInitData) CertificateKey() string               { return "" }
func (t *testInitData) SetCertificateKey(key string)         {}
func (t *testInitData) SkipCertificateKeyPrint() bool        { return false }
func (t *testInitData) Cfg() *kubeadmapi.InitConfiguration   { return nil }
func (t *testInitData) DryRun() bool                         { return false }
func (t *testInitData) SkipTokenPrint() bool                 { return false }
func (t *testInitData) IgnorePreflightErrors() sets.String   { return nil }
func (t *testInitData) CertificateWriteDir() string          { return "" }
func (t *testInitData) CertificateDir() string               { return "" }
func (t *testInitData) KubeConfigDir() string                { return "" }
func (t *testInitData) KubeConfigPath() string               { return "" }
func (t *testInitData) ManifestDir() string                  { return "" }
func (t *testInitData) KubeletDir() string                   { return "" }
func (t *testInitData) ExternalCA() bool                     { return false }
func (t *testInitData) OutputWriter() io.Writer              { return nil }
func (t *testInitData) Client() (clientset.Interface, error) { return nil, nil }
func (t *testInitData) Tokens() []string                     { return nil }
func (t *testInitData) KustomizeDir() string                 { return "" }
