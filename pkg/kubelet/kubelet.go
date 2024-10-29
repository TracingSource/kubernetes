package kubelet

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"k8s.io/utils/mount"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	cloudprovider "k8s.io/cloud-provider"
	internalapi "k8s.io/cri-api/pkg/apis"
	"k8s.io/klog"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubeletconfiginternal "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/apis/podresources"
	"k8s.io/kubernetes/pkg/kubelet/cadvisor"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	"k8s.io/kubernetes/pkg/kubelet/config"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/dockershim"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	"k8s.io/kubernetes/pkg/kubelet/pleg"
	proberesults "k8s.io/kubernetes/pkg/kubelet/prober/results"
	"k8s.io/kubernetes/pkg/kubelet/remote"
	"k8s.io/kubernetes/pkg/kubelet/server"
	"k8s.io/kubernetes/pkg/kubelet/server/streaming"
	"k8s.io/kubernetes/pkg/kubelet/status"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/kubelet/util"
	"k8s.io/kubernetes/pkg/kubelet/util/format"
	"k8s.io/kubernetes/pkg/kubelet/util/sliceutils"
	"k8s.io/kubernetes/pkg/util/oom"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util/hostutil"
	"k8s.io/kubernetes/pkg/volume/util/subpath"
)

const (
	// Max amount of time to wait for the container runtime to come up.
	maxWaitForContainerRuntime = 30 * time.Second

	// nodeStatusUpdateRetry specifies how many times kubelet retries when posting node status failed.
	nodeStatusUpdateRetry = 5

	// ContainerLogsDir is the location of container logs.
	ContainerLogsDir = "/var/log/containers"

	// MaxContainerBackOff is the max backoff period, exported for the e2e test
	MaxContainerBackOff = 300 * time.Second

	// Capacity of the channel for storing pods to kill. A small number should
	// suffice because a goroutine is dedicated to check the channel and does
	// not block on anything else.
	podKillingChannelCapacity = 50

	// Period for performing global cleanup tasks.
	housekeepingPeriod = time.Second * 2

	// Period for performing eviction monitoring.
	// TODO ensure this is in sync with internal cadvisor housekeeping.
	evictionMonitoringPeriod = time.Second * 10

	// The path in containers' filesystems where the hosts file is mounted.
	etcHostsPath = "/etc/hosts"

	// Capacity of the channel for receiving pod lifecycle events. This number
	// is a bit arbitrary and may be adjusted in the future.
	plegChannelCapacity = 1000

	// Generic PLEG relies on relisting for discovering container events.
	// A longer period means that kubelet will take longer to detect container
	// changes and to update pod status. On the other hand, a shorter period
	// will cause more frequent relisting (e.g., container runtime operations),
	// leading to higher cpu usage.
	// Note that even though we set the period to 1s, the relisting itself can
	// take more than 1s to finish if the container runtime responds slowly
	// and/or when there are many container changes in one cycle.
	plegRelistPeriod = time.Second * 1

	// backOffPeriod is the period to back off when pod syncing results in an
	// error. It is also used as the base period for the exponential backoff
	// container restarts and image pulls.
	backOffPeriod = time.Second * 10

	// ContainerGCPeriod is the period for performing container garbage collection.
	ContainerGCPeriod = time.Minute
	// ImageGCPeriod is the period for performing image garbage collection.
	ImageGCPeriod = 5 * time.Minute

	// Minimum number of dead containers to keep in a pod
	minDeadContainerInPod = 1
)

// SyncHandler is an interface implemented by Kubelet, for testability
type SyncHandler interface {
	HandlePodAdditions(pods []*v1.Pod)
	HandlePodUpdates(pods []*v1.Pod)
	HandlePodRemoves(pods []*v1.Pod)
	HandlePodReconcile(pods []*v1.Pod)
	HandlePodSyncs(pods []*v1.Pod)
	HandlePodCleanups() error
}

// Option is a functional option type for Kubelet
type Option func(*Kubelet)

// Bootstrap 由 Kubelet{} 结构体实现(就在当前文件)
//
// Bootstrap is a bootstrapping interface for kubelet,
// targets the initialization protocol
type Bootstrap interface {
	GetConfiguration() kubeletconfiginternal.KubeletConfiguration
	BirthCry()
	StartGarbageCollection()
	ListenAndServe(
		address net.IP, port uint, tlsOptions *server.TLSOptions,
		auth server.AuthInterface, enableCAdvisorJSONEndpoints,
		enableDebuggingHandlers, enableContentionProfiling bool,
	)
	ListenAndServeReadOnly(address net.IP, port uint, enableCAdvisorJSONEndpoints bool)
	ListenAndServePodResources()
	Run(<-chan kubetypes.PodUpdate)
	RunOnce(<-chan kubetypes.PodUpdate) ([]RunPodResult, error)
}

// Builder creates and initializes a Kubelet instance
type Builder func(
	kubeCfg *kubeletconfiginternal.KubeletConfiguration,
	kubeDeps *Dependencies,
	crOptions *config.ContainerRuntimeOptions,
	containerRuntime string,
	runtimeCgroups string,
	hostnameOverride string,
	nodeIP string,
	providerID string,
	cloudProvider string,
	certDirectory string,
	rootDirectory string,
	registerNode bool,
	registerWithTaints []api.Taint,
	allowedUnsafeSysctls []string,
	remoteRuntimeEndpoint string,
	remoteImageEndpoint string,
	experimentalMounterPath string,
	experimentalKernelMemcgNotification bool,
	experimentalCheckNodeCapabilitiesBeforeMount bool,
	experimentalNodeAllocatableIgnoreEvictionThreshold bool,
	minimumGCAge metav1.Duration,
	maxPerPodContainerCount int32,
	maxContainerCount int32,
	masterServiceNamespace string,
	registerSchedulable bool,
	nonMasqueradeCIDR string,
	keepTerminatedPodVolumes bool,
	nodeLabels map[string]string,
	seccompProfileRoot string,
	bootstrapCheckpointPath string,
	nodeStatusMaxImages int32,
) (Bootstrap, error)

