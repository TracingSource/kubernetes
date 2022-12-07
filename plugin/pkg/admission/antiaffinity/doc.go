// Package antiaffinity provides the LimitPodHardAntiAffinityTopology
// admission controller. It rejects any pod that specifies "hard"
// (RequiredDuringScheduling) anti-affinity with a TopologyKey other
// than v1.LabelHostname.  Because anti-affinity is symmetric, without
// this admission controller, a user could maliciously or accidentally
// specify that their pod (once it has scheduled) should block other
// pods from scheduling into the same zone or some other large
// topology, essentially DoSing the cluster.  In the future we will
// address this problem more fully by using quota and priority, but
// for now this admission controller provides a simple protection, on
// the assumption that the only legitimate use of hard pod
// anti-affinity is to exclude other pods from the same node.
package antiaffinity // import "k8s.io/kubernetes/plugin/pkg/admission/antiaffinity"
