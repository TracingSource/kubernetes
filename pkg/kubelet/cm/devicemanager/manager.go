package devicemanager

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	errorsutil "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	"k8s.io/kubernetes/pkg/features"
	podresourcesapi "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/errors"
	cputopology "k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/topology"
	"k8s.io/kubernetes/pkg/kubelet/cm/devicemanager/checkpoint"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager/bitmask"
	"k8s.io/kubernetes/pkg/kubelet/config"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	"k8s.io/kubernetes/pkg/kubelet/pluginmanager/cache"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
	"k8s.io/kubernetes/pkg/util/selinux"
)

// ActivePodsFunc is a function that returns a list of pods to reconcile.
type ActivePodsFunc func() []*v1.Pod

// monitorCallback is the function called when a device's health state changes,
// or new devices are reported, or old devices are deleted.
// Updated contains the most recent state of the Device.
type monitorCallback func(resourceName string, devices []pluginapi.Device)

// 实现了 pkg/kubelet/cm/devicemanager/types.go -> Manager 接口.
//
// ManagerImpl is the structure in charge of managing Device Plugins.
type ManagerImpl struct {
	socketname string // kubelet.sock
	socketdir  string // /var/lib/kubelet/device-plugins

	endpoints map[string]endpointInfo // Key is ResourceName
	mutex     sync.Mutex

	server *grpc.Server
	wg     sync.WaitGroup

	// activePods is a method for listing active pods on the node
	// so the amount of pluginResources requested by existing pods
	// could be counted when updating allocated devices
	activePods ActivePodsFunc

	// sourcesReady provides the readiness of kubelet configuration sources
	// such as apiserver update readiness.
	// We use it to determine when we can purge inactive pods from checkpointed state.
	sourcesReady config.SourcesReady

	// callback is used for updating devices' states in one time call.
	// e.g. a new device is advertised, two old devices are deleted and a running device fails.
	callback monitorCallback

	// allDevices is a map by resource name of all the devices currently registered to the device manager
	allDevices map[string]map[string]pluginapi.Device

	// healthyDevices contains all of the registered healthy resourceNames and their exported device IDs.
	healthyDevices map[string]sets.String

	// unhealthyDevices contains all of the unhealthy devices and their exported device IDs.
	unhealthyDevices map[string]sets.String

	// 当前 node 节点上已分配给 Pod 的 deviceID 设备集合.
	// key 为扩展资源的名称(与 cpu/memory 同级), value 为设备列表(deviceID 列表)
	//
	// allocatedDevices contains allocated deviceIds, keyed by resourceName.
	allocatedDevices map[string]sets.String

	// podDevices 存储着已分配给所有pod的所有container的设备列表, 作为缓存.
	//
	// podDevices contains pod to allocated device mapping.
	podDevices podDevices
	// kubelet 会在 /var/lib/kubelet/device-plugins/kubelet_internal_checkpoint 文件中,
	// 存放所有分配给 pod/container 的设备信息, 防止 kubelet 自身发生重启后数据丢失.
	checkpointManager checkpointmanager.CheckpointManager

	// List of NUMA Nodes available on the underlying machine
	numaNodes []int

	// topologyAffinityStore: pkg/kubelet/cm/topologymanager/topology_manager.go -> manager{}
	//
	// Store of Topology Affinties that the Device Manager can query.
	topologyAffinityStore topologymanager.Store
}

type endpointInfo struct {
	e    endpoint
	opts *pluginapi.DevicePluginOptions
}

type sourcesReadyStub struct{}

func (s *sourcesReadyStub) AddSource(source string) {}
func (s *sourcesReadyStub) AllReady() bool          { return true }

// 	@param topologyAffinityStore: pkg/kubelet/cm/topologymanager/topology_manager.go -> manager{}
//
// NewManagerImpl creates a new manager.
func NewManagerImpl(
	numaNodeInfo cputopology.NUMANodeInfo, topologyAffinityStore topologymanager.Store,
) (*ManagerImpl, error) {
	return newManagerImpl(pluginapi.KubeletSocket, numaNodeInfo, topologyAffinityStore)
}

// 	@param socketPath: /var/lib/kubelet/device-plugins/kubelet.sock
func newManagerImpl(
	socketPath string, numaNodeInfo cputopology.NUMANodeInfo,
	topologyAffinityStore topologymanager.Store,
) (*ManagerImpl, error) {
	klog.V(2).Infof("Creating Device Plugin manager at %s", socketPath)

	if socketPath == "" || !filepath.IsAbs(socketPath) {
		return nil, fmt.Errorf(errBadSocket+" %s", socketPath)
	}

	var numaNodes []int
	for node := range numaNodeInfo {
		numaNodes = append(numaNodes, node)
	}

	dir, file := filepath.Split(socketPath)
	manager := &ManagerImpl{
		endpoints: make(map[string]endpointInfo),

		socketname:            file,
		socketdir:             dir,
		allDevices:            make(map[string]map[string]pluginapi.Device),
		healthyDevices:        make(map[string]sets.String),
		unhealthyDevices:      make(map[string]sets.String),
		allocatedDevices:      make(map[string]sets.String),
		podDevices:            make(podDevices),
		numaNodes:             numaNodes,
		topologyAffinityStore: topologyAffinityStore,
	}
	manager.callback = manager.genericDeviceUpdateCallback

	// The following structures are populated with real implementations in manager.Start()
	// Before that, initializes them to perform no-op operations.
	manager.activePods = func() []*v1.Pod { return []*v1.Pod{} }
	manager.sourcesReady = &sourcesReadyStub{}
	checkpointManager, err := checkpointmanager.NewCheckpointManager(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize checkpoint manager: %v", err)
	}
	manager.checkpointManager = checkpointManager

	return manager, nil
}

