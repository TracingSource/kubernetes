package patchnode

import (
	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
)

// AnnotateCRISocket annotates the node with the given crisocket
func AnnotateCRISocket(client clientset.Interface, nodeName string, criSocket string) error {

	klog.V(1).Infof("[patchnode] Uploading the CRI Socket information %q to the Node API object %q as an annotation\n", criSocket, nodeName)

	return apiclient.PatchNode(client, nodeName, func(n *v1.Node) {
		annotateNodeWithCRISocket(n, criSocket)
	})
}

func annotateNodeWithCRISocket(n *v1.Node, criSocket string) {
	if n.ObjectMeta.Annotations == nil {
		n.ObjectMeta.Annotations = make(map[string]string)
	}
	n.ObjectMeta.Annotations[constants.AnnotationKubeadmCRISocket] = criSocket
}
