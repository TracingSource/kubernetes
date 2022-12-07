package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	clientretry "k8s.io/client-go/util/retry"
	"k8s.io/klog"
)

var updateLabelBackoff = wait.Backoff{
	Steps:    5,
	Duration: 100 * time.Millisecond,
	Jitter:   1.0,
}

// AddOrUpdateLabelsOnNode updates the labels on the node and returns true on
// success and false on failure.
func AddOrUpdateLabelsOnNode(kubeClient clientset.Interface, labelsToUpdate map[string]string, node *v1.Node) bool {
	err := addOrUpdateLabelsOnNode(kubeClient, node.Name, labelsToUpdate)
	if err != nil {
		utilruntime.HandleError(
			fmt.Errorf(
				"unable to update labels %+v for Node %q: %v",
				labelsToUpdate,
				node.Name,
				err))
		return false
	}

	klog.V(4).Infof("Updated labels %+v to Node %v", labelsToUpdate, node.Name)
	return true
}

func addOrUpdateLabelsOnNode(kubeClient clientset.Interface, nodeName string, labelsToUpdate map[string]string) error {
	firstTry := true
	return clientretry.RetryOnConflict(updateLabelBackoff, func() error {
		var err error
		var node *v1.Node
		// First we try getting node from the API server cache, as it's cheaper. If it fails
		// we get it from etcd to be sure to have fresh data.
		if firstTry {
			node, err = kubeClient.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{ResourceVersion: "0"})
			firstTry = false
		} else {
			node, err = kubeClient.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		}
		if err != nil {
			return err
		}

		// Make a copy of the node and update the labels.
		newNode := node.DeepCopy()
		if newNode.Labels == nil {
			newNode.Labels = make(map[string]string)
		}
		for key, value := range labelsToUpdate {
			newNode.Labels[key] = value
		}

		oldData, err := json.Marshal(node)
		if err != nil {
			return fmt.Errorf("failed to marshal the existing node %#v: %v", node, err)
		}
		newData, err := json.Marshal(newNode)
		if err != nil {
			return fmt.Errorf("failed to marshal the new node %#v: %v", newNode, err)
		}
		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, &v1.Node{})
		if err != nil {
			return fmt.Errorf("failed to create a two-way merge patch: %v", err)
		}
		if _, err := kubeClient.CoreV1().Nodes().Patch(node.Name, types.StrategicMergePatchType, patchBytes); err != nil {
			return fmt.Errorf("failed to patch the node: %v", err)
		}
		return nil
	})
}
