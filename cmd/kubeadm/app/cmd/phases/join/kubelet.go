package phases

import (
	"fmt"
	"os"

	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/apimachinery/pkg/util/wait"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeletphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/kubelet"
	patchnodephase "k8s.io/kubernetes/cmd/kubeadm/app/phases/patchnode"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
)

var (
	kubeadmJoinFailMsg = dedent.Dedent(`
		Unfortunately, an error has occurred:
			%v

		This error is likely caused by:
			- The kubelet is not running
			- The kubelet is unhealthy due to a misconfiguration of the node in some way (required cgroups disabled)

		If you are on a systemd-powered system, you can try to troubleshoot the error with the following commands:
			- 'systemctl status kubelet'
			- 'journalctl -xeu kubelet'
		`)
)

// NewKubeletStartPhase 启动新主机上的 kubelet 服务
// (kubelet 会自动向 apiserver 注册自身所在的 Node 信息).
//
// NewKubeletStartPhase creates a kubeadm workflow phase that start kubelet on a node.
func NewKubeletStartPhase() workflow.Phase {
	return workflow.Phase{
		Name:  "kubelet-start [api-server-endpoint]",
		Short: "Write kubelet settings, certificates and (re)start the kubelet",
		Long:  "Write a file with KubeletConfiguration and an environment file with node specific kubelet settings, and then (re)start kubelet.",
		Run:   runKubeletStartJoinPhase,
		InheritFlags: []string{
			options.CfgPath,
			options.NodeCRISocket,
			options.NodeName,
			options.FileDiscovery,
			options.TokenDiscovery,
			options.TokenDiscoveryCAHash,
			options.TokenDiscoverySkipCAHash,
			options.TLSBootstrapToken,
			options.TokenStr,
		},
	}
}

func getKubeletStartJoinData(c workflow.RunData) (*kubeadmapi.JoinConfiguration, *kubeadmapi.InitConfiguration, *clientcmdapi.Config, error) {
	data, ok := c.(JoinData)
	if !ok {
		return nil, nil, nil, errors.New("kubelet-start phase invoked with an invalid data struct")
	}
	cfg := data.Cfg()
	initCfg, err := data.InitCfg()
	if err != nil {
		return nil, nil, nil, err
	}
	tlsBootstrapCfg, err := data.TLSBootstrapCfg()
	if err != nil {
		return nil, nil, nil, err
	}
	return cfg, initCfg, tlsBootstrapCfg, nil
}

