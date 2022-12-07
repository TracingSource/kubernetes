package kubelet

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodResourceLimitsDefaulting(t *testing.T) {
	tk := newTestKubelet(t, true)
	defer tk.Cleanup()
	tk.kubelet.nodeLister = &testNodeLister{
		nodes: []*v1.Node{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: string(tk.kubelet.nodeName),
				},
				Status: v1.NodeStatus{
					Allocatable: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("6"),
						v1.ResourceMemory: resource.MustParse("4Gi"),
					},
				},
			},
		},
	}
	cases := []struct {
		pod      *v1.Pod
		expected *v1.Pod
	}{
		{
			pod:      getPod("0", "0"),
			expected: getPod("6", "4Gi"),
		},
		{
			pod:      getPod("1", "0"),
			expected: getPod("1", "4Gi"),
		},
		{
			pod:      getPod("", ""),
			expected: getPod("6", "4Gi"),
		},
		{
			pod:      getPod("0", "1Mi"),
			expected: getPod("6", "1Mi"),
		},
	}
	as := assert.New(t)
	for idx, tc := range cases {
		actual, _, err := tk.kubelet.defaultPodLimitsForDownwardAPI(tc.pod, nil)
		as.Nil(err, "failed to default pod limits: %v", err)
		if !apiequality.Semantic.DeepEqual(tc.expected, actual) {
			as.Fail("test case [%d] failed.  Expected: %+v, Got: %+v", idx, tc.expected, actual)
		}
	}
}

func getPod(cpuLimit, memoryLimit string) *v1.Pod {
	resources := v1.ResourceRequirements{}
	if cpuLimit != "" || memoryLimit != "" {
		resources.Limits = make(v1.ResourceList)
	}
	if cpuLimit != "" {
		resources.Limits[v1.ResourceCPU] = resource.MustParse(cpuLimit)
	}
	if memoryLimit != "" {
		resources.Limits[v1.ResourceMemory] = resource.MustParse(memoryLimit)
	}
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:      "foo",
					Resources: resources,
				},
			},
		},
	}
}
