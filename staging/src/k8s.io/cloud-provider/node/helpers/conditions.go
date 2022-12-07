package helpers

import (
	"encoding/json"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
)

// GetNodeCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func GetNodeCondition(status *v1.NodeStatus, conditionType v1.NodeConditionType) (int, *v1.NodeCondition) {
	if status == nil {
		return -1, nil
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return i, &status.Conditions[i]
		}
	}
	return -1, nil
}

// SetNodeCondition updates specific node condition with patch operation.
func SetNodeCondition(c clientset.Interface, node types.NodeName, condition v1.NodeCondition) error {
	condition.LastHeartbeatTime = metav1.NewTime(time.Now())
	patch, err := json.Marshal(map[string]interface{}{
		"status": map[string]interface{}{
			"conditions": []v1.NodeCondition{condition},
		},
	})
	if err != nil {
		return err
	}
	_, err = c.CoreV1().Nodes().PatchStatus(string(node), patch)
	return err
}