// runKubeletStartJoinPhase executes the kubelet TLS bootstrap process.
// This process is executed by the kubelet and completes with the node joining the cluster
// with a dedicates set of credentials as required by the node authorizer
func runKubeletStartJoinPhase(c workflow.RunData) (err error) {
	cfg, initCfg, tlsBootstrapCfg, err := getKubeletStartJoinData(c)
	if err != nil {
		return err
	}
	// /etc/kubernetes/bootstrap-kubelet.conf
	bootstrapKubeConfigFile := kubeadmconstants.GetBootstrapKubeletKubeConfigPath()
	// Deletes the bootstrapKubeConfigFile,
	// so the credential used for TLS bootstrap is removed from disk
	defer os.Remove(bootstrapKubeConfigFile)

	// Write the bootstrap kubelet config file or the TLS-Bootstrapped kubelet config file down to disk
	klog.V(1).Infof(
		"[kubelet-start] writing bootstrap kubelet config file at %s",
		bootstrapKubeConfigFile,
	)
	// 将 tlsBootstrapCfg 的内容, 写入到 /etc/kubernetes/bootstrap-kubelet.conf 文件.
	// 一般来说, 这个文件的内容等同于 /etc/kubernetes/admin.conf
	//
	// 这是一个临时文件, 用完就会删除.
	err = kubeconfigutil.WriteToDisk(bootstrapKubeConfigFile, tlsBootstrapCfg)
	if err != nil {
		return errors.Wrap(err, "couldn't save bootstrap-kubelet.conf to disk")
	}

	// Write the ca certificate to disk so kubelet can use it for authentication
	cluster := tlsBootstrapCfg.Contexts[tlsBootstrapCfg.CurrentContext].Cluster
	if _, err := os.Stat(cfg.CACertPath); os.IsNotExist(err) {
		klog.V(1).Infof("[kubelet-start] writing CA certificate at %s", cfg.CACertPath)
		err := certutil.WriteCert(
			cfg.CACertPath, tlsBootstrapCfg.Clusters[cluster].CertificateAuthorityData,
		)
		if err != nil {
			return errors.Wrap(err, "couldn't save the CA certificate to disk")
		}
	}

	kubeletVersion, err := version.ParseSemantic(
		initCfg.ClusterConfiguration.KubernetesVersion,
	)
	if err != nil {
		return err
	}

	bootstrapClient, err := kubeconfigutil.ClientSetFromFile(bootstrapKubeConfigFile)
	if err != nil {
		return errors.Errorf(
			"couldn't create client from kubeconfig file %q", bootstrapKubeConfigFile,
		)
	}
	// 此时, kubelet 的配置文件, 数据目录等都还不存在, 这里要开始创建了, 首先停止 kubelet.
	//
	// Configure the kubelet.
	// In this short timeframe, kubeadm is trying to stop/restart the kubelet
	// Try to stop the kubelet service so no race conditions occur when configuring it
	klog.V(1).Infoln("[kubelet-start] Stopping the kubelet")
	kubeletphase.TryStopKubelet()

	// 从初始 master 主节点中, 获取 kube-system/kubelet-config-1.17 的 ConfigMap 对象.
	// 然后写入到 /var/lib/kubelet/config.yaml 文件中.
	//
	// Write the configuration for the kubelet (using the bootstrap token credentials)
	// to disk so the kubelet can start
	err = kubeletphase.DownloadConfig(
		bootstrapClient, kubeletVersion, kubeadmconstants.KubeletRunDirectory,
	)
	if err != nil {
		return err
	}

	// 构造 kubelet 启动参数, 并写入 /var/lib/kubelet/kubeadm-flags.env
	//
	// Write env file with flags for the kubelet to use.
	// We only want to register the joining node with the specified taints
	// if the node is not a control-plane.
	// The mark-control-plane phase will register the taints otherwise.
	registerTaintsUsingFlags := cfg.ControlPlane == nil
	err = kubeletphase.WriteKubeletDynamicEnvFile(
		&initCfg.ClusterConfiguration, &initCfg.NodeRegistration,
		registerTaintsUsingFlags, kubeadmconstants.KubeletRunDirectory,
	)
	if err != nil {
		return err
	}

	// 启动 kubelet, 直到这里, 才由 kubelet 向 apiserver 注册 Node 对象.
	// (kubelet get nodes 才会看到新添加的主机)
	//
	// Try to start the kubelet service in case it's inactive
	fmt.Println("[kubelet-start] Starting the kubelet")
	kubeletphase.TryStartKubelet()

	// 再次等待 40s, 等待 kubelet 将 3 大件启动完成.

	// 等待生成 /etc/kubernetes/kubelet.conf 文件
	//
	// Now the kubelet will perform the TLS Bootstrap, transforming
	// /etc/kubernetes/bootstrap-kubelet.conf to /etc/kubernetes/kubelet.conf
	// Wait for the kubelet to create the /etc/kubernetes/kubelet.conf kubeconfig file.
	// If this process times out, display a somewhat user-friendly message.
	waiter := apiclient.NewKubeWaiter(nil, kubeadmconstants.TLSBootstrapTimeout, os.Stdout)
	if err := waiter.WaitForKubeletAndFunc(waitForTLSBootstrappedClient); err != nil {
		fmt.Printf(kubeadmJoinFailMsg, err)
		return err
	}

	// 根据 /etc/kubernetes/kubelet.conf 文件构建 kube client
	//
	// When we know the /etc/kubernetes/kubelet.conf file is available, get the client
	client, err := kubeconfigutil.ClientSetFromFile(kubeadmconstants.GetKubeletKubeConfigPath())
	if err != nil {
		return err
	}

	// Patch 一下当前的 Node 作为验证
	klog.V(1).Infoln("[kubelet-start] preserving the crisocket information for the node")
	err = patchnodephase.AnnotateCRISocket(
		client, cfg.NodeRegistration.Name, cfg.NodeRegistration.CRISocket,
	)
	if err != nil {
		return errors.Wrap(err, "error uploading crisocket")
	}

	return nil
}

// waitForTLSBootstrappedClient 等待 kubelet 使用 /etc/kubernetes/bootstrap-kubelet.conf 建立与 apiserver 的受限连接,
// 并生成 /etc/kubernetes/kubelet.conf 文件.
//
// waitForTLSBootstrappedClient waits for the /etc/kubernetes/kubelet.conf file to be available
func waitForTLSBootstrappedClient() error {
	fmt.Println("[kubelet-start] Waiting for the kubelet to perform the TLS Bootstrap...")

	// Loop on every falsy return. Return with an error if raised. Exit successfully if true is returned.
	return wait.PollImmediate(kubeadmconstants.APICallRetryInterval, kubeadmconstants.TLSBootstrapTimeout, func() (bool, error) {
		// Check that we can create a client set out of the kubelet kubeconfig. This ensures not
		// only that the kubeconfig file exists, but that other files required by it also exist (like
		// client certificate and key)
		_, err := kubeconfigutil.ClientSetFromFile(kubeadmconstants.GetKubeletKubeConfigPath())
		return (err == nil), nil
	})
}
