package markcontrolplane

import (
	"fmt"

	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
)

// MarkControlPlane 为 master 节点添加 label(标签) 与 taints(污点) 配置.
//
// 	由于 kubelet 会自动向 apiserver 注册自己,
// 	所以在调用本函数之前, node 应该已经存在于 apiserver 里了, 这里只是 patch 一下而已.
//
// caller:
// 	1. cmd/kubeadm/app/cmd/phases/init/markcontrolplane.go -> runMarkControlPlane()
// 	kubeadm init 过程中被调用(3大件已启动).
//
// MarkControlPlane taints the control-plane and sets the control-plane label
func MarkControlPlane(
	client clientset.Interface, controlPlaneName string, taints []v1.Taint,
) error {

	fmt.Printf(
		"[mark-control-plane] Marking the node %s as control-plane by adding the label \"%s=''\"\n", 
		controlPlaneName, constants.LabelNodeRoleMaster,
	)

	if len(taints) > 0 {
		taintStrs := []string{}
		for _, taint := range taints {
			taintStrs = append(taintStrs, taint.ToString())
		}
		fmt.Printf(
			"[mark-control-plane] Marking the node %s as control-plane by adding the taints %v\n", 
			controlPlaneName, taintStrs,
		)
	}

	return apiclient.PatchNode(client, controlPlaneName, func(n *v1.Node) {
		markControlPlaneNode(n, taints)
	})
}

func taintExists(taint v1.Taint, taints []v1.Taint) bool {
	for _, t := range taints {
		if t == taint {
			return true
		}
	}

	return false
}

func markControlPlaneNode(n *v1.Node, taints []v1.Taint) {
	n.ObjectMeta.Labels[constants.LabelNodeRoleMaster] = ""

	for _, nt := range n.Spec.Taints {
		if !taintExists(nt, taints) {
			taints = append(taints, nt)
		}
	}

	n.Spec.Taints = taints
}
