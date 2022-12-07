package phases

import (
	"github.com/pkg/errors"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/options"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	cmdutil "k8s.io/kubernetes/cmd/kubeadm/app/cmd/util"
	markcontrolplanephase "k8s.io/kubernetes/cmd/kubeadm/app/phases/markcontrolplane"
)

var (
	markControlPlaneExample = cmdutil.Examples(`
		# Applies control-plane label and taint to the current node, functionally equivalent to what executed by kubeadm init.
		kubeadm init phase mark-control-plane --config config.yml

		# Applies control-plane label and taint to a specific node
		kubeadm init phase mark-control-plane --node-name myNode
		`)
)


// MarkControlPlane 为 master 节点添加 label(标签) 与 taints(污点) 配置.
//
// 	由于 kubelet 会自动向 apiserver 注册自己,
// 	所以在调用本函数之前, node 应该已经存在于 apiserver 里了, 这里只是 patch 一下而已.
//
// NewMarkControlPlanePhase creates a kubeadm workflow phase
// that implements mark-controlplane checks.
func NewMarkControlPlanePhase() workflow.Phase {
	return workflow.Phase{
		Name:    "mark-control-plane",
		Short:   "Mark a node as a control-plane",
		Example: markControlPlaneExample,
		InheritFlags: []string{
			options.NodeName,
			options.CfgPath,
		},
		Run: runMarkControlPlane,
	}
}

// runMarkControlPlane executes mark-control-plane checks logic.
func runMarkControlPlane(c workflow.RunData) error {
	data, ok := c.(InitData)
	if !ok {
		return errors.New("mark-control-plane phase invoked with an invalid data struct")
	}

	client, err := data.Client()
	if err != nil {
		return err
	}

	nodeRegistration := data.Cfg().NodeRegistration
	return markcontrolplanephase.MarkControlPlane(
		client, nodeRegistration.Name, nodeRegistration.Taints,
	)
}
