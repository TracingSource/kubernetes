// +build windows

package stats

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/cm"
)

func (sp *summaryProviderImpl) GetSystemContainersStats(nodeConfig cm.NodeConfig, podStats []statsapi.PodStats, updateStats bool) (stats []statsapi.ContainerStats) {
	stats = append(stats, sp.getSystemPodsCPUAndMemoryStats(nodeConfig, podStats, updateStats))
	return stats
}

func (sp *summaryProviderImpl) GetSystemContainersCPUAndMemoryStats(nodeConfig cm.NodeConfig, podStats []statsapi.PodStats, updateStats bool) (stats []statsapi.ContainerStats) {
	stats = append(stats, sp.getSystemPodsCPUAndMemoryStats(nodeConfig, podStats, updateStats))
	return stats
}

func (sp *summaryProviderImpl) getSystemPodsCPUAndMemoryStats(nodeConfig cm.NodeConfig, podStats []statsapi.PodStats, updateStats bool) statsapi.ContainerStats {
	now := metav1.NewTime(time.Now())
	podsSummary := statsapi.ContainerStats{
		StartTime: now,
		CPU:       &statsapi.CPUStats{},
		Memory:    &statsapi.MemoryStats{},
		Name:      statsapi.SystemContainerPods,
	}

	// Sum up all pod's stats.
	var usageCoreNanoSeconds uint64
	var usageNanoCores uint64
	var availableBytes uint64
	var usageBytes uint64
	var workingSetBytes uint64
	for _, pod := range podStats {
		if pod.CPU != nil {
			podsSummary.CPU.Time = now
			if pod.CPU.UsageCoreNanoSeconds != nil {
				usageCoreNanoSeconds = usageCoreNanoSeconds + *pod.CPU.UsageCoreNanoSeconds
			}
			if pod.CPU.UsageNanoCores != nil {
				usageNanoCores = usageNanoCores + *pod.CPU.UsageNanoCores
			}
		}

		if pod.Memory != nil {
			podsSummary.Memory.Time = now
			if pod.Memory.AvailableBytes != nil {
				availableBytes = availableBytes + *pod.Memory.AvailableBytes
			}
			if pod.Memory.UsageBytes != nil {
				usageBytes = usageBytes + *pod.Memory.UsageBytes
			}
			if pod.Memory.WorkingSetBytes != nil {
				workingSetBytes = workingSetBytes + *pod.Memory.WorkingSetBytes
			}
		}
	}

	// Set results only if they are not zero.
	if usageCoreNanoSeconds != 0 {
		podsSummary.CPU.UsageCoreNanoSeconds = &usageCoreNanoSeconds
	}
	if usageNanoCores != 0 {
		podsSummary.CPU.UsageNanoCores = &usageNanoCores
	}
	if availableBytes != 0 {
		podsSummary.Memory.AvailableBytes = &availableBytes
	}
	if usageBytes != 0 {
		podsSummary.Memory.UsageBytes = &usageBytes
	}
	if workingSetBytes != 0 {
		podsSummary.Memory.WorkingSetBytes = &workingSetBytes
	}

	return podsSummary
}
