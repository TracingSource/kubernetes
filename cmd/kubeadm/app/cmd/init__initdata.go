package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
)

// 	@implementOf: cmd/kubeadm/app/cmd/phases/init/data.go -> InitData
//
// initData defines all the runtime information used when running the kubeadm init workflow;
// this data is shared across all the phases that are included in the workflow.
type initData struct {
	// cfg kubeadm init 时通过 --config 指定的配置文件(不过是经过填充的, 包含全部内容)
	cfg                     *kubeadmapi.InitConfiguration
	skipTokenPrint          bool
	dryRun                  bool
	dryRunDir               string
	kubeconfigDir           string // /etc/kubernetes
	kubeconfigPath          string // /etc/kubernetes/admin.conf
	certificatesDir         string // /etc/kubernetes/pki
	ignorePreflightErrors   sets.String
	externalCA              bool
	client                  clientset.Interface
	outputWriter            io.Writer
	uploadCerts             bool
	skipCertificateKeyPrint bool
	kustomizeDir            string
}

// UploadCerts returns Uploadcerts flag.
func (d *initData) UploadCerts() bool {
	return d.uploadCerts
}

// CertificateKey returns the key used to encrypt the certs.
func (d *initData) CertificateKey() string {
	return d.cfg.CertificateKey
}

// SetCertificateKey set the key used to encrypt the certs.
func (d *initData) SetCertificateKey(key string) {
	d.cfg.CertificateKey = key
}

// SkipCertificateKeyPrint returns the skipCertificateKeyPrint flag.
func (d *initData) SkipCertificateKeyPrint() bool {
	return d.skipCertificateKeyPrint
}

// Cfg returns initConfiguration.
func (d *initData) Cfg() *kubeadmapi.InitConfiguration {
	return d.cfg
}

// DryRun returns the DryRun flag.
func (d *initData) DryRun() bool {
	return d.dryRun
}

// SkipTokenPrint returns the SkipTokenPrint flag.
func (d *initData) SkipTokenPrint() bool {
	return d.skipTokenPrint
}

// IgnorePreflightErrors returns the IgnorePreflightErrors flag.
func (d *initData) IgnorePreflightErrors() sets.String {
	return d.ignorePreflightErrors
}

// CertificateWriteDir returns the path to the certificate folder
// or the temporary folder path in case of DryRun.
func (d *initData) CertificateWriteDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return d.certificatesDir
}

// KubeConfigPath /etc/kubernetes/pki
//
// CertificateDir returns the CertificateDir as originally specified by the user.
func (d *initData) CertificateDir() string {
	return d.certificatesDir
}

// KubeConfigPath /etc/kubernetes
//
// KubeConfigDir returns the path of the Kubernetes configuration folder
// or the temporary folder path in case of DryRun.
func (d *initData) KubeConfigDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return d.kubeconfigDir
}

// KubeConfigPath /etc/kubernetes/admin.conf
//
// KubeConfigPath returns the path to the kubeconfig file to use for connecting to Kubernetes
func (d *initData) KubeConfigPath() string {
	if d.dryRun {
		d.kubeconfigPath = filepath.Join(d.dryRunDir, kubeadmconstants.AdminKubeConfigFileName)
	}
	return d.kubeconfigPath
}

// ManifestDir returns the path where manifest should be stored
// or the temporary folder path in case of DryRun.
func (d *initData) ManifestDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return kubeadmconstants.GetStaticPodDirectory()
}

// KubeletDir 一般返回 /var/lib/kubelet
//
// KubeletDir returns path of the kubelet configuration folder
// or the temporary folder in case of DryRun.
func (d *initData) KubeletDir() string {
	if d.dryRun {
		return d.dryRunDir
	}
	return kubeadmconstants.KubeletRunDirectory
}

// ExternalCA returns true if an external CA is provided by the user.
func (d *initData) ExternalCA() bool {
	return d.externalCA
}

// OutputWriter returns the io.Writer used to write output to by this command.
func (d *initData) OutputWriter() io.Writer {
	return d.outputWriter
}

// Client 读取 /etc/kubernetes/admin.conf 配置文件构建 kube client 并返回.
//
// Client returns a Kubernetes client to be used by kubeadm.
// This function is implemented as a singleton,
// thus avoiding to recreate the client when it is used by different phases.
// Important. This function must be called after the admin.conf kubeconfig file is created.
func (d *initData) Client() (clientset.Interface, error) {
	if d.client == nil {
		if d.dryRun {
			svcSubnetCIDR, err := kubeadmconstants.GetKubernetesServiceCIDR(
				d.cfg.Networking.ServiceSubnet,
				features.Enabled(d.cfg.FeatureGates, features.IPv6DualStack),
			)
			if err != nil {
				return nil, errors.Wrapf(
					err, "unable to get internal Kubernetes Service IP from the given service CIDR (%s)",
					d.cfg.Networking.ServiceSubnet,
				)
			}
			// If we're dry-running, we should create a faked client that answers some GETs
			// in order to be able to do the full init flow and just logs the rest of requests
			dryRunGetter := apiclient.NewInitDryRunGetter(
				d.cfg.NodeRegistration.Name, svcSubnetCIDR.String(),
			)
			d.client = apiclient.NewDryRunClient(dryRunGetter, os.Stdout)
		} else {
			// If we're acting for real, we should create a connection to the API server
			// and wait for it to come up
			var err error
			d.client, err = kubeconfigutil.ClientSetFromFile(d.KubeConfigPath())
			if err != nil {
				return nil, err
			}
		}
	}
	return d.client, nil
}

// Tokens returns an array of token strings.
func (d *initData) Tokens() []string {
	tokens := []string{}
	for _, bt := range d.cfg.BootstrapTokens {
		tokens = append(tokens, bt.Token.String())
	}
	return tokens
}

// KustomizeDir returns the folder where kustomize patches for static pod manifest are stored
func (d *initData) KustomizeDir() string {
	return d.kustomizeDir
}
