package ipvs

import (
	"k8s.io/apimachinery/pkg/util/sets"
	utilversion "k8s.io/apimachinery/pkg/util/version"
	utilipset "k8s.io/kubernetes/pkg/util/ipset"

	"fmt"
	"strings"

	"k8s.io/klog"
)

const (
	// MinIPSetCheckVersion is the min ipset version we need.  IPv6 is supported in ipset 6.x
	MinIPSetCheckVersion = "6.0"

	kubeLoopBackIPSetComment = "Kubernetes endpoints dst ip:port, source ip for solving hairpin purpose"
	// kubeLoopBackIPSet 该集合中存储的都是当前集群中的endpoints, 但不是全部, 而是恰好在当前宿主机上的ep.
	// 比如一个service下有两个pod, 那该service衍生的ep会有两个目标: PodIP1和PodIP2.
	// 如果pod分别调度在两个不同的节点上, pod1在worker1上, pod2在worker2上, 
	// 那么worker1上的kubeLoopBackIPSet就只有PodIP1, worker2则只有PodIP2.
	// 名字叫作KUBE-LOOP-BACK应该是为了体现就近原则, 如果访问的目标Pod刚好在本机, 那么直接在本机处理不香吗?
	kubeLoopBackIPSet        = "KUBE-LOOP-BACK"

	kubeClusterIPSetComment = "Kubernetes service cluster ip + port for masquerade purpose"
	// kubeClusterIPSet 的type为ip,port 存的是各service对象的clusterIP:port
	// kube-dns服务有3个port(53/UDP,53/TCP,9153/TCP), 那么在这个ipset中, 
	// 就会有3个成员, clusterIP相同而port不同.
	kubeClusterIPSet        = "KUBE-CLUSTER-IP"

	kubeExternalIPSetComment = "Kubernetes service external ip + port for masquerade and filter purpose"
	kubeExternalIPSet        = "KUBE-EXTERNAL-IP"

	kubeLoadBalancerSetComment = "Kubernetes service lb portal"
	kubeLoadBalancerSet        = "KUBE-LOAD-BALANCER"

	kubeLoadBalancerLocalSetComment = "Kubernetes service load balancer ip + port with externalTrafficPolicy=local"
	kubeLoadBalancerLocalSet        = "KUBE-LOAD-BALANCER-LOCAL"

	kubeLoadbalancerFWSetComment = "Kubernetes service load balancer ip + port for load balancer with sourceRange"
	kubeLoadbalancerFWSet        = "KUBE-LOAD-BALANCER-FW"

	kubeLoadBalancerSourceIPSetComment = "Kubernetes service load balancer ip + port + source IP for packet filter purpose"
	kubeLoadBalancerSourceIPSet        = "KUBE-LOAD-BALANCER-SOURCE-IP"

	kubeLoadBalancerSourceCIDRSetComment = "Kubernetes service load balancer ip + port + source cidr for packet filter purpose"
	kubeLoadBalancerSourceCIDRSet        = "KUBE-LOAD-BALANCER-SOURCE-CIDR"

	kubeNodePortSetTCPComment = "Kubernetes nodeport TCP port for masquerade purpose"
	kubeNodePortSetTCP        = "KUBE-NODE-PORT-TCP"

	kubeNodePortLocalSetTCPComment = "Kubernetes nodeport TCP port with externalTrafficPolicy=local"
	kubeNodePortLocalSetTCP        = "KUBE-NODE-PORT-LOCAL-TCP"

	kubeNodePortSetUDPComment = "Kubernetes nodeport UDP port for masquerade purpose"
	kubeNodePortSetUDP        = "KUBE-NODE-PORT-UDP"

	kubeNodePortLocalSetUDPComment = "Kubernetes nodeport UDP port with externalTrafficPolicy=local"
	kubeNodePortLocalSetUDP        = "KUBE-NODE-PORT-LOCAL-UDP"

	kubeNodePortSetSCTPComment = "Kubernetes nodeport SCTP port for masquerade purpose with type 'hash ip:port'"
	kubeNodePortSetSCTP        = "KUBE-NODE-PORT-SCTP-HASH"

	kubeNodePortLocalSetSCTPComment = "Kubernetes nodeport SCTP port with externalTrafficPolicy=local with type 'hash ip:port'"
	kubeNodePortLocalSetSCTP        = "KUBE-NODE-PORT-LOCAL-SCTP-HASH"
)

// IPSetVersioner can query the current ipset version.
type IPSetVersioner interface {
	// returns "X.Y"
	GetVersion() (string, error)
}

