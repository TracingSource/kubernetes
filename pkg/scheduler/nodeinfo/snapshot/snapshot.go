package snapshot

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	schedulerlisters "k8s.io/kubernetes/pkg/scheduler/listers"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// 	@implementOf: pkg/scheduler/listers/listers.go -> SharedLister
//
// Snapshot is a snapshot of cache NodeInfo and NodeTree order.
// The scheduler takes a snapshot at the beginning of each scheduling cycle
// and uses it for its operations in that cycle.
type Snapshot struct {
	// key 为 node 节点名称, value 为该 node 的详细信息, 包含 pod 列表, 镜像概要等.
	//
	// NodeInfoMap a map of node name to a snapshot of its NodeInfo.
	NodeInfoMap map[string]*schedulernodeinfo.NodeInfo
	// NodeInfoMap 成员中 value 组成的列表
	//
	// NodeInfoList is the list of nodes as ordered in the cache's nodeTree.
	NodeInfoList []*schedulernodeinfo.NodeInfo
	// NodeInfoList 中, 存在 podAffinity 的部分
	//
	// HavePodsWithAffinityNodeInfoList is the list of nodes with at least one
	// pod declaring affinity terms.
	HavePodsWithAffinityNodeInfoList []*schedulernodeinfo.NodeInfo
	Generation                       int64
}

var _ schedulerlisters.SharedLister = &Snapshot{}

// NewEmptySnapshot initializes a Snapshot struct and returns it.
func NewEmptySnapshot() *Snapshot {
	return &Snapshot{
		NodeInfoMap: make(map[string]*schedulernodeinfo.NodeInfo),
	}
}

// 除了 _test.go 测试文件, 没有其他调用者.
//
// NewSnapshot initializes a Snapshot struct and returns it.
func NewSnapshot(nodeInfoMap map[string]*schedulernodeinfo.NodeInfo) *Snapshot {
	nodeInfoList := make([]*schedulernodeinfo.NodeInfo, 0, len(nodeInfoMap))
	havePodsWithAffinityNodeInfoList := make([]*schedulernodeinfo.NodeInfo, 0, len(nodeInfoMap))
	for _, v := range nodeInfoMap {
		nodeInfoList = append(nodeInfoList, v)
		if len(v.PodsWithAffinity()) > 0 {
			havePodsWithAffinityNodeInfoList = append(havePodsWithAffinityNodeInfoList, v)
		}
	}

	s := NewEmptySnapshot()
	s.NodeInfoMap = nodeInfoMap
	s.NodeInfoList = nodeInfoList
	s.HavePodsWithAffinityNodeInfoList = havePodsWithAffinityNodeInfoList

	return s
}

//
//
// 	@return: key 为 node 节点名称, value 为该 node 的详细信息, 包含 pod 列表, 镜像概要等.
//
// CreateNodeInfoMap obtains a list of pods and pivots that list into a map
// where the keys are node names and the values are the aggregated information
// for that node.
func CreateNodeInfoMap(
	pods []*v1.Pod, nodes []*v1.Node,
) map[string]*schedulernodeinfo.NodeInfo {
	// 梳理 pods 与所属 node 的关联关系, 并以 .spec.nodeName 为 key, pod 本身为 value.
	nodeNameToInfo := make(map[string]*schedulernodeinfo.NodeInfo)
	for _, pod := range pods {
		nodeName := pod.Spec.NodeName
		// value 是结构体对象, key 不存在时需要先创建, 已存在则新增即可.
		if _, ok := nodeNameToInfo[nodeName]; !ok {
			nodeNameToInfo[nodeName] = schedulernodeinfo.NewNodeInfo()
		}
		nodeNameToInfo[nodeName].AddPod(pod)
	}
	// 梳理目标 nodes 列表上的所有镜像, 以镜像名称为 key 构建 map, value 是存在该镜像的
	// Node 的 nodeName 集合.
	imageExistenceMap := createImageExistenceMap(nodes)

	for _, node := range nodes {
		if _, ok := nodeNameToInfo[node.Name]; !ok {
			nodeNameToInfo[node.Name] = schedulernodeinfo.NewNodeInfo()
		}
		nodeInfo := nodeNameToInfo[node.Name]
		nodeInfo.SetNode(node)
		nodeInfo.SetImageStates(getNodeImageStates(node, imageExistenceMap))
	}
	return nodeNameToInfo
}

