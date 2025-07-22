package kubelet

import (
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/certificate"
	"k8s.io/client-go/util/flowcontrol"
	cloudprovider "k8s.io/cloud-provider"
	internalapi "k8s.io/cri-api/pkg/apis"
	api "k8s.io/kubernetes/pkg/apis/core"
	kubeletconfiginternal "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/cadvisor"
	"k8s.io/kubernetes/pkg/kubelet/cloudresource"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	"k8s.io/kubernetes/pkg/kubelet/config"
	"k8s.io/kubernetes/pkg/kubelet/configmap"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/dockershim"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	"k8s.io/kubernetes/pkg/kubelet/images"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/logs"
	"k8s.io/kubernetes/pkg/kubelet/network/dns"
	"k8s.io/kubernetes/pkg/kubelet/nodelease"
	oomwatcher "k8s.io/kubernetes/pkg/kubelet/oom"
	"k8s.io/kubernetes/pkg/kubelet/pleg"
	"k8s.io/kubernetes/pkg/kubelet/pluginmanager"
	kubepod "k8s.io/kubernetes/pkg/kubelet/pod"
	"k8s.io/kubernetes/pkg/kubelet/prober"
	proberesults "k8s.io/kubernetes/pkg/kubelet/prober/results"
	"k8s.io/kubernetes/pkg/kubelet/runtimeclass"
	"k8s.io/kubernetes/pkg/kubelet/secret"
	serverstats "k8s.io/kubernetes/pkg/kubelet/server/stats"
	"k8s.io/kubernetes/pkg/kubelet/stats"
	"k8s.io/kubernetes/pkg/kubelet/status"
	"k8s.io/kubernetes/pkg/kubelet/util/queue"
	"k8s.io/kubernetes/pkg/kubelet/volumemanager"
	"k8s.io/kubernetes/pkg/security/apparmor"
	utilipt "k8s.io/kubernetes/pkg/util/iptables"
	"k8s.io/utils/mount"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util/hostutil"
	"k8s.io/kubernetes/pkg/volume/util/subpath"
	corelisters "k8s.io/client-go/listers/core/v1"
)

