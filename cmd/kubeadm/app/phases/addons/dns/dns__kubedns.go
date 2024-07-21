package dns

import (
	"fmt"

	"github.com/pkg/errors"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	"k8s.io/kubernetes/cmd/kubeadm/app/images"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	utilsnet "k8s.io/utils/net"
)

func kubeDNSAddon(cfg *kubeadmapi.ClusterConfiguration, client clientset.Interface) error {
	if err := CreateServiceAccount(client); err != nil {
		return err
	}

	dnsip, err := kubeadmconstants.GetDNSIP(cfg.Networking.ServiceSubnet, features.Enabled(cfg.FeatureGates, features.IPv6DualStack))
	if err != nil {
		return err
	}

	var dnsBindAddr, dnsProbeAddr string
	if utilsnet.IsIPv6(dnsip) {
		dnsBindAddr = "::1"
		dnsProbeAddr = "[" + dnsBindAddr + "]"
	} else {
		dnsBindAddr = "127.0.0.1"
		dnsProbeAddr = dnsBindAddr
	}

	dnsDeploymentBytes, err := kubeadmutil.ParseTemplate(KubeDNSDeployment,
		struct{ DeploymentName, KubeDNSImage, DNSMasqImage, SidecarImage, DNSBindAddr, DNSProbeAddr, DNSDomain, ControlPlaneTaintKey string }{
			DeploymentName:       kubeadmconstants.KubeDNSDeploymentName,
			KubeDNSImage:         images.GetDNSImage(cfg, kubeadmconstants.KubeDNSKubeDNSImageName),
			DNSMasqImage:         images.GetDNSImage(cfg, kubeadmconstants.KubeDNSDnsMasqNannyImageName),
			SidecarImage:         images.GetDNSImage(cfg, kubeadmconstants.KubeDNSSidecarImageName),
			DNSBindAddr:          dnsBindAddr,
			DNSProbeAddr:         dnsProbeAddr,
			DNSDomain:            cfg.Networking.DNSDomain,
			ControlPlaneTaintKey: kubeadmconstants.LabelNodeRoleMaster,
		})
	if err != nil {
		return errors.Wrap(err, "error when parsing kube-dns deployment template")
	}

	dnsServiceBytes, err := kubeadmutil.ParseTemplate(KubeDNSService, struct{ DNSIP string }{
		DNSIP: dnsip.String(),
	})
	if err != nil {
		return errors.Wrap(err, "error when parsing kube-proxy configmap template")
	}

	if err := createKubeDNSAddon(dnsDeploymentBytes, dnsServiceBytes, client); err != nil {
		return err
	}
	fmt.Println("[addons] Applied essential addon: kube-dns")
	return nil
}

// CreateServiceAccount creates the necessary serviceaccounts that kubeadm uses/might use, if they don't already exist.
func CreateServiceAccount(client clientset.Interface) error {

	return apiclient.CreateOrUpdateServiceAccount(client, &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KubeDNSServiceAccountName,
			Namespace: metav1.NamespaceSystem,
		},
	})
}

func createKubeDNSAddon(deploymentBytes, serviceBytes []byte, client clientset.Interface) error {
	kubednsDeployment := &apps.Deployment{}
	if err := kuberuntime.DecodeInto(clientsetscheme.Codecs.UniversalDecoder(), deploymentBytes, kubednsDeployment); err != nil {
		return errors.Wrap(err, "unable to decode kube-dns deployment")
	}

	// Create the Deployment for kube-dns or update it in case it already exists
	if err := apiclient.CreateOrUpdateDeployment(client, kubednsDeployment); err != nil {
		return err
	}

	kubednsService := &v1.Service{}
	return createDNSService(kubednsService, serviceBytes, client)
}