// Dependencies 一个可以作为依赖注入的集合, 里面集成了各种客户端, 各种插件.
//
// 初始化函数在 cmd/kubelet/app/server.go -> UnsecuredDependencies()
// 但是该函数构造的对象有很多是 nil 值, 还有一部分是在 cmd/kubelet/app/server.go -> run() 中赋值的.
//
// Dependencies is a bin for things we might consider "injected dependencies"
// -- objects constructed at runtime that are necessary for running the Kubelet.
// This is a temporary solution for grouping these objects
// while we figure out a more comprehensive dependency injection story for the Kubelet.
type Dependencies struct {
	Options []Option

	// Injected Dependencies
	Auth                    server.AuthInterface
	CAdvisorInterface       cadvisor.Interface
	Cloud                   cloudprovider.Interface
	ContainerManager        cm.ContainerManager
	DockerClientConfig      *dockershim.ClientConfig
	EventClient             v1core.EventsGetter
	HeartbeatClient         clientset.Interface
	OnHeartbeatFailure      func()
	KubeClient              clientset.Interface
	Mounter                 mount.Interface
	HostUtil                hostutil.HostUtils
	OOMAdjuster             *oom.OOMAdjuster
	OSInterface             kubecontainer.OSInterface
	PodConfig               *config.PodConfig
	Recorder                record.EventRecorder
	Subpather               subpath.Interface
	VolumePlugins           []volume.VolumePlugin
	DynamicPluginProber     volume.DynamicPluginProber
	TLSOptions              *server.TLSOptions
	KubeletConfigController *kubeletconfig.Controller
}

// caller:
// 	1. NewMainKubelet()
//
// makePodSourceConfig creates a config.PodConfig from the given
// KubeletConfiguration or returns an error.
func makePodSourceConfig(
	kubeCfg *kubeletconfiginternal.KubeletConfiguration, kubeDeps *Dependencies,
	nodeName types.NodeName, bootstrapCheckpointPath string,
) (*config.PodConfig, error) {
	manifestURLHeader := make(http.Header)
	if len(kubeCfg.StaticPodURLHeader) > 0 {
		for k, v := range kubeCfg.StaticPodURLHeader {
			for i := range v {
				manifestURLHeader.Add(k, v[i])
			}
		}
	}

	// source of all configuration
	cfg := config.NewPodConfig(config.PodConfigNotificationIncremental, kubeDeps.Recorder)

	// define file config source
	if kubeCfg.StaticPodPath != "" {
		klog.Infof("Adding pod path: %v", kubeCfg.StaticPodPath)
		config.NewSourceFile(
			kubeCfg.StaticPodPath, nodeName, kubeCfg.FileCheckFrequency.Duration,
			cfg.Channel(kubetypes.FileSource),
		)
	}

	// define url config source
	if kubeCfg.StaticPodURL != "" {
		klog.Infof(
			"Adding pod url %q with HTTP header %v", kubeCfg.StaticPodURL, manifestURLHeader,
		)
		config.NewSourceURL(
			kubeCfg.StaticPodURL, manifestURLHeader, nodeName,
			kubeCfg.HTTPCheckFrequency.Duration, cfg.Channel(kubetypes.HTTPSource),
		)
	}

	// Restore from the checkpoint path
	// NOTE: This MUST happen before creating the apiserver source
	// below, or the checkpoint would override the source of truth.

	var updatechannel chan<- interface{}
	if bootstrapCheckpointPath != "" {
		klog.Infof("Adding checkpoint path: %v", bootstrapCheckpointPath)
		updatechannel = cfg.Channel(kubetypes.ApiserverSource)
		err := cfg.Restore(bootstrapCheckpointPath, updatechannel)
		if err != nil {
			return nil, err
		}
	}

	if kubeDeps.KubeClient != nil {
		klog.Infof("Watching apiserver")
		if updatechannel == nil {
			updatechannel = cfg.Channel(kubetypes.ApiserverSource)
		}
		config.NewSourceApiserver(kubeDeps.KubeClient, nodeName, updatechannel)
	}
	return cfg, nil
}

// 	@param remoteRuntimeEndpoint: 有两种可能
// 	1. /var/run/dockershim.sock
// 	2. /var/run/containerd/containerd.sock
// 	@param remoteImageEndpoint: 一般与前者取值相同.
//
// caller:
// 	1. NewMainKubelet()
func getRuntimeAndImageServices(
	remoteRuntimeEndpoint string, remoteImageEndpoint string,
	runtimeRequestTimeout metav1.Duration,
) (internalapi.RuntimeService, internalapi.ImageManagerService, error) {
	rs, err := remote.NewRemoteRuntimeService(
		remoteRuntimeEndpoint, runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return nil, nil, err
	}
	is, err := remote.NewRemoteImageService(
		remoteImageEndpoint, runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return nil, nil, err
	}
	return rs, is, err
}

////////////////////////////////////////////////////////////////////////////////
// NewMainKubelet()
////////////////////////////////////////////////////////////////////////////////

type serviceLister interface {
	List(labels.Selector) ([]*v1.Service, error)
}

