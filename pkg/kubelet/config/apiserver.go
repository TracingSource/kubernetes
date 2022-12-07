package config

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// NewSourceApiserver 开始从 apiserver 监听目标主机上的 Pod 列表.
//
// 	@param c: 常规的 kube client
// 	@param nodeName: kubelet 当前所在的主机名称
// 	@param updates: 当 apiserver 中发生新的 Pod 事件时, 将该事件发送到 updates 通道中.
//         updates 中的事件的处理函数在 pkg/util/config/config.go -> Mux.listen()
//
// caller:
// 	1. pkg/kubelet/kubelet.go -> makePodSourceConfig()
//
// NewSourceApiserver creates a config source that watches and pulls from the apiserver.
func NewSourceApiserver(c clientset.Interface, nodeName types.NodeName, updates chan<- interface{}) {
	lw := cache.NewListWatchFromClient(
		c.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, 
		fields.OneTermEqualSelector(api.PodHostField, string(nodeName)),
	)
	newSourceApiserverFromLW(lw, updates)
}

// newSourceApiserverFromLW holds creates a config source that watches 
// and pulls from the apiserver.
func newSourceApiserverFromLW(lw cache.ListerWatcher, updates chan<- interface{}) {
	send := func(objs []interface{}) {
		var pods []*v1.Pod
		for _, o := range objs {
			pods = append(pods, o.(*v1.Pod))
		}
		updates <- kubetypes.PodUpdate{
			Pods: pods, Op: kubetypes.SET, Source: kubetypes.ApiserverSource,
		}
	}
	// 注意这里传入的 resyncPeriod 值为 0, 即不主动遍历本地缓存.
	r := cache.NewReflector(
		lw, &v1.Pod{}, cache.NewUndeltaStore(send, cache.MetaNamespaceKeyFunc), 0,
	)
	go r.Run(wait.NeverStop)
}
