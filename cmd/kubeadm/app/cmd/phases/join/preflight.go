package phases

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	"k8s.io/klog"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	"k8s.io/kubernetes/cmd/kubeadm/app/preflight"
	utilsexec "k8s.io/utils/exec"
)

var (
	preflightExample = cmdutil.Examples(`
		# Run join pre-flight checks using a config file.
		kubeadm join phase preflight --config kubeadm-config.yml
		`)

	notReadyToJoinControlPlaneTemp = template.Must(template.New("join").Parse(dedent.Dedent(`
		One or more conditions for hosting a new control plane instance is not satisfied.

		{{.Error}}

		Please ensure that:
		* The cluster has a stable controlPlaneEndpoint address.
		* The certificates that must be shared among control plane instances are provided.

		`)))
)

// NewPreflightPhase creates a kubeadm workflow phase that implements preflight checks for a new node join
func NewPreflightPhase() workflow.Phase {
	return workflow.Phase{
		Name:    "preflight [api-server-endpoint]",
		Short:   "Run join pre-flight checks",
		Long:    "Run pre-flight checks for kubeadm join.",
		Example: preflightExample,
		Run:     runPreflight,
		InheritFlags: []string{
			options.CfgPath,
			options.IgnorePreflightErrors,
			options.TLSBootstrapToken,
			options.TokenStr,
			options.ControlPlane,
			options.APIServerAdvertiseAddress,
			options.APIServerBindPort,
			options.NodeCRISocket,
			options.NodeName,
			options.FileDiscovery,
			options.TokenDiscovery,
			options.TokenDiscoveryCAHash,
			options.TokenDiscoverySkipCAHash,
			options.CertificateKey,
		},
	}
}

// runPreflight executes preflight checks logic.
func runPreflight(c workflow.RunData) (err error) {
	j, ok := c.(JoinData)
	if !ok {
		return errors.New("preflight phase invoked with an invalid data struct")
	}
	fmt.Println("[preflight] Running pre-flight checks")

	// Start with general checks
	klog.V(1).Infoln("[preflight] Running general checks")
	err = preflight.RunJoinNodeChecks(
		utilsexec.New(), j.Cfg(), j.IgnorePreflightErrors(),
	)
	if err != nil {
		return err
	}

	initCfg, err := j.InitCfg()
	if err != nil {
		return err
	}

	// Continue with more specific checks based on the init configuration
	klog.V(1).Infoln("[preflight] Running configuration dependant checks")
	if j.Cfg().ControlPlane != nil {
		// Checks if the cluster configuration supports
		// joining a new control plane instance and if all the necessary certificates are provided
		hasCertificateKey := len(j.CertificateKey()) > 0
		err = checkIfReadyForAdditionalControlPlane(
			&initCfg.ClusterConfiguration, hasCertificateKey,
		)
		if err != nil {
			// outputs the not ready for hosting a new control plane instance message
			ctx := map[string]string{
				"Error": err.Error(),
			}

			var msg bytes.Buffer
			notReadyToJoinControlPlaneTemp.Execute(&msg, ctx)
			return errors.New(msg.String())
		}

		// run kubeadm init preflight checks for checking all the prerequisites
		fmt.Println("[preflight] Running pre-flight checks before initializing the new control plane instance")
		err = preflight.RunInitNodeChecks(
			utilsexec.New(), initCfg, j.IgnorePreflightErrors(), true, hasCertificateKey,
		)
		if err != nil {
			return err
		}

		fmt.Println("[preflight] Pulling images required for setting up a Kubernetes cluster")
		fmt.Println("[preflight] This might take a minute or two, depending on the speed of your internet connection")
		fmt.Println("[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'")
		err = preflight.RunPullImagesCheck(
			utilsexec.New(), initCfg, j.IgnorePreflightErrors(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// 	@param hasCertificateKey: 如果 kubeadm join 添加新的 master 节点时, 
// 	没有指定额外的 crt/key 对, 那么就表示, 所有 master 节点都使用相同的 ca.crt/ca.key ...
// 	这要求部署人员事先将这些文件拷贝过去, 这是 precheck 的条件之一.
//
// checkIfReadyForAdditionalControlPlane ensures that the cluster is in a state that supports
// joining an additional control plane instance and if the node is ready to preflight
func checkIfReadyForAdditionalControlPlane(
	initConfiguration *kubeadmapi.ClusterConfiguration, hasCertificateKey bool,
) error {
	// blocks if the cluster was created without a stable control plane endpoint
	if initConfiguration.ControlPlaneEndpoint == "" {
		return errors.New(
			"unable to add a new control plane instance a cluster "+
			"that doesn't have a stable controlPlaneEndpoint address",
		)
	}

	if !hasCertificateKey {
		// checks if the certificates that must be equal across controlplane instances are provided
		if ret, err := certs.SharedCertificateExists(initConfiguration); !ret {
			return err
		}
	}

	return nil
}