// StartGarbageCollection 清理停止的容器/不使用的镜像
//
// caller:
// 	1. cmd/kubelet/app/server.go -> createAndInitKubelet()
//
// StartGarbageCollection starts garbage collection threads.
func (kl *Kubelet) StartGarbageCollection() {
	loggedContainerGCFailure := false
	go wait.Until(func() {
		if err := kl.containerGC.GarbageCollect(); err != nil {
			klog.Errorf("Container garbage collection failed: %v", err)
			kl.recorder.Eventf(
				kl.nodeRef, v1.EventTypeWarning, events.ContainerGCFailed, err.Error(),
			)
			loggedContainerGCFailure = true
		} else {
			var vLevel klog.Level = 4
			if loggedContainerGCFailure {
				vLevel = 1
				loggedContainerGCFailure = false
			}

			klog.V(vLevel).Infof("Container garbage collection succeeded")
		}
	}, ContainerGCPeriod, wait.NeverStop)

	// when the high threshold is set to 100, stub the image GC manager
	if kl.kubeletConfiguration.ImageGCHighThresholdPercent == 100 {
		klog.V(2).Infof("ImageGCHighThresholdPercent is set 100, Disable image GC")
		return
	}

	prevImageGCFailed := false
	go wait.Until(func() {
		if err := kl.imageManager.GarbageCollect(); err != nil {
			if prevImageGCFailed {
				klog.Errorf("Image garbage collection failed multiple times in a row: %v", err)
				// Only create an event for repeated failures
				kl.recorder.Eventf(
					kl.nodeRef, v1.EventTypeWarning, events.ImageGCFailed, err.Error(),
				)
			} else {
				klog.Errorf("Image garbage collection failed once. Stats initialization may not have completed yet: %v", err)
			}
			prevImageGCFailed = true
		} else {
			var vLevel klog.Level = 4
			if prevImageGCFailed {
				vLevel = 1
				prevImageGCFailed = false
			}

			klog.V(vLevel).Infof("Image garbage collection succeeded")
		}
	}, ImageGCPeriod, wait.NeverStop)
}

// Run 启动各种子级协程, 监控volume的, 监控node状态的, 初始化iptables表什么的...
// 最后进入阻塞循环, 监听 apiserver 的 Pod 变动事件.
//
// 	@param updates: pkg/kubelet/config/config.go -> PodConfig.updates 通道.
//         接收来自 apiserver 的常规 Pod 与 manifests 目录下的 staticPod 变动事件.
//
// caller:
// 	1. cmd/kubelet/app/server.go -> startKubelet()
//
// Run starts the kubelet reacting to config updates
func (kl *Kubelet) Run(updates <-chan kubetypes.PodUpdate) {
	if kl.logServer == nil {
		kl.logServer = http.StripPrefix("/logs/", http.FileServer(http.Dir("/var/log/")))
	}
	if kl.kubeClient == nil {
		klog.Warning("No api server defined - no node status update will be sent.")
	}

	// Start the cloud provider sync manager
	if kl.cloudResourceSyncManager != nil {
		go kl.cloudResourceSyncManager.Run(wait.NeverStop)
	}

	if err := kl.initializeModules(); err != nil {
		kl.recorder.Eventf(
			kl.nodeRef, v1.EventTypeWarning, events.KubeletSetupFailed, err.Error(),
		)
		klog.Fatal(err)
	}

	// Start volume manager
	go kl.volumeManager.Run(kl.sourcesReady, wait.NeverStop)

	if kl.kubeClient != nil {
		// 同步心跳, 上报 node 节点的信息, 这其中会设置一些 docker 运行必要的东西.
		// 注意, 这里是通过定时器定时调用的, 永不停止, 这也是这个函数要求的.
		//
		// Start syncing node status immediately,
		// this may set up things the runtime needs to run.
		go wait.Until(kl.syncNodeStatus, kl.nodeStatusUpdateFrequency, wait.NeverStop)
		go kl.fastStatusUpdateOnce()
		// 1.16.2 版本中 NodeLease 是默认开启的, 到 1.17.2 就不是再是 feature 了.
		// start syncing lease
		go kl.nodeLeaseController.Run(wait.NeverStop)
	}
	go wait.Until(kl.updateRuntimeUp, 5*time.Second, wait.NeverStop)

	// Set up iptables util rules
	if kl.makeIPTablesUtilChains {
		kl.initNetworkUtil()
	}

	// Start a goroutine responsible for killing pods (that are not properly
	// handled by pod workers).
	go wait.Until(kl.podKiller, 1*time.Second, wait.NeverStop)

	// Start component sync loops.
	kl.statusManager.Start()
	kl.probeManager.Start()

	// Start syncing RuntimeClasses if enabled.
	if kl.runtimeClassManager != nil {
		kl.runtimeClassManager.Start(wait.NeverStop)
	}

	// Start the pod lifecycle event generator.
	kl.pleg.Start()
	kl.syncLoop(updates, kl)
}

// Get pods which should be resynchronized.
// Currently, the following pod should be resynchronized:
//   * pod whose work is ready.
//   * internal modules that request sync of a pod.
func (kl *Kubelet) getPodsToSync() []*v1.Pod {
	allPods := kl.podManager.GetPods()
	podUIDs := kl.workQueue.GetWork()
	podUIDSet := sets.NewString()
	for _, podUID := range podUIDs {
		podUIDSet.Insert(string(podUID))
	}
	var podsToSync []*v1.Pod
	for _, pod := range allPods {
		if podUIDSet.Has(string(pod.UID)) {
			// The work of the pod is ready
			podsToSync = append(podsToSync, pod)
			continue
		}
		for _, podSyncLoopHandler := range kl.PodSyncLoopHandlers {
			if podSyncLoopHandler.ShouldSync(pod) {
				podsToSync = append(podsToSync, pod)
				break
			}
		}
	}
	return podsToSync
}

