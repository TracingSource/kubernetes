package qos

import (
	"strconv"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/kubernetes/pkg/apis/scheduling"
)

const (
	standardMemoryAmount = 8000000000
)

var (
	cpuLimit = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceName(v1.ResourceCPU): resource.MustParse("10"),
						},
					},
				},
			},
		},
	}

	memoryLimitCPURequest = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceName(v1.ResourceCPU): resource.MustParse("0"),
						},
						Limits: v1.ResourceList{
							v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
						},
					},
				},
			},
		},
	}

	zeroMemoryLimit = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							v1.ResourceName(v1.ResourceMemory): resource.MustParse("0"),
						},
					},
				},
			},
		},
	}

	noRequestLimit = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{},
				},
			},
		},
	}

	equalRequestLimitCPUMemory = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
							v1.ResourceName(v1.ResourceCPU):    resource.MustParse("5m"),
						},
						Limits: v1.ResourceList{
							v1.ResourceName(v1.ResourceCPU):    resource.MustParse("5m"),
							v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
						},
					},
				},
			},
		},
	}

	cpuUnlimitedMemoryLimitedWithRequests = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceName(v1.ResourceMemory): resource.MustParse(strconv.FormatInt(standardMemoryAmount/2, 10)),
							v1.ResourceName(v1.ResourceCPU):    resource.MustParse("5m"),
						},
						Limits: v1.ResourceList{
							v1.ResourceName(v1.ResourceMemory): resource.MustParse("10G"),
						},
					},
				},
			},
		},
	}

	requestNoLimit = v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceName(v1.ResourceMemory): resource.MustParse(strconv.FormatInt(standardMemoryAmount-1, 10)),
							v1.ResourceName(v1.ResourceCPU):    resource.MustParse("5m"),
						},
					},
				},
			},
		},
	}

	systemCritical = scheduling.SystemCriticalPriority

	critical = v1.Pod{
		Spec: v1.PodSpec{
			Priority: &systemCritical,
			Containers: []v1.Container{
				{
					Resources: v1.ResourceRequirements{},
				},
			},
		},
	}
)

type oomTest struct {
	pod             *v1.Pod
	memoryCapacity  int64
	lowOOMScoreAdj  int // The max oom_score_adj score the container should be assigned.
	highOOMScoreAdj int // The min oom_score_adj score the container should be assigned.
}

func TestGetContainerOOMScoreAdjust(t *testing.T) {
	oomTests := []oomTest{
		{
			pod:             &cpuLimit,
			memoryCapacity:  4000000000,
			lowOOMScoreAdj:  999,
			highOOMScoreAdj: 999,
		},
		{
			pod:             &memoryLimitCPURequest,
			memoryCapacity:  8000000000,
			lowOOMScoreAdj:  999,
			highOOMScoreAdj: 999,
		},
		{
			pod:             &zeroMemoryLimit,
			memoryCapacity:  7230457451,
			lowOOMScoreAdj:  1000,
			highOOMScoreAdj: 1000,
		},
		{
			pod:             &noRequestLimit,
			memoryCapacity:  4000000000,
			lowOOMScoreAdj:  1000,
			highOOMScoreAdj: 1000,
		},
		{
			pod:             &equalRequestLimitCPUMemory,
			memoryCapacity:  123456789,
			lowOOMScoreAdj:  -998,
			highOOMScoreAdj: -998,
		},
		{
			pod:             &cpuUnlimitedMemoryLimitedWithRequests,
			memoryCapacity:  standardMemoryAmount,
			lowOOMScoreAdj:  495,
			highOOMScoreAdj: 505,
		},
		{
			pod:             &requestNoLimit,
			memoryCapacity:  standardMemoryAmount,
			lowOOMScoreAdj:  2,
			highOOMScoreAdj: 2,
		},
		{
			pod:             &critical,
			memoryCapacity:  4000000000,
			lowOOMScoreAdj:  -998,
			highOOMScoreAdj: -998,
		},
	}
	for _, test := range oomTests {
		oomScoreAdj := GetContainerOOMScoreAdjust(test.pod, &test.pod.Spec.Containers[0], test.memoryCapacity)
		if oomScoreAdj < test.lowOOMScoreAdj || oomScoreAdj > test.highOOMScoreAdj {
			t.Errorf("oom_score_adj should be between %d and %d, but was %d", test.lowOOMScoreAdj, test.highOOMScoreAdj, oomScoreAdj)
		}
	}
}
