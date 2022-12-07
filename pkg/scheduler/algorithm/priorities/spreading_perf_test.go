package priorities

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/listers/fake"
	nodeinfosnapshot "k8s.io/kubernetes/pkg/scheduler/nodeinfo/snapshot"
	st "k8s.io/kubernetes/pkg/scheduler/testing"
)

// The tests in this file compare the performance of SelectorSpreadPriority
// against EvenPodsSpreadPriority with a similar rule.

var (
	tests = []struct {
		name            string
		existingPodsNum int
		allNodesNum     int
	}{
		{
			name:            "100nodes",
			existingPodsNum: 1000,
			allNodesNum:     100,
		},
		{
			name:            "1000nodes",
			existingPodsNum: 10000,
			allNodesNum:     1000,
		},
	}
)

func BenchmarkTestDefaultEvenPodsSpreadPriority(b *testing.B) {
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			pod := st.MakePod().Name("p").Label("foo", "").
				SpreadConstraint(1, v1.LabelHostname, softSpread, st.MakeLabelSelector().Exists("foo").Obj()).
				SpreadConstraint(1, v1.LabelZoneFailureDomain, softSpread, st.MakeLabelSelector().Exists("foo").Obj()).Obj()
			existingPods, allNodes, filteredNodes := st.MakeNodesAndPodsForEvenPodsSpread(pod.Labels, tt.existingPodsNum, tt.allNodesNum, tt.allNodesNum)
			snapshot := nodeinfosnapshot.NewSnapshot(nodeinfosnapshot.CreateNodeInfoMap(existingPods, allNodes))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tpSpreadMap, err := buildPodTopologySpreadMap(pod, filteredNodes, snapshot.NodeInfoList)
				if err != nil {
					b.Fatal(err)
				}
				meta := &priorityMetadata{
					podTopologySpreadMap: tpSpreadMap,
				}
				var gotList framework.NodeScoreList
				for _, n := range filteredNodes {
					score, err := CalculateEvenPodsSpreadPriorityMap(pod, meta, snapshot.NodeInfoMap[n.Name])
					if err != nil {
						b.Fatal(err)
					}
					gotList = append(gotList, score)
				}
				err = CalculateEvenPodsSpreadPriorityReduce(pod, meta, snapshot, gotList)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkTestSelectorSpreadPriority(b *testing.B) {
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			pod := st.MakePod().Name("p").Label("foo", "").Obj()
			existingPods, allNodes, filteredNodes := st.MakeNodesAndPodsForEvenPodsSpread(pod.Labels, tt.existingPodsNum, tt.allNodesNum, tt.allNodesNum)
			snapshot := nodeinfosnapshot.NewSnapshot(nodeinfosnapshot.CreateNodeInfoMap(existingPods, allNodes))
			services := []*v1.Service{{Spec: v1.ServiceSpec{Selector: map[string]string{"foo": ""}}}}
			ss := SelectorSpread{
				serviceLister:     fake.ServiceLister(services),
				controllerLister:  fake.ControllerLister(nil),
				replicaSetLister:  fake.ReplicaSetLister(nil),
				statefulSetLister: fake.StatefulSetLister(nil),
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				meta := &priorityMetadata{
					podSelector: getSelector(pod, ss.serviceLister, ss.controllerLister, ss.replicaSetLister, ss.statefulSetLister),
				}
				_, err := runMapReducePriority(ss.CalculateSpreadPriorityMap, ss.CalculateSpreadPriorityReduce, meta, pod, snapshot, filteredNodes)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