// Kubelet.podKillingCh 成员在 NewMainKubelet() 函数中被初始化,
// 在 Kubelet.podKiller() 函数中处理其中的数据.
//
// caller:
// 	1. Kubelet.HandlePodRemoves() 只有这一处
//
// deletePod deletes the pod from the internal state of the kubelet by:
// 1.  stopping the associated pod worker asynchronously
// 2.  signaling to kill the pod by sending on the podKillingCh channel
//
// deletePod returns an error if not all sources are ready or the pod is not
// found in the runtime cache.
func (kl *Kubelet) deletePod(pod *v1.Pod) error {
	if pod == nil {
		return fmt.Errorf("deletePod does not allow nil pod")
	}
	if !kl.sourcesReady.AllReady() {
		// If the sources aren't ready, skip deletion, as we may accidentally delete pods
		// for sources that haven't reported yet.
		return fmt.Errorf("skipping delete because sources aren't ready yet")
	}
	kl.podWorkers.ForgetWorker(pod.UID)

	// Runtime cache may not have been updated to with the pod, but it's okay
	// because the periodic cleanup routine will attempt to delete again later.
	runningPods, err := kl.runtimeCache.GetPods()
	if err != nil {
		return fmt.Errorf("error listing containers: %v", err)
	}
	runningPod := kubecontainer.Pods(runningPods).FindPod("", pod.UID)
	if runningPod.IsEmpty() {
		return fmt.Errorf("pod not found")
	}
	podPair := kubecontainer.PodPair{APIPod: pod, RunningPod: &runningPod}

	kl.podKillingCh <- &podPair
	// TODO: delete the mirror pod here?

	// We leave the volume/directory cleanup to the periodic cleanup routine.
	return nil
}

// rejectPod records an event about the pod with the given reason and message,
// and updates the pod to the failed phase in the status manage.
func (kl *Kubelet) rejectPod(pod *v1.Pod, reason, message string) {
	kl.recorder.Eventf(pod, v1.EventTypeWarning, reason, message)
	kl.statusManager.SetPodStatus(pod, v1.PodStatus{
		Phase:   v1.PodFailed,
		Reason:  reason,
		Message: "Pod " + message})
}

// 	@param pod: 等待在当前节点上运行的 pod(已被调度到当前 Node)
// 	@param pods: 所有运行在当前节点上的 pod 列表
//
// caller:
// 	1. Kubelet.HandlePodAdditions() 处理常规的 Pod 新增事件(已被调度到当前 Node)
// 	2. pkg/kubelet/runonce.go -> Kubelet.runOnce() kubelet 一次性运行, 一般不会用到
//
// canAdmitPod determines if a pod can be admitted(被承认), and gives a reason
// if it cannot.
// "pod" is new pod, while "pods" are all admitted pods
// The function returns a boolean value indicating whether the pod can be admitted,
// a brief single-word reason and a message explaining why the pod cannot be admitted.
func (kl *Kubelet) canAdmitPod(pods []*v1.Pod, pod *v1.Pod) (bool, string, string) {
	// the kubelet will invoke each pod admit handler in sequence
	// if any handler rejects, the pod is rejected.
	// TODO: move out of disk check into a pod admitter
	// TODO: out of resource eviction should have a pod admitter call-out
	attrs := &lifecycle.PodAdmitAttributes{Pod: pod, OtherPods: pods}
	for _, podAdmitHandler := range kl.admitHandlers {
		if result := podAdmitHandler.Admit(attrs); !result.Admit {
			return false, result.Reason, result.Message
		}
	}

	return true, "", ""
}

// caller:
//	1. pkg/kubelet/kubelet__sync_pod.go -> Kubelet.syncPod()
func (kl *Kubelet) canRunPod(pod *v1.Pod) lifecycle.PodAdmitResult {
	attrs := &lifecycle.PodAdmitAttributes{Pod: pod}
	// Get "OtherPods". Rejected pods are failed, so only include admitted pods that are alive.
	attrs.OtherPods = kl.filterOutTerminatedPods(kl.podManager.GetPods())

	for _, handler := range kl.softAdmitHandlers {
		if result := handler.Admit(attrs); !result.Admit {
			return result
		}
	}

	return lifecycle.PodAdmitResult{Admit: true}
}

// syncLoop kubelet 主循环
//
// 	@param updates: pkg/kubelet/config/config.go -> PodConfig.updates 通道.
//         接收来自 apiserver 的常规 Pod 与 manifests 目录下的 staticPod 变动事件.
// 	@param handler: 其实传入的是 kl(Kubetlet) 对象本身.
//
// caller:
// 	1. Kubelet.Run()
//
// syncLoop is the main loop for processing changes. It watches for changes from
// three channels (file, apiserver, and http) and creates a union of them. For
// any new change seen, will run a sync against desired state and running state. If
// no changes are seen to the configuration, will synchronize the last known desired
// state every sync-frequency seconds. Never returns.
func (kl *Kubelet) syncLoop(updates <-chan kubetypes.PodUpdate, handler SyncHandler) {
	klog.Info("Starting kubelet main sync loop.")
	// The syncTicker wakes up kubelet to checks if there are any pod workers
	// that need to be sync'd. A one-second period is sufficient because the
	// sync interval is defaulted to 10s.
	syncTicker := time.NewTicker(time.Second)
	defer syncTicker.Stop()
	housekeepingTicker := time.NewTicker(housekeepingPeriod)
	defer housekeepingTicker.Stop()
	plegCh := kl.pleg.Watch()
	const (
		base   = 100 * time.Millisecond
		max    = 5 * time.Second
		factor = 2
	)
	duration := base
	// Responsible for checking limits in resolv.conf
	// The limits do not have anything to do with individual pods
	// Since this is called in syncLoop, we don't need to call it anywhere else
	if kl.dnsConfigurer != nil && kl.dnsConfigurer.ResolverConfig != "" {
		kl.dnsConfigurer.CheckLimitsForResolvConf()
	}

	for {
		if err := kl.runtimeState.runtimeErrors(); err != nil {
			klog.Errorf("skipping pod synchronization - %v", err)
			// exponential backoff
			time.Sleep(duration)
			duration = time.Duration(math.Min(float64(max), factor*float64(duration)))
			continue
		}
		// reset backoff if we have a success
		duration = base

		kl.syncLoopMonitor.Store(kl.clock.Now())
		if !kl.syncLoopIteration(updates, handler, syncTicker.C, housekeepingTicker.C, plegCh) {
			break
		}
		kl.syncLoopMonitor.Store(kl.clock.Now())
	}
}