// device plugin 初次注册自身或宿主机上的设备状态或数量发生变动时,
// 由 device plugin 上报最新的设备列表信息, kubelet 最终会运行到这里进行更新.
//
// 	@param devices: 当前节点上全量的设备列表信息.
//
// caller:
// 	1. pkg/kubelet/cm/devicemanager/endpoint.go -> endpointImpl.callback()
// 	作为 endpointImpl 的 cb() 成员方法被调用.
func (m *ManagerImpl) genericDeviceUpdateCallback(
	resourceName string, devices []pluginapi.Device,
) {
	m.mutex.Lock()
	m.healthyDevices[resourceName] = sets.NewString()
	m.unhealthyDevices[resourceName] = sets.NewString()
	m.allDevices[resourceName] = make(map[string]pluginapi.Device)
	for _, dev := range devices {
		m.allDevices[resourceName][dev.ID] = dev
		if dev.Health == pluginapi.Healthy {
			m.healthyDevices[resourceName].Insert(dev.ID)
		} else {
			m.unhealthyDevices[resourceName].Insert(dev.ID)
		}
	}
	m.mutex.Unlock()
	m.writeCheckpoint()
}

func (m *ManagerImpl) removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	var errs []error
	for _, name := range names {
		filePath := filepath.Join(dir, name)
		if filePath == m.checkpointFile() {
			continue
		}
		stat, err := os.Stat(filePath)
		if err != nil {
			klog.Errorf("Failed to stat file %s: %v", filePath, err)
			continue
		}
		if stat.IsDir() {
			continue
		}
		err = os.RemoveAll(filePath)
		if err != nil {
			errs = append(errs, err)
			klog.Errorf("Failed to remove file %s: %v", filePath, err)
			continue
		}
	}
	return errorsutil.NewAggregate(errs)
}

// checkpointFile returns device plugin checkpoint file path.
func (m *ManagerImpl) checkpointFile() string {
	return filepath.Join(m.socketdir, kubeletDeviceManagerCheckpoint)
}

// caller:
// 	1. pkg/kubelet/cm/container_manager_linux.go -> containerManagerImpl.Start()
//
// Start starts the Device Plugin Manager and start initialization of
// podDevices and allocatedDevices information from checkpointed state and
// starts device plugin registration service.
func (m *ManagerImpl) Start(activePods ActivePodsFunc, sourcesReady config.SourcesReady) error {
	klog.V(2).Infof("Starting Device Plugin manager")

	m.activePods = activePods
	m.sourcesReady = sourcesReady

	// Loads in allocatedDevices information from disk.
	err := m.readCheckpoint()
	if err != nil {
		klog.Warningf(
			"Continue after failing to read checkpoint file. "+
				"Device allocation info may NOT be up-to-date. Err: %v", err,
		)
	}

	socketPath := filepath.Join(m.socketdir, m.socketname)
	if err = os.MkdirAll(m.socketdir, 0750); err != nil {
		return err
	}
	if selinux.SELinuxEnabled() {
		if err := selinux.SetFileLabel(m.socketdir, config.KubeletPluginsDirSELinuxLabel); err != nil {
			klog.Warningf("Unprivileged containerized plugins might not work. Could not set selinux context on %s: %v", m.socketdir, err)
		}
	}

	// Removes all stale sockets in m.socketdir. Device plugins can monitor
	// this and use it as a signal to re-register with the new Kubelet.
	if err := m.removeContents(m.socketdir); err != nil {
		klog.Errorf("Fail to clean up stale contents under %s: %v", m.socketdir, err)
	}

	s, err := net.Listen("unix", socketPath)
	if err != nil {
		klog.Errorf(errListenSocket+" %v", err)
		return err
	}

	m.wg.Add(1)
	m.server = grpc.NewServer([]grpc.ServerOption{}...)

	pluginapi.RegisterRegistrationServer(m.server, m)
	go func() {
		defer m.wg.Done()
		m.server.Serve(s)
	}()

	klog.V(2).Infof("Serving device plugin registration server on %q", socketPath)

	return nil
}

// GetWatcherHandler returns the plugin handler
func (m *ManagerImpl) GetWatcherHandler() cache.PluginHandler {
	if f, err := os.Create(m.socketdir + "DEPRECATION"); err != nil {
		klog.Errorf("Failed to create deprecation file at %s", m.socketdir)
	} else {
		f.Close()
		klog.V(4).Infof("created deprecation file %s", f.Name())
	}
	// 类型转换, 而非函数调用
	return cache.PluginHandler(m)
}

// ValidatePlugin validates a plugin if the version is correct and the name has
// the format of an extended resource
func (m *ManagerImpl) ValidatePlugin(
	pluginName string, endpoint string, versions []string,
) error {
	klog.V(2).Infof(
		"Got Plugin %s at endpoint %s with versions %v",
		pluginName, endpoint, versions,
	)

	if !m.isVersionCompatibleWithPlugin(versions) {
		return fmt.Errorf(
			"manager version, %s, is not among plugin supported versions %v",
			pluginapi.Version, versions,
		)
	}

	if !v1helper.IsExtendedResourceName(v1.ResourceName(pluginName)) {
		return fmt.Errorf(
			"invalid name of device plugin socket: %s",
			fmt.Sprintf(errInvalidResourceName, pluginName),
		)
	}

	return nil
}