// IPSet 继承 utilipset.IPSet 结构
//
// IPSet wraps util/ipset which is used by IPVS proxier.
type IPSet struct {
	utilipset.IPSet
	// activeEntries is the current active entries of the ipset.
	activeEntries sets.String
	// handle is the util ipset interface handle.
	handle utilipset.Interface
}

// @param handle: pkg/util/ipset/ipset.go -> New(), 之后会通过exec接口进行操作.
//
// caller: 
// 	1. proxier.go -> NewProxier()
//
// NewIPSet initialize a new IPSet struct
func NewIPSet(
	handle utilipset.Interface, name string, setType utilipset.Type, 
	isIPv6 bool, comment string,
) *IPSet {
	hashFamily := utilipset.ProtocolFamilyIPV4
	if isIPv6 {
		hashFamily = utilipset.ProtocolFamilyIPV6
		// In dual-stack both ipv4 and ipv6 ipset's can co-exist. To
		// ensure unique names the prefix for ipv6 is changed from
		// "KUBE-" to "KUBE-6-". The "KUBE-" prefix is kept for
		// backward compatibility. The maximum name length of an ipset
		// is 31 characters which must be taken into account.  The
		// ipv4 names are not altered to minimize the risk for
		// problems on upgrades.
		if strings.HasPrefix(name, "KUBE-") {
			name = strings.Replace(name, "KUBE-", "KUBE-6-", 1)
			if len(name) > 31 {
				klog.Warningf("ipset name truncated; [%s] -> [%s]", name, name[:31])
				name = name[:31]
			}
		}
	}
	set := &IPSet{
		IPSet: utilipset.IPSet{
			Name:       name,
			SetType:    setType,
			HashFamily: hashFamily,
			Comment:    comment,
		},
		activeEntries: sets.NewString(),
		handle:        handle,
	}
	return set
}

func (set *IPSet) validateEntry(entry *utilipset.Entry) bool {
	return entry.Validate(&set.IPSet)
}

func (set *IPSet) isEmpty() bool {
	return len(set.activeEntries.UnsortedList()) == 0
}

func (set *IPSet) getComment() string {
	return fmt.Sprintf("\"%s\"", set.Comment)
}

func (set *IPSet) resetEntries() {
	set.activeEntries = sets.NewString()
}

func (set *IPSet) syncIPSetEntries() {
	appliedEntries, err := set.handle.ListEntries(set.Name)
	if err != nil {
		klog.Errorf("Failed to list ip set entries, error: %v", err)
		return
	}

	// currentIPSetEntries represents Endpoints watched from API Server.
	currentIPSetEntries := sets.NewString()
	for _, appliedEntry := range appliedEntries {
		currentIPSetEntries.Insert(appliedEntry)
	}

	if !set.activeEntries.Equal(currentIPSetEntries) {
		// Clean legacy entries
		for _, entry := range currentIPSetEntries.Difference(set.activeEntries).List() {
			if err := set.handle.DelEntry(entry, set.Name); err != nil {
				if !utilipset.IsNotFoundError(err) {
					klog.Errorf("Failed to delete ip set entry: %s from ip set: %s, error: %v", entry, set.Name, err)
				}
			} else {
				klog.V(3).Infof("Successfully delete legacy ip set entry: %s from ip set: %s", entry, set.Name)
			}
		}
		// Create active entries
		for _, entry := range set.activeEntries.Difference(currentIPSetEntries).List() {
			if err := set.handle.AddEntry(entry, &set.IPSet, true); err != nil {
				klog.Errorf("Failed to add entry: %v to ip set: %s, error: %v", entry, set.Name, err)
			} else {
				klog.V(3).Infof("Successfully add entry: %v to ip set: %s", entry, set.Name)
			}
		}
	}
}

// ensureIPSet 创建ipset项
func ensureIPSet(set *IPSet) error {
	if err := set.handle.CreateSet(&set.IPSet, true); err != nil {
		klog.Errorf("Failed to make sure ip set: %v exist, error: %v", set, err)
		return err
	}
	return nil
}

// checkMinVersion checks if ipset current version satisfies required min version
func checkMinVersion(vstring string) bool {
	version, err := utilversion.ParseGeneric(vstring)
	if err != nil {
		klog.Errorf("vstring (%s) is not a valid version string: %v", vstring, err)
		return false
	}

	minVersion, err := utilversion.ParseGeneric(MinIPSetCheckVersion)
	if err != nil {
		klog.Errorf("MinCheckVersion (%s) is not a valid version string: %v", MinIPSetCheckVersion, err)
		return false
	}
	return !version.LessThan(minVersion)
}
