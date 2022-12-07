package equal

import apiv1 "k8s.io/api/core/v1"

// KubeletConfigOkEq returns true if the two conditions are semantically equivalent in the context of dynamic config
func KubeletConfigOkEq(a, b *apiv1.NodeCondition) bool {
	return a.Message == b.Message && a.Reason == b.Reason && a.Status == b.Status
}