// RegisterPlugin starts the endpoint and registers it
// TODO: Start the endpoint and wait for the First ListAndWatch call
//       before registering the plugin
func (m *ManagerImpl) RegisterPlugin(
	pluginName string, endpoint string, versions []string,
) error {
	klog.V(2).Infof("Registering Plugin %s at endpoint %s", pluginName, endpoint)

	e, err := newEndpointImpl(endpoint, pluginName, m.callback)
	if err != nil {
		return fmt.Errorf("failed to dial device plugin with socketPath %s: %v", endpoint, err)
	}

	options, err := e.client.GetDevicePluginOptions(context.Background(), &pluginapi.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get device plugin options: %v", err)
	}

	m.registerEndpoint(pluginName, options, e)
	go m.runEndpoint(pluginName, e)

	return nil
}

// DeRegisterPlugin deregisters the plugin
// TODO work on the behavior for deregistering plugins
// e.g: Should we delete the resource
func (m *ManagerImpl) DeRegisterPlugin(pluginName string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Note: This will mark the resource unhealthy as per the behavior
	// in runEndpoint
	if eI, ok := m.endpoints[pluginName]; ok {
		eI.e.stop()
	}
}

func (m *ManagerImpl) isVersionCompatibleWithPlugin(versions []string) bool {
	// TODO(vikasc): Currently this is fine as we only have a single supported version. When we do need to support
	// multiple versions in the future, we may need to extend this function to return a supported version.
	// E.g., say kubelet supports v1beta1 and v1beta2, and we get v1alpha1 and v1beta1 from a device plugin,
	// this function should return v1beta1
	for _, version := range versions {
		for _, supportedVersion := range pluginapi.SupportedVersions {
			if version == supportedVersion {
				return true
			}
		}
	}
	return false
}

// caller:
// 	1. ManagerImpl.Allocate()
// 	2. ManagerImpl.GetDeviceRunContainerOptions()
// 	kubelet 在调用 docker/cri 接口创建容器前会调用到这里.
func (m *ManagerImpl) allocatePodResources(pod *v1.Pod) error {
	devicesToReuse := make(map[string]sets.String)
	for _, container := range pod.Spec.InitContainers {
		if err := m.allocateContainerResources(pod, &container, devicesToReuse); err != nil {
			return err
		}
		m.podDevices.addContainerAllocatedResources(string(pod.UID), container.Name, devicesToReuse)
	}
	for _, container := range pod.Spec.Containers {
		if err := m.allocateContainerResources(pod, &container, devicesToReuse); err != nil {
			return err
		}
		// TODO 啥意思?
		m.podDevices.removeContainerAllocatedResources(string(pod.UID), container.Name, devicesToReuse)
	}
	return nil
}

// Allocate 在一个 Pod 已被调度到当前节点, 但还没有运行时, kubelet 需要为其做准入判断,
// device manager 会尝试为其分配扩展资源并写入本地 podDevice 缓存.
//
// 之后在启动 Pod 前, kubelet 会调用 ManagerImpl.GetDeviceRunContainerOptions()
//
// caller:
// 	1. pkg/kubelet/cm/container_manager_linux.go -> containerManagerImpl.UpdatePluginResources()
// 	实际创建容器时被调用, 用于为其分配确定的 deviceID 资源列表.
//
// Allocate is the call that you can use to allocate a set of devices
// from the registered device plugins.
func (m *ManagerImpl) Allocate(
	node *schedulernodeinfo.NodeInfo, attrs *lifecycle.PodAdmitAttributes,
) error {
	pod := attrs.Pod
	// 调用 device plugin 为 pod 分配设备列表, 会写入 kubelet 本地的 podDevice 缓存中.
	err := m.allocatePodResources(pod)
	if err != nil {
		klog.Errorf(
			"Failed to allocate device plugin resource for pod %s: %v", 
			string(pod.UID), err,
		)
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// quick return if no pluginResources requested
	if _, podRequireDevicePluginResource := m.podDevices[string(pod.UID)]; !podRequireDevicePluginResource {
		return nil
	}

	// 更新 node.status.allocatable 字段
	m.sanitizeNodeAllocatable(node)
	return nil
}

// Register 由某个 device plugin 通过 grpc 调用, 注册一个扩展资源
// 比如 pkg/kubelet/cm/devicemanager/device_plugin_stub.go -> Stub.Register
//
// Register registers a device plugin.
func (m *ManagerImpl) Register(
	ctx context.Context, r *pluginapi.RegisterRequest,
) (*pluginapi.Empty, error) {
	klog.Infof(
		"Got registration request from device plugin with resource name %q",
		r.ResourceName,
	)
	metrics.DevicePluginRegistrationCount.WithLabelValues(r.ResourceName).Inc()
	metrics.DeprecatedDevicePluginRegistrationCount.WithLabelValues(r.ResourceName).Inc()
	var versionCompatible bool
	for _, v := range pluginapi.SupportedVersions {
		if r.Version == v {
			versionCompatible = true
			break
		}
	}
	if !versionCompatible {
		errorString := fmt.Sprintf(errUnsupportedVersion, r.Version, pluginapi.SupportedVersions)
		klog.Infof(
			"Bad registration request from device plugin with resource name %q: %s",
			r.ResourceName, errorString,
		)
		return &pluginapi.Empty{}, fmt.Errorf(errorString)
	}

	if !v1helper.IsExtendedResourceName(v1.ResourceName(r.ResourceName)) {
		errorString := fmt.Sprintf(errInvalidResourceName, r.ResourceName)
		klog.Infof("Bad registration request from device plugin: %s", errorString)
		return &pluginapi.Empty{}, fmt.Errorf(errorString)
	}

	// TODO: for now, always accepts newest device plugin. Later may consider to
	// add some policies here, e.g., verify whether an old device plugin with the
	// same resource name is still alive to determine whether we want to accept
	// the new registration.
	go m.addEndpoint(r)

	return &pluginapi.Empty{}, nil
}

// Stop is the function that can stop the gRPC server.
// Can be called concurrently, more than once, and is safe to call
// without a prior Start.
func (m *ManagerImpl) Stop() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, eI := range m.endpoints {
		eI.e.stop()
	}

	if m.server == nil {
		return nil
	}
	m.server.Stop()
	m.wg.Wait()
	m.server = nil
	return nil
}