// syncLoopIteration ...
//
// 	@param configCh: pkg/kubelet/config/config.go -> PodConfig.updates 通道.
//         接收来自 apiserver 的常规 Pod 与 manifests 目录下的 staticPod 变动事件.
// 	@param handler: 其实传入的是 kl(Kubetlet) 对象本身.
// 	@param syncCh: 1s的定时器
// 	@param housekeepingCh: 2s的定时器
//
// caller:
// 	1. Kubelet.syncLoop()
//
// syncLoopIteration reads from various channels and dispatches pods to the
// given handler.
//
// Arguments:
// 1.  configCh:       a channel to read config events from
// 2.  handler:        the SyncHandler to dispatch pods to
// 3.  syncCh:         a channel to read periodic sync events from
// 4.  housekeepingCh: a channel to read housekeeping events from
// 5.  plegCh:         a channel to read PLEG updates from
//
// Events are also read from the kubelet liveness manager's update channel.
//
// The workflow is to read from one of the channels, handle that event, and
// update the timestamp in the sync loop monitor.
//
// Here is an appropriate place to note that despite the syntactical
// similarity to the switch statement, the case statements in a select are
// evaluated in a pseudorandom order if there are multiple channels ready to
// read from when the select is evaluated.  In other words, case statements
// are evaluated in random order, and you can not assume that the case
// statements evaluate in order if multiple channels have events.
//
// With that in mind, in truly no particular order, the different channels
// are handled as follows:
//
// * configCh: dispatch the pods for the config change to the appropriate
//             handler callback for the event type
// * plegCh: update the runtime cache; sync pod
// * syncCh: sync all pods waiting for sync
// * housekeepingCh: trigger cleanup of pods
// * liveness manager: sync pods that have failed or in which one or more
//                     containers have failed liveness checks
func (kl *Kubelet) syncLoopIteration(
	configCh <-chan kubetypes.PodUpdate,
	handler SyncHandler,
	syncCh <-chan time.Time,
	housekeepingCh <-chan time.Time,
	plegCh <-chan *pleg.PodLifecycleEvent,
) bool {
	select {
	case u, open := <-configCh:
		// Update from a config source; dispatch it to the right handler callback.
		if !open {
			klog.Errorf("Update channel is closed. Exiting the sync loop.")
			return false
		}

		switch u.Op {
		case kubetypes.ADD:
			klog.V(2).Infof("SyncLoop (ADD, %q): %q", u.Source, format.Pods(u.Pods))
			// After restarting, kubelet will get all existing pods through
			// ADD as if they are new pods. These pods will then go through the
			// admission process and *may* be rejected. This can be resolved
			// once we have checkpointing.
			handler.HandlePodAdditions(u.Pods)
		case kubetypes.UPDATE:
			klog.V(2).Infof("SyncLoop (UPDATE, %q): %q", u.Source, format.PodsWithDeletionTimestamps(u.Pods))
			handler.HandlePodUpdates(u.Pods)
		case kubetypes.REMOVE:
			klog.V(2).Infof("SyncLoop (REMOVE, %q): %q", u.Source, format.Pods(u.Pods))
			handler.HandlePodRemoves(u.Pods)
		case kubetypes.RECONCILE:
			klog.V(4).Infof("SyncLoop (RECONCILE, %q): %q", u.Source, format.Pods(u.Pods))
			handler.HandlePodReconcile(u.Pods)
		case kubetypes.DELETE:
			klog.V(2).Infof("SyncLoop (DELETE, %q): %q", u.Source, format.Pods(u.Pods))
			// DELETE is treated as a UPDATE because of graceful deletion.
			handler.HandlePodUpdates(u.Pods)
		case kubetypes.RESTORE:
			klog.V(2).Infof("SyncLoop (RESTORE, %q): %q", u.Source, format.Pods(u.Pods))
			// These are pods restored from the checkpoint. Treat them as new
			// pods.
			handler.HandlePodAdditions(u.Pods)
		case kubetypes.SET:
			// TODO: Do we want to support this?
			klog.Errorf("Kubelet does not support snapshot update")
		}

		if u.Op != kubetypes.RESTORE {
			// If the update type is RESTORE, it means that the update is from
			// the pod checkpoints and may be incomplete.
			// Do not mark the source as ready.

			// Mark the source ready after receiving at least one update from the
			// source. Once all the sources are marked ready, various cleanup
			// routines will start reclaiming resources. It is important that this
			// takes place only after kubelet calls the update handler to process
			// the update to ensure the internal pod cache is up-to-date.
			kl.sourcesReady.AddSource(u.Source)
		}
	case e := <-plegCh:
		if isSyncPodWorthy(e) {
			// PLEG event for a pod; sync it.
			if pod, ok := kl.podManager.GetPodByUID(e.ID); ok {
				klog.V(2).Infof("SyncLoop (PLEG): %q, event: %#v", format.Pod(pod), e)
				handler.HandlePodSyncs([]*v1.Pod{pod})
			} else {
				// If the pod no longer exists, ignore the event.
				klog.V(4).Infof("SyncLoop (PLEG): ignore irrelevant event: %#v", e)
			}
		}

		if e.Type == pleg.ContainerDied {
			if containerID, ok := e.Data.(string); ok {
				kl.cleanUpContainersInPod(e.ID, containerID)
			}
		}
	case <-syncCh:
		// Sync pods waiting for sync
		podsToSync := kl.getPodsToSync()
		if len(podsToSync) == 0 {
			break
		}
		klog.V(4).Infof("SyncLoop (SYNC): %d pods; %s", len(podsToSync), format.Pods(podsToSync))
		handler.HandlePodSyncs(podsToSync)
	case update := <-kl.livenessManager.Updates():
		if update.Result == proberesults.Failure {
			// The liveness manager detected a failure; sync the pod.

			// We should not use the pod from livenessManager, because it is never updated after
			// initialization.
			pod, ok := kl.podManager.GetPodByUID(update.PodUID)
			if !ok {
				// If the pod no longer exists, ignore the update.
				klog.V(4).Infof("SyncLoop (container unhealthy): ignore irrelevant update: %#v", update)
				break
			}
			klog.V(1).Infof("SyncLoop (container unhealthy): %q", format.Pod(pod))
			handler.HandlePodSyncs([]*v1.Pod{pod})
		}
	case <-housekeepingCh:
		if !kl.sourcesReady.AllReady() {
			// If the sources aren't ready or volume manager has not yet synced the states,
			// skip housekeeping, as we may accidentally delete pods from unready sources.
			klog.V(4).Infof("SyncLoop (housekeeping, skipped): sources aren't ready yet.")
		} else {
			klog.V(4).Infof("SyncLoop (housekeeping)")
			if err := handler.HandlePodCleanups(); err != nil {
				klog.Errorf("Failed cleaning pods: %v", err)
			}
		}
	}
	return true
}

