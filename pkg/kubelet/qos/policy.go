package qos

import (
	v1 "k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	// PodInfraOOMAdj is very docker specific.
	// For arbitrary runtime, it may not make sense to set sandbox level oom score,
	// e.g. a sandbox could only be a namespace without a process.
	// TODO: Handle infra container oom score adj in a runtime agnostic way.
	PodInfraOOMAdj int = -998
	// KubeletOOMScoreAdj is the OOM score adjustment for Kubelet
	KubeletOOMScoreAdj int = -999
	// DockerOOMScoreAdj is the OOM score adjustment for Docker
	DockerOOMScoreAdj int = -999
	// KubeProxyOOMScoreAdj is the OOM score adjustment for kube-proxy
	KubeProxyOOMScoreAdj  int = -999
	guaranteedOOMScoreAdj int = -998
	besteffortOOMScoreAdj int = 1000
)

// GetContainerOOMScoreAdjust 根据目标 pod 的 requests/limits 值判断其 QoS 类型,
// 并返回相应的 oom_score_adj 值.
//
// 	@param memoryCapacity: 宿主机的可用内存.
//
// 	@return oomScoreAdjust: OS发生oom时 kill 进程的优先级, 分值越大越容易被 kill.
//
// GetContainerOOMScoreAdjust returns the amount by which the OOM score of all
// processes in the container should be adjusted.
// The OOM score of a process is the percentage of memory it consumes multiplied
// by 10 (barring exceptional cases) + a configurable quantity
// which is between -1000 and 1000.
// Containers with higher OOM scores are killed if the system runs out of memory.
// See https://lwn.net/Articles/391222/ for more information.
func GetContainerOOMScoreAdjust(
	pod *v1.Pod, container *v1.Container, memoryCapacity int64,
) int {
	if types.IsCriticalPod(pod) {
		// Critical pods should be the last to get killed.
		return guaranteedOOMScoreAdj
	}

	switch v1qos.GetPodQOS(pod) {
	case v1.PodQOSGuaranteed:
		// Guaranteed containers should be the last to get killed.
		return guaranteedOOMScoreAdj
	case v1.PodQOSBestEffort:
		return besteffortOOMScoreAdj
	}

	// Burstable containers are a middle tier, between Guaranteed and Best-Effort.
	// Ideally, we want to protect Burstable containers that consume less memory
	// than requested.
	// The formula below is a heuristic(启发式的, 试探性的).
	// A container requesting for 10% of a system's memory will have
	// an OOM score adjust of 900.
	// If a process in container Y uses over 10% of memory, its OOM score will be 1000.
	// The idea is that containers which use more than their request will have
	// an OOM score of 1000 and will be prime targets for OOM kills.
	//
	// Note that this is a heuristic, it won't work if a container has many
	// small processes.
	memoryRequest := container.Resources.Requests.Memory().Value()
	// 如下公式可简化成: 1000 * (1 - memoryRequest/memoryCapacity), 取值在 (0, 1000)
	// 视容器 reqeusts 的内存占比而定, 占比越高的越**不容易**被 kill.
	//
	// 但这只是 k8s 的角度, 对于 OS 来说, 在发生系统级 OOM 时, 如果优先级一致,
	// 先干掉的肯定是占用资源多的(只 kill 小内存进程的话, 要 kill 多少个才够啊).
	//
	oomScoreAdjust := 1000 - (1000*memoryRequest)/memoryCapacity
	// (保底承诺)保证 guaranteed 类型的 Pod 比 burstable 类型的存活机率大(两点).
	// A guaranteed pod using 100% of memory can have an OOM score of 10.
	// Ensure that burstable pods have a higher OOM score adjustment.
	if int(oomScoreAdjust) < (1000 + guaranteedOOMScoreAdj) {
		return (1000 + guaranteedOOMScoreAdj)
	}
	// (保底承诺)保证 burstable 类型的 Pod 比 besteffort 类型的存活机率大(一点).
	//
	// Give burstable pods a higher chance of survival over besteffort pods.
	if int(oomScoreAdjust) == besteffortOOMScoreAdj {
		return int(oomScoreAdjust - 1)
	}
	return int(oomScoreAdjust)
}