func (m *ManagerImpl) registerEndpoint(
	resourceName string, options *pluginapi.DevicePluginOptions, e endpoint,
) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.endpoints[resourceName] = endpointInfo{e: e, opts: options}
	klog.V(2).Infof("Registered endpoint %v", e)
}

func (m *ManagerImpl) runEndpoint(resourceName string, e endpoint) {
	e.run() // 启动 ListAndWatch() 函数进行监听
	e.stop()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if old, ok := m.endpoints[resourceName]; ok && old.e == e {
		m.markResourceUnhealthy(resourceName)
	}

	klog.V(2).Infof("Endpoint (%s, %v) became unhealthy", resourceName, e)
}

// caller:
// 	1. ManagerImpl.Register() kubelet 接收到来自 device plugin 的注册请求时被调用
func (m *ManagerImpl) addEndpoint(r *pluginapi.RegisterRequest) {
	new, err := newEndpointImpl(filepath.Join(m.socketdir, r.Endpoint), r.ResourceName, m.callback)
	if err != nil {
		klog.Errorf("Failed to dial device plugin with request %v: %v", r, err)
		return
	}
	m.registerEndpoint(r.ResourceName, r.Options, new)
	go func() {
		m.runEndpoint(r.ResourceName, new)
	}()
}

func (m *ManagerImpl) markResourceUnhealthy(resourceName string) {
	klog.V(2).Infof("Mark all resources Unhealthy for resource %s", resourceName)
	healthyDevices := sets.NewString()
	if _, ok := m.healthyDevices[resourceName]; ok {
		healthyDevices = m.healthyDevices[resourceName]
		m.healthyDevices[resourceName] = sets.NewString()
	}
	if _, ok := m.unhealthyDevices[resourceName]; !ok {
		m.unhealthyDevices[resourceName] = sets.NewString()
	}
	m.unhealthyDevices[resourceName] = m.unhealthyDevices[resourceName].Union(healthyDevices)
}

// GetCapacity is expected to be called when Kubelet updates its node status.
// The first returned variable contains the registered device plugin resource capacity.
// The second returned variable contains the registered device plugin resource allocatable.
// The third returned variable contains previously registered resources that are no longer active.
// Kubelet uses this information to update resource capacity/allocatable in its node status.
// After the call, device plugin can remove the inactive resources from its internal list as the
// change is already reflected in Kubelet node status.
// Note in the special case after Kubelet restarts, device plugin resource capacities can
// temporarily drop to zero till corresponding device plugins re-register. This is OK because
// cm.UpdatePluginResource() run during predicate Admit guarantees we adjust nodeinfo
// capacity for already allocated pods so that they can continue to run. However, new pods
// requiring device plugin resources will not be scheduled till device plugin re-registers.
func (m *ManagerImpl) GetCapacity() (v1.ResourceList, v1.ResourceList, []string) {
	needsUpdateCheckpoint := false
	var capacity = v1.ResourceList{}
	var allocatable = v1.ResourceList{}
	deletedResources := sets.NewString()
	m.mutex.Lock()
	for resourceName, devices := range m.healthyDevices {
		eI, ok := m.endpoints[resourceName]
		if (ok && eI.e.stopGracePeriodExpired()) || !ok {
			// The resources contained in endpoints and (un)healthyDevices
			// should always be consistent. Otherwise, we run with the risk
			// of failing to garbage collect non-existing resources or devices.
			if !ok {
				klog.Errorf("unexpected: healthyDevices and endpoints are out of sync")
			}
			delete(m.endpoints, resourceName)
			delete(m.healthyDevices, resourceName)
			deletedResources.Insert(resourceName)
			needsUpdateCheckpoint = true
		} else {
			capacity[v1.ResourceName(resourceName)] = *resource.NewQuantity(int64(devices.Len()), resource.DecimalSI)
			allocatable[v1.ResourceName(resourceName)] = *resource.NewQuantity(int64(devices.Len()), resource.DecimalSI)
		}
	}
	for resourceName, devices := range m.unhealthyDevices {
		eI, ok := m.endpoints[resourceName]
		if (ok && eI.e.stopGracePeriodExpired()) || !ok {
			if !ok {
				klog.Errorf("unexpected: unhealthyDevices and endpoints are out of sync")
			}
			delete(m.endpoints, resourceName)
			delete(m.unhealthyDevices, resourceName)
			deletedResources.Insert(resourceName)
			needsUpdateCheckpoint = true
		} else {
			capacityCount := capacity[v1.ResourceName(resourceName)]
			unhealthyCount := *resource.NewQuantity(int64(devices.Len()), resource.DecimalSI)
			capacityCount.Add(unhealthyCount)
			capacity[v1.ResourceName(resourceName)] = capacityCount
		}
	}
	m.mutex.Unlock()
	if needsUpdateCheckpoint {
		m.writeCheckpoint()
	}
	return capacity, allocatable, deletedResources.UnsortedList()
}