// 传入的参数中, mirrorPod 有可能为 nil, 即目标 Pod 不是某个 static 的镜像记录.
// ...我猜的.
//
// caller:
// 	1. Kubelet.HandlePodAdditions()
// 	2. Kubelet.HandlePodUpdates()
// 	3. ...
//
// dispatchWork starts the asynchronous sync of the pod in a pod worker.
// If the pod is terminated, dispatchWork
func (kl *Kubelet) dispatchWork(pod *v1.Pod, syncType kubetypes.SyncPodType, mirrorPod *v1.Pod, start time.Time) {
	if kl.podIsTerminated(pod) {
		if pod.DeletionTimestamp != nil {
			// If the pod is in a terminated state, there is no pod worker to
			// handle the work item. Check if the DeletionTimestamp has been
			// set, and force a status update to trigger a pod deletion request
			// to the apiserver.
			kl.statusManager.TerminatePod(pod)
		}
		return
	}
	// Run the sync in an async worker.
	kl.podWorkers.UpdatePod(&UpdatePodOptions{
		Pod:        pod,
		MirrorPod:  mirrorPod,
		UpdateType: syncType,
		OnCompleteFunc: func(err error) {
			if err != nil {
				metrics.PodWorkerDuration.WithLabelValues(syncType.String()).Observe(metrics.SinceInSeconds(start))
				metrics.DeprecatedPodWorkerLatency.WithLabelValues(syncType.String()).Observe(metrics.SinceInMicroseconds(start))
			}
		},
	})
	// Note the number of containers for new pods.
	if syncType == kubetypes.SyncPodCreate {
		metrics.ContainersPerPodCount.Observe(float64(len(pod.Spec.Containers)))
	}
}

// TODO: handle mirror pods in a separate component (issue #17251)
func (kl *Kubelet) handleMirrorPod(mirrorPod *v1.Pod, start time.Time) {
	// Mirror pod ADD/UPDATE/DELETE operations are considered an UPDATE to the
	// corresponding static pod. Send update to the pod worker if the static
	// pod exists.
	if pod, ok := kl.podManager.GetPodByMirrorPod(mirrorPod); ok {
		kl.dispatchWork(pod, kubetypes.SyncPodUpdate, mirrorPod, start)
	}
}

// HandlePodAdditions 将传入的 Pod 列表信息, 写入到本地缓存中.
//
// caller:
// 	1. Kubelet.syncLoopIteration() 在 kubelet 启动过各中被调用.
//     会被调用2次, 分别传入 mirrorPod 与常规 Pod 列表, 其中后者包括前者.
//
// HandlePodAdditions is the callback in SyncHandler for pods being added from
// a config source.
func (kl *Kubelet) HandlePodAdditions(pods []*v1.Pod) {
	start := kl.clock.Now()
	sort.Sort(sliceutils.PodsByCreationTime(pods))
	for _, pod := range pods {
		existingPods := kl.podManager.GetPods()
		// Always add the pod to the pod manager. Kubelet relies on the pod
		// manager as the source of truth for the desired state. If a pod does
		// not exist in the pod manager, it means that it has been deleted in
		// the apiserver and no action (other than cleanup) is required.
		kl.podManager.AddPod(pod)

		// 判断是不是某个 static pod 的 mirror pod.
		if kubetypes.IsMirrorPod(pod) {
			kl.handleMirrorPod(pod, start)
			continue
		}

		if !kl.podIsTerminated(pod) {
			// Only go through the admission process if the pod is not terminated.

			// We failed pods that we rejected, so activePods include all admitted
			// pods that are alive.
			activePods := kl.filterOutTerminatedPods(existingPods)

			// Check if we can admit the pod; if not, reject it.
			if ok, reason, message := kl.canAdmitPod(activePods, pod); !ok {
				kl.rejectPod(pod, reason, message)
				continue
			}
		}
		// 到这里, 说明目标 Pod 不是某个 mirror pod,
		// 返回的 mirrorPod 应该可能为 nil 吧...
		mirrorPod, _ := kl.podManager.GetMirrorPodByPod(pod)
		kl.dispatchWork(pod, kubetypes.SyncPodCreate, mirrorPod, start)
		kl.probeManager.AddPod(pod)
	}
}