// 	@initAt: pkg/kubelet/kubelet__new.go -> NewMainKubelet()
//
// Kubelet is the main kubelet implementation.
type Kubelet struct {
	kubeletConfiguration kubeletconfiginternal.KubeletConfiguration

	// hostname is the hostname the kubelet detected or was given via flag/config
	hostname string
	// hostnameOverridden indicates the hostname was overridden via flag/config
	hostnameOverridden bool

	// 当前 Node 的名称, 一般是 hostname
	// 在 pkg/kubelet/kubelet__new.go -> NewMainKubelet() 中被赋值.
	nodeName        types.NodeName
	runtimeCache    kubecontainer.RuntimeCache
	kubeClient      clientset.Interface
	heartbeatClient clientset.Interface
	iptClient       utilipt.Interface
	// 默认为 /var/lib/kubelet
	rootDirectory   string

	lastObservedNodeAddressesMux sync.RWMutex
	lastObservedNodeAddresses    []v1.NodeAddress

	// onRepeatedHeartbeatFailure is called when a heartbeat operation fails
	// more than once. optional.
	onRepeatedHeartbeatFailure func()

	// 在 NewMainKubelet() 函数中被初始化, 主要调用 Kubelet.syncPod() 方法
	// podWorkers handle syncing Pods in response to events.
	podWorkers PodWorkers

	// resyncInterval is the interval between periodic full reconciliations of
	// pods on this node.
	resyncInterval time.Duration

	// sourcesReady records the sources seen by the kubelet, it is thread-safe.
	sourcesReady config.SourcesReady

	// podManager is a facade that abstracts away the various sources of pods
	// this Kubelet services.
	podManager kubepod.Manager

	// Needed to observe and respond to situations that could impact node stability
	evictionManager eviction.Manager
	// logServer 默认为 true, 提供宿主机上 /var/log 目录的访问能力.
	//
	// Optional, defaults to /logs/ from /var/log
	logServer http.Handler
	// Optional, defaults to simple Docker implementation
	runner kubecontainer.ContainerCommandRunner

	// cAdvisor used for container information.
	cadvisor cadvisor.Interface

	// apiserver 是否已经得知当前 node 的存在.
	//
	// Set to true to have the node register itself with the apiserver.
	registerNode bool
	// List of taints to add to a node object when the kubelet registers itself.
	registerWithTaints []api.Taint
	// Set to true to have the node register itself as schedulable.
	registerSchedulable bool
	// for internal book keeping; access only from within registerWithApiserver
	registrationCompleted bool

	// dnsConfigurer is used for setting up DNS resolver configuration when launching pods.
	dnsConfigurer *dns.Configurer

	// masterServiceNamespace is the namespace that the master service is exposed in.
	masterServiceNamespace string
	// serviceLister knows how to list services
	serviceLister serviceLister
	// nodeLister knows how to list nodes
	nodeLister corelisters.NodeLister

	// a list of node labels to register
	nodeLabels map[string]string

	// Last timestamp when runtime responded on ping.
	// Mutex is used to protect this value.
	runtimeState *runtimeState

	// Volume plugins.
	volumePluginMgr *volume.VolumePluginMgr
	// 由 pkg/kubelet/prober/prober_manager.go -> NewManager() 构建
	//
	// Handles container probing.
	probeManager prober.Manager
	// Manages container health check results.
	livenessManager proberesults.Manager
	startupManager  proberesults.Manager

	// How long to keep idle streaming command execution/port forwarding
	// connections open before terminating them
	streamingConnectionIdleTimeout time.Duration

	// The EventRecorder to use
	recorder record.EventRecorder
	// containerGC 在 NewMainKubelet() 函数中被初始化.
	//
	// Policy for handling garbage collection of dead containers.
	containerGC kubecontainer.ContainerGC
	// imageManager 在 NewMainKubelet() 函数中被初始化.
	//
	// Manager for image garbage collection.
	imageManager images.ImageGCManager

	// Manager for container logs.
	containerLogManager logs.ContainerLogManager

	// Secret manager.
	secretManager secret.Manager

	// ConfigMap manager.
	configMapManager configmap.Manager

	// Cached MachineInfo returned by cadvisor.
	machineInfo *cadvisorapi.MachineInfo

	// Handles certificate rotations.
	serverCertificateManager certificate.Manager
	// statusManager 用于与 apiserver 保持通信, pod状态变动时主动通知, 否则定时上报.
	//
	// Syncs pods statuses with apiserver; also used as a cache of statuses.
	statusManager status.Manager

	// VolumeManager runs a set of asynchronous loops that figure out which
	// volumes need to be attached/mounted/unmounted/detached based on the pods
	// scheduled on this node and makes it so.
	volumeManager volumemanager.VolumeManager

	// Cloud provider interface.
	cloud cloudprovider.Interface
	// Handles requests to cloud provider with timeout
	cloudResourceSyncManager cloudresource.SyncManager

	// Indicates that the node initialization happens in an external cloud controller
	externalCloudProvider bool
	// Reference to this node.
	nodeRef *v1.ObjectReference
	// containerRuntimeName 有 docker 与 remote 两种, 一般为 docker
	//
	// The name of the container runtime
	containerRuntimeName string

	// redirectContainerStreaming enables container streaming redirect.
	redirectContainerStreaming bool
	// containerRuntime 有 docker, remote 两种
	//
	// 	@initAt: pkg/kubelet/kubelet__new.go -> NewMainKubelet() {
	//	pkg/kubelet/kuberuntime/kuberuntime_manager.go -> NewKubeGenericRuntimeManager()
	// }
	//
	// Container runtime.
	containerRuntime kubecontainer.Runtime

	// Streaming runtime handles container streaming.
	streamingRuntime kubecontainer.StreamingRuntime

	// 连接 dockershim.sock 或 containerd.sock 的 grpc 客户端对象
	//
	// Container runtime service (needed by container runtime Start()).
	// TODO(CD): try to make this available without holding a reference in this
	//           struct. For example, by adding a getter to generic runtime.
	runtimeService internalapi.RuntimeService

	// reasonCache caches the failure reason of the last creation of all containers,
	// which is used for generating ContainerStatus.
	reasonCache *ReasonCache

	// nodeStatusUpdateFrequency 计算本地 nodeStatus 变动的时间间隔.
	// 默认 10s, 见 pkg/kubelet/apis/config/v1beta1/defaults.go -> SetDefaults_KubeletConfiguration()
	//
	// 注意, 由于 node lease 的存在, 计算 node status 是独立的, 只有发生有意义的变化,
	// 或者超过上报间隔时, 才会真正上报.
	//
	// nodeStatusUpdateFrequency specifies how often kubelet computes node status.
	// If node lease feature is not enabled, it is also the frequency that
	// kubelet posts node status to master.
	//
	// In that case, be cautious when changing the constant, it must work with
	// nodeMonitorGracePeriod in nodecontroller. There are several constraints:
	// 1. nodeMonitorGracePeriod must be N times more than nodeStatusUpdateFrequency, where
	//    N means number of retries allowed for kubelet to post node status. It is pointless
	//    to make nodeMonitorGracePeriod be less than nodeStatusUpdateFrequency, since there
	//    will only be fresh values from Kubelet at an interval of nodeStatusUpdateFrequency.
	//    The constant must be less than podEvictionTimeout.
	// 2. nodeStatusUpdateFrequency needs to be large enough for kubelet to generate node
	//    status. Kubelet may fail to update node status reliably if the value is too small,
	//    as it takes time to gather all necessary node information.
	nodeStatusUpdateFrequency time.Duration

	// nodeStatusReportFrequency 上报 apiserver 的时间间隔.
	//
	// 默认 5m, 见 pkg/kubelet/apis/config/v1beta1/defaults.go -> SetDefaults_KubeletConfiguration()
	// 
	// nodeStatusReportFrequency is the frequency that kubelet posts node
	// status to master. It is only used when node lease feature is enabled.
	nodeStatusReportFrequency time.Duration

	// lastStatusReportTime is the time when node status was last reported.
	lastStatusReportTime time.Time
	// 每次调用 Kubelet.syncNodeStatus() 向 apiserver 同步 node 节点状态时,
	// 需要使用这个锁.
	//
	// syncNodeStatusMux is a lock on updating the node status,
	// because this path is not thread-safe.
	// This lock is used by Kubelet.syncNodeStatus function
	// and shouldn't be used anywhere else.
	syncNodeStatusMux sync.Mutex

	// updatePodCIDRMux is a lock on updating pod CIDR, because this path is not thread-safe.
	// This lock is used by Kubelet.syncNodeStatus function and shouldn't be used anywhere else.
	updatePodCIDRMux sync.Mutex

	// updateRuntimeMux is a lock on updating runtime, because this path is not thread-safe.
	// This lock is used by Kubelet.updateRuntimeUp function and shouldn't be used anywhere else.
	updateRuntimeMux sync.Mutex

	// nodeLeaseController claims and renews the node lease for this Kubelet
	nodeLeaseController nodelease.Controller

	// Generates pod events.
	pleg pleg.PodLifecycleEventGenerator

	// 生产者: pkg/kubelet/pleg/generic.go -> GenericPLEG.updateCache() 在 g.cache.Set()
	//
	// Store kubecontainer.PodStatus for all pods.
	podCache kubecontainer.Cache

	// os is a facade for various syscalls that need to be mocked during testing.
	os kubecontainer.OSInterface

	// Watcher of out of memory events.
	oomWatcher oomwatcher.Watcher

	// Monitor resource usage
	resourceAnalyzer serverstats.ResourceAnalyzer

	// Whether or not we should have the QOS cgroup hierarchy for resource management
	cgroupsPerQOS bool

	// 一般为 ""
	// If non-empty, pass this to the container runtime as the root cgroup.
	cgroupRoot string

	// Mounter to use for volumes.
	mounter mount.Interface

	// hostutil to interact with filesystems
	hostutil hostutil.HostUtils

	// subpather to execute subpath actions
	subpather subpath.Interface

	// Manager of non-Runtime containers.
	containerManager cm.ContainerManager

	// Maximum Number of Pods which can be run by this Kubelet
	maxPods int

	// Monitor Kubelet's sync loop
	syncLoopMonitor atomic.Value

	// Container restart Backoff
	backOff *flowcontrol.Backoff

	// Channel for sending pods to kill.
	podKillingCh chan *kubecontainer.PodPair

	// Information about the ports which are opened by daemons on Node running this Kubelet server.
	daemonEndpoints *v1.NodeDaemonEndpoints

	// A queue used to trigger pod workers.
	workQueue queue.WorkQueue

	// oneTimeInitializer is used to initialize modules that are dependent on the runtime to be up.
	oneTimeInitializer sync.Once

	// If non-nil, use this IP address for the node
	nodeIP net.IP

	// use this function to validate the kubelet nodeIP
	nodeIPValidator func(net.IP) error

	// If non-nil, this is a unique identifier for the node in an external database, eg. cloudprovider
	providerID string

	// clock is an interface that provides time related functionality in a way that makes it
	// easy to test the code.
	clock clock.Clock

	// 	@assignAs: pkg/kubelet/kubelet_node_status.go -> Kubelet.defaultNodeStatusFuncs()
	//
	// handlers called during the tryUpdateNodeStatus cycle
	setNodeStatusFuncs []func(*v1.Node) error

	lastNodeUnschedulableLock sync.Mutex
	// maintains Node.Spec.Unschedulable value from previous run of tryUpdateNodeStatus()
	lastNodeUnschedulable bool

	// TODO: think about moving this to be centralized in PodWorkers in follow-on.
	// the list of handlers to call during pod admission.
	admitHandlers lifecycle.PodAdmitHandlers

	// softAdmithandlers 在 pod 通过准入但真正运行前执行.
	//
	// softAdmithandlers are applied to the pod after it is admitted by the Kubelet,
	// but before it is run. 
	// A pod rejected by a softAdmitHandler will be left in a Pending state indefinitely.
	// If a rejected pod should not be recreated, or the scheduler is not aware
	// of the rejection rule, 
	// the admission rule should be applied by a softAdmitHandler.
	softAdmitHandlers lifecycle.PodAdmitHandlers

	// the list of handlers to call during pod sync loop.
	lifecycle.PodSyncLoopHandlers

	// the list of handlers to call during pod sync.
	lifecycle.PodSyncHandlers

	// the number of allowed pods per core
	podsPerCore int

	// enableControllerAttachDetach indicates the Attach/Detach controller
	// should manage attachment/detachment of volumes scheduled to this node,
	// and disable kubelet from executing any attach/detach operations
	enableControllerAttachDetach bool

	// trigger deleting containers in a pod
	containerDeletor *podContainerDeletor

	// config iptables util rules
	makeIPTablesUtilChains bool

	// The bit of the fwmark space to mark packets for SNAT.
	iptablesMasqueradeBit int

	// The bit of the fwmark space to mark packets for dropping.
	iptablesDropBit int

	// The AppArmor validator for checking whether AppArmor is supported.
	appArmorValidator apparmor.Validator

	// The handler serving CRI streaming calls (exec/attach/port-forward).
	criHandler http.Handler

	// experimentalHostUserNamespaceDefaulting sets userns=true when users request host namespaces (pid, ipc, net),
	// are using non-namespaced capabilities (mknod, sys_time, sys_module), the pod contains a privileged container,
	// or using host path volumes.
	// This should only be enabled when the container runtime is performing user remapping AND if the
	// experimental behavior is desired.
	experimentalHostUserNamespaceDefaulting bool

	// dockerLegacyService contains some legacy methods for backward compatibility.
	// It should be set only when docker is using non json-file logging driver.
	dockerLegacyService dockershim.DockerLegacyService

	// 	@assignAs: pkg/kubelet/stats/stats_provider.go -> StatsProvider{}
	//
	// StatsProvider provides the node and the container stats.
	*stats.StatsProvider

	// This flag, if set, instructs the kubelet to keep volumes from terminated pods mounted to the node.
	// This can be useful for debugging volume related issues.
	keepTerminatedPodVolumes bool // DEPRECATED

	// pluginmanager runs a set of asynchronous loops that figure out which
	// plugins need to be registered/unregistered based on this node and makes it so.
	pluginManager pluginmanager.PluginManager

	// This flag sets a maximum number of images to report in the node status.
	nodeStatusMaxImages int32
	// runtimeClassManager 这是一个 struct 对象(而非接口对象)
	//
	// Handles RuntimeClass objects for the Kubelet.
	runtimeClassManager *runtimeclass.Manager
}
