// Package podtolerationrestriction is a plugin that first verifies
// any conflict between a pod's tolerations and its namespace's
// tolerations, and rejects the pod if there's a conflict.  If there's
// no conflict, the pod's tolerations are merged with its namespace's
// toleration. Resulting pod's tolerations are verified against its
// namespace's whitelist of tolerations. If the verification is
// successful, the pod is admitted otherwise rejected. If a namespace
// does not have associated default or whitelist of tolerations, then
// cluster level default or whitelist of tolerations are used instead
// if specified. Tolerations to a namespace are assigned via
// scheduler.alpha.kubernetes.io/defaultTolerations and
// scheduler.alpha.kubernetes.io/tolerationsWhitelist annotations
// keys.
package podtolerationrestriction // import "k8s.io/kubernetes/plugin/pkg/admission/podtolerationrestriction"
