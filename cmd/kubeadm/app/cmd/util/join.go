package util

import (
	"bytes"
	"crypto/x509"
	"html/template"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcertutil "k8s.io/client-go/util/cert"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pubkeypin"
)

/*
 kubeadm join 172.16.42.39:6443 --token fl6bwd.3ube1uwcyk3a5o6b \
    --discovery-token-ca-cert-hash sha256:4eb77bb5abf037c8643d21efca30062c57862ab15d112e3166e88bf190f9e936 \
    --control-plane
*/
var joinCommandTemplate = template.Must(template.New("join").Parse(`` +
	`kubeadm join {{.ControlPlaneHostPort}} --token {{.Token}} \
    {{range $h := .CAPubKeyPins}}--discovery-token-ca-cert-hash {{$h}} {{end}}{{if .ControlPlane}}\
    --control-plane {{if .CertificateKey}}--certificate-key {{.CertificateKey}}{{end}}{{end}}`,
))

// GetJoinWorkerCommand returns the kubeadm join command for a given token and
// and Kubernetes cluster (the current cluster in the kubeconfig file)
func GetJoinWorkerCommand(kubeConfigFile, token string, skipTokenPrint bool) (string, error) {
	return getJoinCommand(kubeConfigFile, token, "", false, skipTokenPrint, false)
}

// GetJoinControlPlaneCommand ...
//
// 	@param kubeConfigFile: /etc/kubernetes/admin.conf
// 	@param token: kube-system/bootstrap-token-${xxx} Secret 的内容.
//
// caller:
// 	1. cmd/kubeadm/app/cmd/init.go -> printJoinCommand()
// 	kubeadm init 构建完成的最后一步被调用, 打印 join 命令
//
// GetJoinControlPlaneCommand returns the kubeadm join command for a given token and
// and Kubernetes cluster (the current cluster in the kubeconfig file)
func GetJoinControlPlaneCommand(
	kubeConfigFile, token, key string, skipTokenPrint, skipCertificateKeyPrint bool,
) (string, error) {
	return getJoinCommand(
		kubeConfigFile, token, key,
		true, skipTokenPrint, skipCertificateKeyPrint,
	)
}

// getJoinCommand ...参数介绍见上面的主调函数
//
// caller:
// 	1. GetJoinControlPlaneCommand()
//
func getJoinCommand(
	kubeConfigFile, token, key string,
	controlPlane, skipTokenPrint, skipCertificateKeyPrint bool,
) (string, error) {
	// load the kubeconfig file to get the CA certificate and endpoint
	config, err := clientcmd.LoadFromFile(kubeConfigFile)
	if err != nil {
		return "", errors.Wrap(err, "failed to load kubeconfig")
	}

	// load the default cluster config
	clusterConfig := kubeconfigutil.GetClusterFromKubeConfig(config)
	if clusterConfig == nil {
		return "", errors.New("failed to get default cluster config")
	}

	// load CA certificates from the kubeconfig (either from PEM data or by file path)
	var caCerts []*x509.Certificate
	if clusterConfig.CertificateAuthorityData != nil {
		caCerts, err = clientcertutil.ParseCertsPEM(clusterConfig.CertificateAuthorityData)
		if err != nil {
			return "", errors.Wrap(err, "failed to parse CA certificate from kubeconfig")
		}
	} else if clusterConfig.CertificateAuthority != "" {
		caCerts, err = clientcertutil.CertsFromFile(clusterConfig.CertificateAuthority)
		if err != nil {
			return "", errors.Wrap(err, "failed to load CA certificate referenced by kubeconfig")
		}
	} else {
		return "", errors.New("no CA certificates found in kubeconfig")
	}

	// hash all the CA certs and include their public key pins as trusted values
	publicKeyPins := make([]string, 0, len(caCerts))
	for _, caCert := range caCerts {
		publicKeyPins = append(publicKeyPins, pubkeypin.Hash(caCert))
	}

	ctx := map[string]interface{}{
		"Token":                token,
		"CAPubKeyPins":         publicKeyPins,
		"ControlPlaneHostPort": strings.Replace(clusterConfig.Server, "https://", "", -1),
		"CertificateKey":       key,
		"ControlPlane":         controlPlane,
	}

	if skipTokenPrint {
		ctx["Token"] = template.HTML("<value withheld>")
	}
	if skipCertificateKeyPrint {
		ctx["CertificateKey"] = template.HTML("<value withheld>")
	}

	var out bytes.Buffer
	err = joinCommandTemplate.Execute(&out, ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to render join command template")
	}
	return out.String(), nil
}