// getNodeImageStates 从 imageExistenceMap 中查询目标 node 上所有镜像的 size 信息,
// 以及集群中拥有各个镜像的节点的数量(比如 node 拥有 centos:7 镜像, 而该存在于3个节点上).
//
// 	@param imageExistenceMap: key 为集群中存在的所有镜像名称, value 为拥有该镜像的所有
// node 节点的 nodeName 集合.
//
// getNodeImageStates returns the given node's image states based on the given
// imageExistence map.
func getNodeImageStates(
	node *v1.Node, imageExistenceMap map[string]sets.String,
) map[string]*schedulernodeinfo.ImageStateSummary {
	imageStates := make(map[string]*schedulernodeinfo.ImageStateSummary)

	for _, image := range node.Status.Images {
		for _, name := range image.Names {
			imageStates[name] = &schedulernodeinfo.ImageStateSummary{
				Size:     image.SizeBytes,
				NumNodes: len(imageExistenceMap[name]),
			}
		}
	}
	return imageStates
}

// createImageExistenceMap 梳理目标 nodes 列表上的所有镜像, 以镜像名称为 key 构建 map.
//
// 	@return: key 为 镜像名称, value 是存在该镜像的 Node 的 nodeName 集合.
//
// createImageExistenceMap returns a map recording on which nodes the images exist,
// keyed by the images' names.
func createImageExistenceMap(nodes []*v1.Node) map[string]sets.String {
	imageExistenceMap := make(map[string]sets.String)
	for _, node := range nodes {
		// 遍历当前 node 上存在的所有 image
		for _, image := range node.Status.Images {
			// 虽然同一个 image 的 image id 是相同的, 但这里还是按 image name 为 key
			for _, name := range image.Names {
				if _, ok := imageExistenceMap[name]; !ok {
					imageExistenceMap[name] = sets.NewString(node.Name)
				} else {
					imageExistenceMap[name].Insert(node.Name)
				}
			}
		}
	}
	return imageExistenceMap
}

// Pods returns a PodLister
func (s *Snapshot) Pods() schedulerlisters.PodLister {
	return &podLister{snapshot: s}
}

// NodeInfos returns a NodeInfoLister.
func (s *Snapshot) NodeInfos() schedulerlisters.NodeInfoLister {
	return &nodeInfoLister{snapshot: s}
}

// ListNodes returns the list of nodes in the snapshot.
func (s *Snapshot) ListNodes() []*v1.Node {
	nodes := make([]*v1.Node, 0, len(s.NodeInfoMap))
	for _, n := range s.NodeInfoList {
		if n.Node() != nil {
			nodes = append(nodes, n.Node())
		}
	}
	return nodes
}

type podLister struct {
	snapshot *Snapshot
}

// List returns the list of pods in the snapshot.
func (p *podLister) List(selector labels.Selector) ([]*v1.Pod, error) {
	alwaysTrue := func(p *v1.Pod) bool { return true }
	return p.FilteredList(alwaysTrue, selector)
}

// FilteredList returns a filtered list of pods in the snapshot.
func (p *podLister) FilteredList(podFilter schedulerlisters.PodFilter, selector labels.Selector) ([]*v1.Pod, error) {
	// podFilter is expected to return true for most or all of the pods.
	// We can avoid expensive array growth without wasting too much memory by
	// pre-allocating capacity.
	maxSize := 0
	for _, n := range p.snapshot.NodeInfoMap {
		maxSize += len(n.Pods())
	}
	pods := make([]*v1.Pod, 0, maxSize)
	for _, n := range p.snapshot.NodeInfoMap {
		for _, pod := range n.Pods() {
			if podFilter(pod) && selector.Matches(labels.Set(pod.Labels)) {
				pods = append(pods, pod)
			}
		}
	}
	return pods, nil
}

type nodeInfoLister struct {
	snapshot *Snapshot
}

// List returns the list of nodes in the snapshot.
func (n *nodeInfoLister) List() ([]*schedulernodeinfo.NodeInfo, error) {
	return n.snapshot.NodeInfoList, nil
}

// HavePodsWithAffinityList returns the list of nodes with at least one pods with inter-pod affinity
func (n *nodeInfoLister) HavePodsWithAffinityList() ([]*schedulernodeinfo.NodeInfo, error) {
	return n.snapshot.HavePodsWithAffinityNodeInfoList, nil
}

// Returns the NodeInfo of the given node name.
func (n *nodeInfoLister) Get(nodeName string) (*schedulernodeinfo.NodeInfo, error) {
	if v, ok := n.snapshot.NodeInfoMap[nodeName]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("nodeinfo not found for node name %q", nodeName)
}