// Checkpoints device to container allocation information to disk.
func (m *ManagerImpl) writeCheckpoint() error {
	m.mutex.Lock()
	registeredDevs := make(map[string][]string)
	for resource, devices := range m.healthyDevices {
		registeredDevs[resource] = devices.UnsortedList()
	}
	data := checkpoint.New(m.podDevices.toCheckpointData(), registeredDevs)
	m.mutex.Unlock()
	err := m.checkpointManager.CreateCheckpoint(kubeletDeviceManagerCheckpoint, data)
	if err != nil {
		err2 := fmt.Errorf("failed to write checkpoint file %q: %v", kubeletDeviceManagerCheckpoint, err)
		klog.Warning(err2)
		return err2
	}
	return nil
}

// Reads device to container allocation information from disk, and populates
// m.allocatedDevices accordingly.
func (m *ManagerImpl) readCheckpoint() error {
	registeredDevs := make(map[string][]string)
	devEntries := make([]checkpoint.PodDevicesEntry, 0)
	cp := checkpoint.New(devEntries, registeredDevs)
	err := m.checkpointManager.GetCheckpoint(kubeletDeviceManagerCheckpoint, cp)
	if err != nil {
		if err == errors.ErrCheckpointNotFound {
			klog.Warningf("Failed to retrieve checkpoint for %q: %v", kubeletDeviceManagerCheckpoint, err)
			return nil
		}
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	podDevices, registeredDevs := cp.GetData()
	m.podDevices.fromCheckpointData(podDevices)
	m.allocatedDevices = m.podDevices.devices()
	for resource := range registeredDevs {
		// During start up, creates empty healthyDevices list so that the resource capacity
		// will stay zero till the corresponding device plugin re-registers.
		m.healthyDevices[resource] = sets.NewString()
		m.unhealthyDevices[resource] = sets.NewString()
		m.endpoints[resource] = endpointInfo{e: newStoppedEndpointImpl(resource), opts: nil}
	}
	return nil
}

// updateAllocatedDevices gets a list of active pods and then frees any Devices that are bound to
// terminated pods. Returns error on failure.
func (m *ManagerImpl) updateAllocatedDevices(activePods []*v1.Pod) {
	if !m.sourcesReady.AllReady() {
		return
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	activePodUids := sets.NewString()
	for _, pod := range activePods {
		activePodUids.Insert(string(pod.UID))
	}
	allocatedPodUids := m.podDevices.pods()
	podsToBeRemoved := allocatedPodUids.Difference(activePodUids)
	if len(podsToBeRemoved) <= 0 {
		return
	}
	klog.V(3).Infof("pods to be removed: %v", podsToBeRemoved.List())
	m.podDevices.delete(podsToBeRemoved.List())
	// Regenerated allocatedDevices after we update pod allocation information.
	m.allocatedDevices = m.podDevices.devices()
}

// 	@param podUID: 待分配扩展资源的 pod 的 uid
// 	@param podUID: 表示的 pod 中的某一 container 名称(resources{}字段都是配置在 container 中的).
// 	@param resource: 扩展资源名称, 与 cpu/memory 平级
// 	@param required: 需要分配的数量
//
// Returns list of device Ids we need to allocate with Allocate rpc call.
// Returns empty list in case we don't need to issue the Allocate rpc call.
func (m *ManagerImpl) devicesToAllocate(
	podUID, contName, resource string, required int, reusableDevices sets.String,
) (sets.String, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	needed := required
	// 获取已分配给当前 container 的设备列表(一般发生在容器发生重启的场景)
	//
	// Gets list of devices that have already been allocated.
	// This can happen if a container restarts for example.
	devices := m.podDevices.containerDevices(podUID, contName, resource)
	if devices != nil {
		klog.V(3).Infof(
			"Found pre-allocated devices for resource %s container %q in Pod %q: %v",
			resource, contName, podUID, devices.List(),
		)
		needed = needed - devices.Len()
		// 容器重启不应影响已分配的设备列表, 如果数量不一致则报错.
		//
		// A pod's resource is not expected to change once admitted by the API server,
		// so just fail loudly here. We can revisit this part if this no longer holds.
		if needed != 0 {
			return nil, fmt.Errorf(
				"pod %q container %q changed request for resource %q from %d to %d",
				podUID, contName, resource, devices.Len(), required,
			)
		}
	}
	// 如果当前容器没有 request 扩展资源, 或是重启后找到了原本分配给自己的资源,
	// 则什么都不用做, 直接退出.
	if needed == 0 {
		// No change, no work.
		return nil, nil
	}
	klog.V(3).Infof(
		"Needs to allocate %d %q for pod %q container %q",
		needed, resource, podUID, contName,
	)
	// 对应的 device plugin worker没有在工作?
	//
	// Needs to allocate additional devices.
	if _, ok := m.healthyDevices[resource]; !ok {
		return nil, fmt.Errorf("can't allocate unregistered device %s", resource)
	}

	// devices 表示作为结果返回的设备列表
	devices = sets.NewString()
	// 先从可复用的设备中选一部分.
	//
	// Allocates from reusableDevices list first.
	for device := range reusableDevices {
		devices.Insert(device)
		needed--
		if needed == 0 {
			return devices, nil
		}
	}
	// Needs to allocate additional devices.
	if m.allocatedDevices[resource] == nil {
		m.allocatedDevices[resource] = sets.NewString()
	}
	// 从正常工作的 device 中, 排除已经被分配了的, 得到剩余可分配的列表.
	//
	// Gets Devices in use.
	devicesInUse := m.allocatedDevices[resource]
	// Gets a list of available devices.
	available := m.healthyDevices[resource].Difference(devicesInUse)
	// 剩余的不够分了, 则报错.
	if available.Len() < needed {
		return nil, fmt.Errorf(
			"requested number of devices unavailable for %s. Requested: %d, Available: %d",
			resource, needed, available.Len(),
		)
	}
	// 从可分配列表(无序)中截取出指定长度的子列表
	// By default, pull devices from the unsorted list of available devices.
	allocated := available.UnsortedList()[:needed]
	// If topology alignment is desired, update allocated to the set of devices
	// with the best alignment.
	hint := m.topologyAffinityStore.GetAffinity(podUID, contName)
	if m.deviceHasTopologyAlignment(resource) && hint.NUMANodeAffinity != nil {
		allocated = m.takeByTopology(resource, available, hint.NUMANodeAffinity, needed)
	}
	// 将截取出的 allocated 部分标记为"已分配".
	//
	// Updates m.allocatedDevices with allocated devices to prevent them
	// from being allocated to other pods/containers, given that we are
	// not holding lock during the rpc call.
	for _, device := range allocated {
		m.allocatedDevices[resource].Insert(device)
		devices.Insert(device)
	}
	return devices, nil
}

func (m *ManagerImpl) takeByTopology(
	resource string, available sets.String, affinity bitmask.BitMask, request int,
) []string {
	// Build a map of NUMA Nodes to the devices associated with them. A
	// device may be associated to multiple NUMA nodes at the same time. If an
	// available device does not have any NUMA Nodes associated with it, add it
	// to a list of NUMA Nodes for the fake NUMANode -1.
	perNodeDevices := make(map[int]sets.String)
	for d := range available {
		var nodes []int
		if m.allDevices[resource][d].Topology != nil {
			for _, node := range m.allDevices[resource][d].Topology.Nodes {
				nodes = append(nodes, int(node.ID))
			}
		}
		if len(nodes) == 0 {
			nodes = []int{-1}
		}
		for _, node := range nodes {
			if _, ok := perNodeDevices[node]; !ok {
				perNodeDevices[node] = sets.NewString()
			}
			perNodeDevices[node].Insert(d)
		}
	}

	// Get a flat list of all of the nodes associated with available devices.
	var nodes []int
	for node := range perNodeDevices {
		nodes = append(nodes, node)
	}

	// Sort the list of nodes by how many devices they contain.
	sort.Slice(nodes, func(i, j int) bool {
		return perNodeDevices[i].Len() < perNodeDevices[j].Len()
	})

	// Generate three sorted lists of devices. Devices in the first list come
	// from valid NUMA Nodes contained in the affinity mask. Devices in the
	// second list come from valid NUMA Nodes not in the affinity mask. Devices
	// in the third list come from devices with no NUMA Node association (i.e.
	// those mapped to the fake NUMA Node -1). Because we loop through the
	// sorted list of NUMA nodes in order, within each list, devices are sorted
	// by their connection to NUMA Nodes with more devices on them.
	var fromAffinity []string
	var notFromAffinity []string
	var withoutTopology []string
	for d := range available {
		// Since the same device may be associated with multiple NUMA Nodes. We
		// need to be careful not to add each device to multiple lists. The
		// logic below ensures this by breaking after the first NUMA node that
		// has the device is encountered.
		for _, n := range nodes {
			if perNodeDevices[n].Has(d) {
				if n == -1 {
					withoutTopology = append(withoutTopology, d)
				} else if affinity.IsSet(n) {
					fromAffinity = append(fromAffinity, d)
				} else {
					notFromAffinity = append(notFromAffinity, d)
				}
				break
			}
		}
	}

	// Concatenate the lists above return the first 'request' devices from it..
	return append(append(fromAffinity, notFromAffinity...), withoutTopology...)[:request]
}

// 	@param container: pod 中包含的 container (可以是 InitContainer)
//
// allocateContainerResources attempts to allocate all of required device
// plugin resources for the input container, issues an Allocate rpc request
// for each new device resource requirement, processes their AllocateResponses,
// and updates the cached containerDevices on success.
func (m *ManagerImpl) allocateContainerResources(
	pod *v1.Pod, container *v1.Container, devicesToReuse map[string]sets.String,
) error {
	podUID := string(pod.UID)
	contName := container.Name
	allocatedDevicesUpdated := false
	// 扩展资源不允许超分, requests/limits 值需要一致.
	// Extended resources are not allowed to be overcommitted.
	// Since device plugin advertises extended resources,
	// therefore Requests must be equal to Limits and iterating
	// over the Limits should be sufficient.
	for k, v := range container.Resources.Limits {
		resource := string(k)
		needed := int(v.Value())
		// 比如 needs 1 cpu, needs 52428800 memory 等, 但这些不算是 device 类型的资源.
		klog.V(3).Infof("needs %d %s", needed, resource)
		if !m.isDevicePluginResource(resource) {
			continue
		}
		// allocatedDevicesUpdated 是一个开关, updateAllocatedDevices 只执行一次.
		//
		// Updates allocatedDevices to garbage collect any stranded resources
		// before doing the device plugin allocation.
		if !allocatedDevicesUpdated {
			m.updateAllocatedDevices(m.activePods())
			allocatedDevicesUpdated = true
		}
		// 从本地查询到可分配的设备列表
		allocDevices, err := m.devicesToAllocate(
			podUID, contName, resource, needed, devicesToReuse[resource],
		)
		if err != nil {
			return err
		}
		if allocDevices == nil || len(allocDevices) <= 0 {
			continue
		}

		startRPCTime := time.Now()
		// Manager.Allocate involves RPC calls to device plugin, which
		// could be heavy-weight.
		// Therefore we want to perform this operation outside mutex lock.
		// Note if Allocate call fails, we may leave container resources
		// partially allocated for the failed container. We rely on updateAllocatedDevices()
		// to garbage collect these resources later. Another side effect is that if
		// we have X resource A and Y resource B in total, and two containers, container1
		// and container2 both require X resource A and Y resource B.
		// Both allocation requests may fail if we serve them in mixed order.
		//
		// TODO: may revisit this part later if we see inefficient resource allocation
		// in real use as the result of this. Should also consider to parallelize device
		// plugin Allocate grpc calls if it becomes common that a container may require
		// resources from multiple device plugins.
		m.mutex.Lock()
		eI, ok := m.endpoints[resource]
		m.mutex.Unlock()
		if !ok {
			m.mutex.Lock()
			m.allocatedDevices = m.podDevices.devices()
			m.mutex.Unlock()
			return fmt.Errorf("unknown Device Plugin %s", resource)
		}

		devs := allocDevices.UnsortedList()
		// TODO: refactor this part of code to just append a ContainerAllocationRequest
		// in a passed in AllocateRequest pointer, and issues a single Allocate call per pod.
		klog.V(3).Infof(
			"Making allocation request for devices %v for device plugin %s",
			devs, resource,
		)
		resp, err := eI.e.allocate(devs)
		metrics.DevicePluginAllocationDuration.WithLabelValues(resource).Observe(metrics.SinceInSeconds(startRPCTime))
		metrics.DeprecatedDevicePluginAllocationLatency.WithLabelValues(resource).Observe(metrics.SinceInMicroseconds(startRPCTime))
		if err != nil {
			// In case of allocation failure, we want to restore m.allocatedDevices
			// to the actual allocated state from m.podDevices.
			m.mutex.Lock()
			m.allocatedDevices = m.podDevices.devices()
			m.mutex.Unlock()
			return err
		}

		if len(resp.ContainerResponses) == 0 {
			return fmt.Errorf("no containers return in allocation response %v", resp)
		}

		// 将 pod/container 与分配给ta的设备列表信息, 写入缓存.
		//
		// Update internal cached podDevices state.
		m.mutex.Lock()
		m.podDevices.insert(
			podUID, contName, resource, allocDevices, resp.ContainerResponses[0],
		)
		m.mutex.Unlock()
	}

	// Checkpoints device to container allocation information.
	return m.writeCheckpoint()
}

// 返回目标 container 传入 runc 的启动选项(尤其是 device, mount 信息)
//
// caller:
// 	1. pkg/kubelet/cm/container_manager_linux.go -> containerManagerImpl.GetResources()
// 	kubelet 在调用 docker/cri 接口创建容器前会调用到这里.
//  只有这一处
//
// GetDeviceRunContainerOptions checks whether we have cached containerDevices
// for the passed-in <pod, container> and returns its DeviceRunContainerOptions
// for the found one. An empty struct is returned in case no cached state is found.
func (m *ManagerImpl) GetDeviceRunContainerOptions(
	pod *v1.Pod, container *v1.Container,
) (*DeviceRunContainerOptions, error) {
	podUID := string(pod.UID)
	contName := container.Name
	// 一般来说, Pod 在被创建前的准入阶段, kubelet 就已经分配过设备了, 不需要重新分配.
	needsReAllocate := false
	// 注意: k 的值不只可以是 cpu, memory, 还可以是 storage 和 ephemeral-storage.
	// 其中, 后两者就可能是 device 类型, 因为这两个是可能需要从宿主机上挂载的设备.
	for k := range container.Resources.Limits {
		resource := string(k)
		if !m.isDevicePluginResource(resource) {
			continue
		}
		err := m.callPreStartContainerIfNeeded(podUID, contName, resource)
		if err != nil {
			return nil, err
		}
		// 如果从本地缓存中没找到对应的设备信息, 则需要重新分配.
		//
		// This is a device plugin resource yet we don't have cached resource state.
		// This is likely due to a race during node restart.
		// We re-issue allocate request to cover this race.
		if m.podDevices.containerDevices(podUID, contName, resource) == nil {
			needsReAllocate = true
		}
	}
	// 调用 device plugin 进行分配.
	if needsReAllocate {
		klog.V(2).Infof("needs re-allocate device plugin resources for pod %s", podUID)
		if err := m.allocatePodResources(pod); err != nil {
			return nil, err
		}
	}
	// 分配完成后, 绑定信息被储存到 podDevices 缓存中, 需要从里面读取.
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// 返回目标 container 传入 runc 的启动选项
	return m.podDevices.deviceRunContainerOptions(string(pod.UID), container.Name), nil
}

// callPreStartContainerIfNeeded issues PreStartContainer grpc call for device plugin resource
// with PreStartRequired option set.
func (m *ManagerImpl) callPreStartContainerIfNeeded(podUID, contName, resource string) error {
	m.mutex.Lock()
	eI, ok := m.endpoints[resource]
	if !ok {
		m.mutex.Unlock()
		return fmt.Errorf("endpoint not found in cache for a registered resource: %s", resource)
	}

	if eI.opts == nil || !eI.opts.PreStartRequired {
		m.mutex.Unlock()
		klog.V(4).Infof("Plugin options indicate to skip PreStartContainer for resource: %s", resource)
		return nil
	}

	devices := m.podDevices.containerDevices(podUID, contName, resource)
	if devices == nil {
		m.mutex.Unlock()
		return fmt.Errorf(
			"no devices found allocated in local cache for pod %s, container %s, resource %s",
			podUID, contName, resource,
		)
	}

	m.mutex.Unlock()
	devs := devices.UnsortedList()
	klog.V(4).Infof("Issuing an PreStartContainer call for container, %s, of pod %s", contName, podUID)
	_, err := eI.e.preStartContainer(devs)
	if err != nil {
		return fmt.Errorf("device plugin PreStartContainer rpc failed with err: %v", err)
	}
	// TODO: Add metrics support for init RPC
	return nil
}

// 更新 node.status.allocatable 字段
//
// 	@param node: 本次调度的 Pod 所在的 Node 节点对象.
//
// caller:
// 	1. ManagerImpl.Allocate()
//
// sanitizeNodeAllocatable scans through allocatedDevices in the device manager
// and if necessary, updates allocatableResource in nodeInfo to at least equal to
// the allocated capacity. This allows pods that have already been scheduled on
// the node to pass GeneralPredicates admission checking even upon device plugin failure.
func (m *ManagerImpl) sanitizeNodeAllocatable(node *schedulernodeinfo.NodeInfo) {
	var newAllocatableResource *schedulernodeinfo.Resource
	// allocatableResource node.status.allocatable{} 块的内容
	allocatableResource := node.AllocatableResource()
	if allocatableResource.ScalarResources == nil {
		allocatableResource.ScalarResources = make(map[v1.ResourceName]int64)
	}
	// 这个更新策略应该不常见.
	// 这个循环的意思是, 在 node.status.allocatable 中,
	// 扩展资源值原本为 10, 但是 device manager 检测到已分配给 Pod 的设备数量为 15,
	// 这种情况下需要将 node.status.allocatable 中的扩展资源值至少设置为 15...
	for resource, devices := range m.allocatedDevices {
		needed := devices.Len()
		quant, ok := allocatableResource.ScalarResources[v1.ResourceName(resource)]
		if ok && int(quant) >= needed {
			continue
		}
		// Needs to update nodeInfo.AllocatableResource to make sure
		// NodeInfo.allocatableResource at least equal to the capacity already allocated.
		if newAllocatableResource == nil {
			newAllocatableResource = allocatableResource.Clone()
		}
		newAllocatableResource.ScalarResources[v1.ResourceName(resource)] = int64(needed)
	}
	if newAllocatableResource != nil {
		node.SetAllocatableResource(newAllocatableResource)
	}
}

// isDevicePluginResource 判断目标资源类型是否有对应的 device plugin 注册.
// cpu/memory 是内置资源类型, 这里会返回 false.
//
// 	@param resource: 扩展的资源类型名称, 与 cpu/memory/ephemeral-storage 同级
func (m *ManagerImpl) isDevicePluginResource(resource string) bool {
	_, registeredResource := m.healthyDevices[resource]
	_, allocatedResource := m.allocatedDevices[resource]
	// Return true if this is either an active device plugin resource or
	// a resource we have previously allocated.
	if registeredResource || allocatedResource {
		return true
	}
	return false
}

// GetDevices returns the devices used by the specified container
func (m *ManagerImpl) GetDevices(podUID, containerName string) []*podresourcesapi.ContainerDevices {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.podDevices.getContainerDevices(podUID, containerName)
}

// ShouldResetExtendedResourceCapacity returns whether the extended resources should be zeroed or not,
// depending on whether the node has been recreated. Absence of the checkpoint file strongly indicates the node
// has been recreated.
func (m *ManagerImpl) ShouldResetExtendedResourceCapacity() bool {
	if utilfeature.DefaultFeatureGate.Enabled(features.DevicePlugins) {
		checkpoints, err := m.checkpointManager.ListCheckpoints()
		if err != nil {
			return false
		}
		return len(checkpoints) == 0
	}
	return false
}