// HandlePodUpdates is the callback in the SyncHandler interface for pods
// being updated from a config source.
func (kl *Kubelet) HandlePodUpdates(pods []*v1.Pod) {
	start := kl.clock.Now()
	for _, pod := range pods {
		kl.podManager.UpdatePod(pod)
		if kubetypes.IsMirrorPod(pod) {
			kl.handleMirrorPod(pod, start)
			continue
		}
		// TODO: Evaluate if we need to validate and reject updates.

		mirrorPod, _ := kl.podManager.GetMirrorPodByPod(pod)
		kl.dispatchWork(pod, kubetypes.SyncPodUpdate, mirrorPod, start)
	}
}

// HandlePodRemoves is the callback in the SyncHandler interface for pods
// being removed from a config source.
func (kl *Kubelet) HandlePodRemoves(pods []*v1.Pod) {
	start := kl.clock.Now()
	for _, pod := range pods {
		kl.podManager.DeletePod(pod)
		if kubetypes.IsMirrorPod(pod) {
			kl.handleMirrorPod(pod, start)
			continue
		}
		// Deletion is allowed to fail because the periodic cleanup routine
		// will trigger deletion again.
		if err := kl.deletePod(pod); err != nil {
			klog.V(2).Infof("Failed to delete pod %q, err: %v", format.Pod(pod), err)
		}
		kl.probeManager.RemovePod(pod)
	}
}

// HandlePodReconcile is the callback in the SyncHandler interface for pods
// that should be reconciled.
func (kl *Kubelet) HandlePodReconcile(pods []*v1.Pod) {
	start := kl.clock.Now()
	for _, pod := range pods {
		// Update the pod in pod manager, status manager will do periodically reconcile according
		// to the pod manager.
		kl.podManager.UpdatePod(pod)

		// Reconcile Pod "Ready" condition if necessary. Trigger sync pod for reconciliation.
		if status.NeedToReconcilePodReadiness(pod) {
			mirrorPod, _ := kl.podManager.GetMirrorPodByPod(pod)
			kl.dispatchWork(pod, kubetypes.SyncPodSync, mirrorPod, start)
		}

		// After an evicted pod is synced, all dead containers in the pod can be removed.
		if eviction.PodIsEvicted(pod.Status) {
			if podStatus, err := kl.podCache.Get(pod.UID); err == nil {
				kl.containerDeletor.deleteContainersInPod("", podStatus, true)
			}
		}
	}
}

// caller:
// 	1. Kubelet.syncLoopIteration() 貌似在只有 status 变动的时候就会调用此函数.
//
// HandlePodSyncs is the callback in the syncHandler interface for pods
// that should be dispatched to pod workers for sync.
func (kl *Kubelet) HandlePodSyncs(pods []*v1.Pod) {
	start := kl.clock.Now()
	for _, pod := range pods {
		mirrorPod, _ := kl.podManager.GetMirrorPodByPod(pod)
		kl.dispatchWork(pod, kubetypes.SyncPodSync, mirrorPod, start)
	}
}

// LatestLoopEntryTime returns the last time in the sync loop monitor.
func (kl *Kubelet) LatestLoopEntryTime() time.Time {
	val := kl.syncLoopMonitor.Load()
	if val == nil {
		return time.Time{}
	}
	return val.(time.Time)
}

// updateRuntimeUp 检测 docker/containerd 的运行状态, 并更新到 runtimeState 成员.
//
// caller:
// 	1. Kubelet.Run()
//
// updateRuntimeUp calls the container runtime status callback, initializing
// the runtime dependent modules when the container runtime first comes up,
// and returns an error if the status check fails.  If the status check is OK,
// update the container runtime uptime in the kubelet runtimeState.
func (kl *Kubelet) updateRuntimeUp() {
	kl.updateRuntimeMux.Lock()
	defer kl.updateRuntimeMux.Unlock()

	s, err := kl.containerRuntime.Status()
	if err != nil {
		klog.Errorf("Container runtime sanity check failed: %v", err)
		return
	}
	if s == nil {
		klog.Errorf("Container runtime status is nil")
		return
	}
	// Periodically log the whole runtime status for debugging.
	// TODO(random-liu): Consider to send node event when optional
	// condition is unmet.
	klog.V(4).Infof("Container runtime status: %v", s)
	networkReady := s.GetRuntimeCondition(kubecontainer.NetworkReady)
	if networkReady == nil || !networkReady.Status {
		klog.Errorf("Container runtime network not ready: %v", networkReady)
		kl.runtimeState.setNetworkState(
			fmt.Errorf("runtime network not ready: %v", networkReady),
		)
	} else {
		// Set nil if the container runtime network is ready.
		kl.runtimeState.setNetworkState(nil)
	}
	// information in RuntimeReady condition will be propagated to NodeReady condition.
	runtimeReady := s.GetRuntimeCondition(kubecontainer.RuntimeReady)
	// If RuntimeReady is not set or is false, report an error.
	if runtimeReady == nil || !runtimeReady.Status {
		err := fmt.Errorf("Container runtime not ready: %v", runtimeReady)
		klog.Error(err)
		kl.runtimeState.setRuntimeState(err)
		return
	}
	kl.runtimeState.setRuntimeState(nil)
	kl.oneTimeInitializer.Do(kl.initializeRuntimeDependentModules)
	kl.runtimeState.setRuntimeSync(kl.clock.Now())
}

