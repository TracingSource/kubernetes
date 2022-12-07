package nodename

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

func TestNodeName(t *testing.T) {
	tests := []struct {
		pod        *v1.Pod
		node       *v1.Node
		name       string
		wantStatus *framework.Status
	}{
		{
			pod:  &v1.Pod{},
			node: &v1.Node{},
			name: "no host specified",
		},
		{
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					NodeName: "foo",
				},
			},
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
			name: "host matches",
		},
		{
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					NodeName: "bar",
				},
			},
			node: &v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
			name:       "host doesn't match",
			wantStatus: framework.NewStatus(framework.UnschedulableAndUnresolvable, predicates.ErrPodNotMatchHostName.GetReason()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodeInfo := schedulernodeinfo.NewNodeInfo()
			nodeInfo.SetNode(test.node)

			p, _ := New(nil, nil)
			gotStatus := p.(framework.FilterPlugin).Filter(context.Background(), nil, test.pod, nodeInfo)
			if !reflect.DeepEqual(gotStatus, test.wantStatus) {
				t.Errorf("status does not match: %v, want: %v", gotStatus, test.wantStatus)
			}
		})
	}
}
