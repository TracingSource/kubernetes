package testing

import (
	"k8s.io/api/core/v1"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// DeepEqualSafePodSpec returns a PodSpec which is ready to be used with apiequality.Semantic.DeepEqual
func DeepEqualSafePodSpec() api.PodSpec {
	grace := int64(30)
	return api.PodSpec{
		RestartPolicy:                 api.RestartPolicyAlways,
		DNSPolicy:                     api.DNSClusterFirst,
		TerminationGracePeriodSeconds: &grace,
		SecurityContext:               &api.PodSecurityContext{},
		SchedulerName:                 api.DefaultSchedulerName,
	}
}

// V1DeepEqualSafePodSpec returns a PodSpec which is ready to be used with apiequality.Semantic.DeepEqual
func V1DeepEqualSafePodSpec() v1.PodSpec {
	grace := int64(30)
	return v1.PodSpec{
		RestartPolicy:                 v1.RestartPolicyAlways,
		DNSPolicy:                     v1.DNSClusterFirst,
		TerminationGracePeriodSeconds: &grace,
		SecurityContext:               &v1.PodSecurityContext{},
	}
}