// GetConfiguration returns the KubeletConfiguration used to configure the kubelet.
func (kl *Kubelet) GetConfiguration() kubeletconfiginternal.KubeletConfiguration {
	return kl.kubeletConfiguration
}

// BirthCry sends an event that the kubelet has started up.
func (kl *Kubelet) BirthCry() {
	// Make an event that kubelet restarted.
	kl.recorder.Eventf(kl.nodeRef, v1.EventTypeNormal, events.StartingKubelet, "Starting kubelet.")
}

// ResyncInterval returns the interval used for periodic syncs.
func (kl *Kubelet) ResyncInterval() time.Duration {
	return kl.resyncInterval
}

// ListenAndServe ...
//
// 	@param address:
// 	@param port:
// 	@param auth:
//
// caller:
// 	1. cmd/kubelet/app/server.go -> startKubelet()
//
// ListenAndServe runs the kubelet HTTP server.
func (kl *Kubelet) ListenAndServe(
	address net.IP, port uint, tlsOptions *server.TLSOptions,
	auth server.AuthInterface,
	enableCAdvisorJSONEndpoints, enableDebuggingHandlers, enableContentionProfiling bool,
) {
	server.ListenAndServeKubeletServer(
		kl, kl.resourceAnalyzer, address, port, tlsOptions, auth,
		enableCAdvisorJSONEndpoints, enableDebuggingHandlers,
		enableContentionProfiling, kl.redirectContainerStreaming, kl.criHandler,
	)
}

// ListenAndServeReadOnly runs the kubelet HTTP server in read-only mode.
func (kl *Kubelet) ListenAndServeReadOnly(
	address net.IP, port uint, enableCAdvisorJSONEndpoints bool,
) {
	server.ListenAndServeKubeletReadOnlyServer(
		kl, kl.resourceAnalyzer, address, port, enableCAdvisorJSONEndpoints,
	)
}

// ListenAndServePodResources runs the kubelet podresources grpc service
func (kl *Kubelet) ListenAndServePodResources() {
	socket, err := util.LocalEndpoint(kl.getPodResourcesDir(), podresources.Socket)
	if err != nil {
		klog.V(2).Infof("Failed to get local endpoint for PodResources endpoint: %v", err)
		return
	}
	server.ListenAndServePodResources(socket, kl.podManager, kl.containerManager)
}

// Delete the eligible dead container instances in a pod.
// Depending on the configuration, the latest dead containers may be kept around.
func (kl *Kubelet) cleanUpContainersInPod(podID types.UID, exitedContainerID string) {
	if podStatus, err := kl.podCache.Get(podID); err == nil {
		removeAll := false
		if syncedPod, ok := kl.podManager.GetPodByUID(podID); ok {
			// generate the api status using the cached runtime status to get up-to-date ContainerStatuses
			apiPodStatus := kl.generateAPIPodStatus(syncedPod, podStatus)
			// When an evicted or deleted pod has already synced, all containers can be removed.
			removeAll = eviction.PodIsEvicted(syncedPod.Status) || (syncedPod.DeletionTimestamp != nil && notRunning(apiPodStatus.ContainerStatuses))
		}
		kl.containerDeletor.deleteContainersInPod(exitedContainerID, podStatus, removeAll)
	}
}

// fastStatusUpdateOnce starts a loop that checks the internal node indexer cache
// for when a CIDR is applied  and tries to update pod CIDR immediately.
// After pod CIDR is updated it fires off a runtime update and a node status update.
// Function returns after one successful node status update.
// Function is executed only during Kubelet start which improves latency to
// ready node by updating pod CIDR, runtime status and node statuses ASAP.
func (kl *Kubelet) fastStatusUpdateOnce() {
	for {
		time.Sleep(100 * time.Millisecond)
		node, err := kl.GetNode()
		if err != nil {
			klog.Errorf(err.Error())
			continue
		}
		if len(node.Spec.PodCIDRs) != 0 {
			podCIDRs := strings.Join(node.Spec.PodCIDRs, ",")
			if _, err := kl.updatePodCIDR(podCIDRs); err != nil {
				klog.Errorf("Pod CIDR update to %v failed %v", podCIDRs, err)
				continue
			}
			kl.updateRuntimeUp()
			kl.syncNodeStatus()
			return
		}
	}
}

// isSyncPodWorthy filters out events that are not worthy of pod syncing
func isSyncPodWorthy(event *pleg.PodLifecycleEvent) bool {
	// ContainerRemoved doesn't affect pod state
	return event.Type != pleg.ContainerRemoved
}

// Gets the streaming server configuration to use with in-process CRI shims.
func getStreamingConfig(kubeCfg *kubeletconfiginternal.KubeletConfiguration, kubeDeps *Dependencies, crOptions *config.ContainerRuntimeOptions) *streaming.Config {
	config := &streaming.Config{
		StreamIdleTimeout:               kubeCfg.StreamingConnectionIdleTimeout.Duration,
		StreamCreationTimeout:           streaming.DefaultConfig.StreamCreationTimeout,
		SupportedRemoteCommandProtocols: streaming.DefaultConfig.SupportedRemoteCommandProtocols,
		SupportedPortForwardProtocols:   streaming.DefaultConfig.SupportedPortForwardProtocols,
	}
	if !crOptions.RedirectContainerStreaming {
		config.Addr = net.JoinHostPort("localhost", "0")
	} else {
		// Use a relative redirect (no scheme or host).
		config.BaseURL = &url.URL{
			Path: "/cri/",
		}
		if kubeDeps.TLSOptions != nil {
			config.TLSConfig = kubeDeps.TLSOptions.Config
		}
	}
	return config
}
